# Contains a python object representing the metadata client
# This should handle anything related to the job itself, like registering the job, updating metadata, etc
# This should not be used for logging, metrics, etc
import requests
import json
import logging
import os
import validators
import logging_loki
import logging
from .. import logger


class Client:
    def __init__(self):
        # We are going to assume that the user has set the credentials in the environment
        # There are other flows but it's the easiest one
        login_string = os.environ.get('GF_AI_TRAINING_CREDS')
        self.custom_logger = logging.getLogger(__name__ + "-ai-training-o11y-client")
        self.set_credentials(login_string)
    
    def set_credentials(self, login_string):
        if not login_string or type(login_string) != str:
            logger.error("No login string provided, please set GF_AI_TRAINING_CREDS environment variable")
            return False
        # Count @ characters in the login string, should be 1
        if login_string.count("@") != 1:
            logger.error("Invalid login string format")
            return False

        token, url = login_string.split("@")
        # Check that the token is exactly 40 characters of hex
        if len(token) != 40 or not all(c in "0123456789abcdef" for c in token):
            logger.error("Invalid token format")
            return "Invalid token format"
        # Validate that the url is a url
        if not validators.url(url):
            logger.error("Invalid url format")
            return "Invalid url format"

        self.url = 'https://' + url
        self.token = token
        return True

    # Returns a boolean indicating if the process was successfully registered
    # POST /api/v1/process/new
    # Takes a JSON object with a key “user_metadata” containing a dictionary of metadata
    # Returns a JSON object with a key “process_uuid” containing the UUID of the process
    def register_process(self, user_metadata):
        # If the process is currently registered, clear everything from it
        if self.process_uuid:
            self.process_uuid = None
            self.user_metadata = None
            self.custom_logger.handlers.clear()

        headers = {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        response = requests.post(f'{self.url}/api/v1/process/new', headers=headers, data=json.dumps(data))
        if response.status_code != 200:
            logging.error(f'Failed to register with error: {response.text}')
            return False
        process_uuid = response.json()['process_uuid']
        self.process_uuid = process_uuid
        self.user_metadata = user_metadata
        self.custom_logger.addHandler(
            logging_loki.LokiHandler(
                url=f"{self.url}/api/v1/process/{self.process_uuid}/model-metrics",
                tags={
                    # This specific label guarantees that these logs never collide with anything not from this exporter
                    "grafana-aitraining-o11y-process-uuid": self.process_uuid,
                    "log-type": "custom",
                },
                # The LokiHandler doesn't currently allow the use of token auth, we're going to have to add it
                # or write our own handler, which seems a lot less elegant
                # We definitely shouldn't use the (supported) basic auth because it's a bad pattern
                #auth= "Bearer " + self.token,
                version="1",
            )
        )
        return True
    
    # Update user_metadata information
    def update_metadata(self, process_uuid, user_metadata):
        if not process_uuid:
            logging.error("No process registered, unable to update metadata")
            return False
        headers = {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        response = requests.post(f'{self.url}/api/v1/process/{process_uuid}/update-metadata', headers=headers, data=json.dumps(data))
        if response.status_code != 200:
            logging.error(f'Failed to update metadata: {response.text}')
            return False
        return True
    
    # Report a state change to the process
    # POST /api/v1/process/{uuid}/state
    # Options are “succeeded” and “failed”
    def report_state(self, state):
        if not self.process_uuid:
            logging.error("No process registered, unable to report state")
            return False
        headers = {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'state': state
        }
        response = requests.post(f'{self.url}/api/v1/process/{self.process_uuid}/state', headers=headers, data=json.dumps(data))
        if response.status_code != 200:
            logging.error(f'Failed to report state: {response.text}')
            return False
        return True
    
    def send_log(self, log):
        if not self.process_uuid:
            logging.error("No process registered, unable to send logs")
            return False
        self.custom_logger.log(log)
        return True