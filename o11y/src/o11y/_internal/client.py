# Contains a python object representing the metadata client
# This should handle anything related to the job itself, like registering the job, updating metadata, etc
# This should not be used for logging, metrics, etc
import requests
import json
import logging
import os
from .util.validate_url import validate_url
from .. import logger
import time

class Client:
    def __init__(self):
        print("Initializing Client...")
        self.process_uuid = None
        self.user_metadata = None
        self.url = None
        self.token = None
        self.tenant_id = None
        login_string = os.environ.get('GF_AI_TRAINING_CREDS')
        self.set_credentials(login_string)

    def set_credentials(self, login_string):
        if not login_string or type(login_string) != str:
            logger.error("No login string provided, please set GF_AI_TRAINING_CREDS environment variable")
            return False
        if login_string.count("@") != 2:
            logger.error("Invalid login string format")
            return False

        token, tenant_id, url = login_string.split("@")
        if not url.startswith("http://") and not url.startswith("https://"):
            url = "http://" + url
            print(f"Updated URL to: {url}")

        self.url = url
        self.token = token
        self.tenant_id = tenant_id
        print("Credentials set successfully.")
        return True

    def register_process(self, data):
        if self.process_uuid:
            print(f"Clearing existing process UUID: {self.process_uuid}")
            self.process_uuid = None
            self.user_metadata = None

        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }

        url = f'{self.url}/api/v1/process/new'
        response = requests.post(url, headers=headers, data=json.dumps(data))

        if response.status_code != 200:
            logging.error(f'Failed to register with error: {response.text}')
            return False
        try:
            process_uuid = response.json()['data']['process_uuid']
        except:
            logging.error(f'Failed to register with error: {response.text}')
            return False
        self.process_uuid = process_uuid
        self.user_metadata = data['user_metadata']
        return True

    def update_metadata(self, process_uuid, user_metadata):
        if not process_uuid:
            logging.error("No process registered, unable to update metadata")
            return False
        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        url = f'{self.url}/api/v1/process/{process_uuid}/update-metadata'
        response = requests.post(url, headers=headers, data=json.dumps(data))

        if response.status_code != 200:
            logging.error(f'Failed to update metadata: {response.text}')
            return False
        return True

    def report_state(self, state):
        if not self.process_uuid:
            logging.error("No process registered, unable to report state")
            return False
        headers = {
            'Authorization': f'Bearer {self.tenant_id}:{self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'state': state
        }
        url = f'{self.url}/api/v1/process/{self.process_uuid}/state'
        response = requests.post(url, headers=headers, data=json.dumps(data))

        if response.status_code != 200:
            logging.error(f'Failed to report state: {response.text}')
            return False
        return True

    def send_model_metrics(self, log, *, x_axis=None):
        if not self.process_uuid:
            logging.error("No process registered, unable to send logs")
            return False
        
        timestamp = str(time.time_ns())

        metadata = {
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
        response = requests.post(
            url,
            headers={
                'Authorization': f'Bearer {self.tenant_id}:{self.token}', 'Content-Type': 'application/json'
            },
            data=json.dumps(json_data)
        )

        if response.status_code != 200:
            logging.error(f'Failed to log model metric: {response.text}')
            return False
        return True
