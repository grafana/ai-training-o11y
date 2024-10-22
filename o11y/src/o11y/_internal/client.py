# SPDX-License-Identifier: Apache-2.0
# Contains a python object representing the metadata client
# This should handle anything related to the job itself, like registering the job, updating metadata, etc
# This should not be used for logging, metrics, etc
import json
import os
import time
import warnings
from typing import Any, Dict, Optional
from urllib.parse import ParseResult, urlparse

import requests

from .. import logger


class Client:
    def __init__(self):
        self.process_uuid = None
        self.user_metadata = None
        self.step = 1
        # TODO: Should we require a URL when creating the client instead of via set_credentials?
        self.url: ParseResult = ParseResult('', '', '', '', '', '')
        login_string = os.environ.get('GF_AI_TRAINING_CREDS')
        self.set_credentials(login_string)

    def set_credentials(self, login_string: Optional[str]) -> bool:
        if not login_string or not isinstance(login_string, str):
            warnings.warn("No login string provided, please set GF_AI_TRAINING_CREDS environment variable")
            return False

        try:
            self.url = self._parse_login_string(login_string)
            return True
        except Exception as e:
            warnings.warn(f"Invalid login string: {str(e)}")
            return False


    def _parse_login_string(self, login_string: str) -> ParseResult:
        parsed_url = urlparse(login_string)
        if not parsed_url.hostname:
            raise ValueError("Invalid login string format. Could not parse hostname from the URL.")
        if parsed_url.scheme != "http" and parsed_url.scheme != "https":
            raise ValueError("Invalid login string format. Scheme must be http or https.")

        return parsed_url


    def register_process(self, data):
        if self.process_uuid:
            self.process_uuid = None
            self.user_metadata = None
            self.step = 1

        headers = {
            'Content-Type': 'application/json'
        }

        url = f'{self.url.geturl()}/api/v1/process/new'
        
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
        logger.info(f"Process registered successfully. UUID: {self.process_uuid}")
        return True

    def update_metadata(self, process_uuid: str, user_metadata: Dict[str, Any]) -> bool:
        if not process_uuid:
            logger.error("No process registered, unable to update metadata")
            return False
        headers = {
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        url = f'{self.url.geturl()}/api/v1/process/{process_uuid}/update-metadata'
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
            'Content-Type': 'application/json'
        }
        data = {
            'state': state
        }
        url = f'{self.url.geturl()}/api/v1/process/{self.process_uuid}/state'
        response = requests.post(url, headers=headers, json=data)

        if response.status_code != 200:
            logger.error(f'Failed to report state: {response.text}')
            return False
        return True

    def send_model_metrics(self, log: Dict[str, Any], *, x_axis: Optional[Dict[str, Any]] = None) -> bool:
        if not self.process_uuid:
            logger.error("No process registered, unable to send logs")
            return False

        if not x_axis:
            x_axis = {
                "step": self.step
            }
            self.step += 1

        step_name, step_value = next(iter(x_axis.items()))

        json_data = [{
            "step_name": step_name,
            "step_value": step_value,
            "metrics": log
        }]

        url = f'{self.url.geturl()}/api/v1/process/{self.process_uuid}/model-metrics'

        headers = {
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
