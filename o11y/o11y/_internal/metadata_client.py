# Contains a python object representing the metadata client
# This should handle anything related to the job itself, like registering the job, updating metadata, etc
# This should not be used for logging, metrics, etc
import requests
import json
import logging
import os
import validators
from .. import logger

class MetadataClient:
    def __init__(self):
        # We are going to assume that the user has set the credentials in the environment
        # There are other flows but it's the easiest one
        login_string = os.environ.get('GF_AI_TRAINING_CREDS')
        self.set_credentials(login_string)
    
    def set_credentials(self, login_string):
        if not login_string or type(login_string) != str:
            logger.error("No login string provided")
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

    # Function for calling an endpoint as documented here:
    def register_process(self, user_metadata):
        headers = {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }
        data = {
            'user_metadata': user_metadata
        }
        response = requests.post(f'{self.url}/api/v1/process/new', headers=headers, data=json.dumps(data))
        if response.status_code != 200:
            logging.error(f'Failed to register process: {response.text}')
            return False
        process_uuid = response.json()['process_uuid']
        self.process_uuid = process_uuid
        self.user_metadata = user_metadata
        return True
    
    # Update user_metadata information
    def update_metadata(self, process_uuid, user_metadata):
        if not process_uuid:
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