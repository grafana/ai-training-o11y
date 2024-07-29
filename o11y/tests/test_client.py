import unittest
import os
from unittest.mock import patch
from o11y._internal.client import Client
import warnings

warnings.filterwarnings("ignore", category=DeprecationWarning, module="pkg_resources")

class TestClient(unittest.TestCase):

    def setUp(self):
        os.environ['GF_AI_TRAINING_CREDS'] = 'testtoken:12345@http://example.com'
        self.client = Client()

    def test_valid_credentials(self):
        login_string = "mytoken:12345@http://example.com"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.token, "mytoken")
        self.assertEqual(self.client.user_id, "12345")
        self.assertEqual(self.client.url, "http://example.com")

    def test_valid_credentials_with_https(self):
        login_string = "mytoken:12345@https://example.com"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.url, "https://example.com")

    def test_valid_credentials_without_http(self):
        login_string = "mytoken:12345@example.com"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.url, "http://example.com")

    def test_invalid_login_string_format(self):
        login_string = "invalid_format"
        result = self.client.set_credentials(login_string)
        self.assertFalse(result)

    def test_empty_login_string(self):
        login_string = ""
        result = self.client.set_credentials(login_string)
        self.assertFalse(result)

    def test_none_login_string(self):
        login_string = None
        result = self.client.set_credentials(login_string)
        self.assertFalse(result)

    def test_complex_url(self):
        login_string = "mytoken:12345@http://example.com:8080/path?query=value"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.url, "http://example.com:8080/path?query=value")

    def test_special_characters(self):
        login_string = "my!token$:user@id#@http://ex&ample.com/path?query=value"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.token, "my!token$")
        self.assertEqual(self.client.user_id, "user@id#")
        self.assertEqual(self.client.url, "http://ex&ample.com/path?query=value")

    def test_ip_address_url(self):
        login_string = "mytoken:12345@http://192.168.1.1:8080"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.url, "http://192.168.1.1:8080")

    @patch('o11y._internal.client.logger')
    def test_logger_error_for_invalid_format(self, mock_logger):
        login_string = "invalid_format"
        result = self.client.set_credentials(login_string)
        self.assertFalse(result)
        mock_logger.error.assert_called_with("Invalid login string format. Expected format: token:user_id@uri")

    @patch('o11y._internal.client.logger')
    def test_logger_error_for_empty_string(self, mock_logger):
        login_string = ""
        result = self.client.set_credentials(login_string)
        self.assertFalse(result)
        mock_logger.error.assert_called_with("No login string provided, please set GF_AI_TRAINING_CREDS environment variable")

    def test_environment_variable_usage(self):
        os.environ['GF_AI_TRAINING_CREDS'] = 'envtoken:67890@http://envexample.com'
        client = Client()  # This should use the environment variable
        self.assertEqual(client.token, "envtoken")
        self.assertEqual(client.user_id, "67890")
        self.assertEqual(client.url, "http://envexample.com")

    def test_very_long_token_and_user_id(self):
        long_token = "a" * 1000
        long_user_id = "b" * 1000
        login_string = f"{long_token}:{long_user_id}@http://example.com"
        result = self.client.set_credentials(login_string)
        self.assertTrue(result)
        self.assertEqual(self.client.token, long_token)
        self.assertEqual(self.client.user_id, long_user_id)

if __name__ == '__main__':
    unittest.main()
