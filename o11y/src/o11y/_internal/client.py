# Contains a python object representing the metadata client
# This should handle anything related to the job itself, like registering the job, updating metadata, etc
# This should not be used for logging, metrics, etc
from typing import Optional, Tuple
import warnings
import requests
import json
import os
import time
import warnings
from typing import Any, Dict, Optional, Tuple

import requests
from urllib.parse import urlparse

from .. import logger

class Client:
    def __init__(self):
        self.process_uuid = None
        self.user_metadata = None
        self.url = None
        self.token = None
        self.tenant_id = None
        login_string = os.environ.get('GF_AI_TRAINING_CREDS')
        self.set_credentials(login_string)

    def set_credentials(self, login_string: Optional[str]) -> bool:
        if not login_string or not isinstance(login_string, str):
            logger.warning("No login string provided or invalid type")
            warnings.warn("No login string provided, please set GF_AI_TRAINING_CREDS environment variable")
            return False

        try:
            token, tenant_id, uri = self._parse_login_string(login_string)
            uri = self._validate_credentials(token, tenant_id, uri)
            self._set_credentials(token, tenant_id, uri)
            return True
        except Exception as e:
            logger.error(f"Error setting credentials: {str(e)}")
            warnings.warn(f"Invalid login string: {str(e)}")
            return False


    def _parse_login_string(self, login_string: str) -> Tuple[str, str, str]:
        parts = login_string.split('@')
        if len(parts) != 2:
            raise ValueError("Invalid login string format. Expected format: token:tenant_id@uri")
        
        credentials, uri = parts
        cred_parts = credentials.split(':')
        if len(cred_parts) != 2:
            raise ValueError("Invalid credentials format. Expected format: token:tenant_id")
        
        token, tenant_id = cred_parts
        return token.strip(), tenant_id.strip(), uri.strip()

    def _validate_credentials(self, token: str, tenant_id: str, uri: str) -> str:
        if not tenant_id.isdigit():
            warnings.warn("Invalid tenant_id: must be purely numeric")
        
        parsed_uri = urlparse(uri)
        if not parsed_uri.scheme:
            uri = "https://" + uri
        elif parsed_uri.scheme not in ["http", "https"]:
            warnings.warn(f"Invalid URI scheme '{parsed_uri.scheme}'. Using https instead.")
            uri = "https://" + parsed_uri.netloc + parsed_uri.path
        
        return uri

    def _set_credentials(self, token: str, tenant_id: str, uri: str) -> None:
        self.url = uri
        self.token = token
        self.tenant_id = tenant_id

    def register_process(self, data):
        if self.process_uuid:
            self.process_uuid = None
            self.user_metadata = None

        if not self.tenant_id or not self.token:
            logger.error("User ID or token is not set.")
            return False

        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }

        url = f'{self.url}/api/v1/process/new'
        
        try:
            response = requests.post(url, headers=headers, json=data)

            if response.status_code != 200:
                return False
            
            process_uuid = response.json()['data']['process_uuid']
        except Exception as e:
            logger.error(f"Exception during process registration: {str(e)}")
            return False

        self.process_uuid = process_uuid
        self.user_metadata = data['user_metadata']
        return True

    def update_metadata(self, process_uuid: str, user_metadata: Dict[str, Any]) -> bool:
        if not process_uuid:
            logger.error("No process registered, unable to update metadata")
            return False
        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        url = f'{self.url}/api/v1/process/{process_uuid}/update-metadata'
        response = requests.post(url, headers=headers, json=data)

        if response.status_code != 200:
            logger.error(f'Failed to update metadata: {response.text}')
            return False
        return True

    def report_state(self, state: str) -> bool:
        if not self.process_uuid:
            logger.error("No process registered, unable to report state")
            return False
        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'state': state
        }
        url = f'{self.url}/api/v1/process/{self.process_uuid}/state'
        response = requests.post(url, headers=headers, json=data)

        if response.status_code != 200:
            logger.error(f'Failed to report state: {response.text}')
            return False
        return True

    def send_model_metrics(self, log: Dict[str, Any], *, x_axis: Optional[Dict[str, Any]] = None) -> bool:
        if not self.process_uuid:
            logger.error("No process registered, unable to send logs")
            return False
        
        timestamp = str(time.time_ns())

        metadata: Dict[str, Any] = {
            "process_uuid": self.process_uuid,
            "type": "model-metrics"
        }

        if x_axis:
            x_key = list(x_axis.keys())[0]
            metadata['x_axis'] = x_key
            metadata['x_value'] = str(x_axis[x_key])

        json_data = {
            "streams": [
                {
                    "stream": {
                        "job": "o11y",
                    },
                    "values": [
                        [
                            timestamp,
                            json.dumps(log),
                            metadata,
                        ]
                    ]
                }
            ]
        }

        url = f'{self.url}/api/v1/process/{self.process_uuid}/model-metrics'

        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }

        try:
            response = requests.post(
                url,
                headers=headers,
                json=json_data
            )

            if response.status_code != 200:
                logger.error(f'Failed to log model metric. Status code: {response.status_code}, Response: {response.text}')
                return False
            return True
        except requests.exceptions.RequestException as e:
            logger.exception(f"An error occurred while sending the request: {e}")
            return False
