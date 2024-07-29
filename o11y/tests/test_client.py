import unittest
from unittest.mock import patch
from o11y._internal.client import Client
import warnings

class TestClient(unittest.TestCase):

    def setUp(self):
        self.client = Client()

    def test_set_credentials_valid(self):
        with warnings.catch_warnings(record=True) as w:
            result = self.client.set_credentials("token123:12345@example.com")
            self.assertTrue(result)
            self.assertEqual(self.client.token, "token123")
            self.assertEqual(self.client.user_id, "12345")
            self.assertEqual(self.client.url, "https://example.com")
            self.assertEqual(len(w), 0)

    def test_set_credentials_invalid(self):
        with warnings.catch_warnings(record=True) as w:
            result = self.client.set_credentials("invalid_format")
            self.assertFalse(result)
            self.assertEqual(len(w), 1)

    def test_set_credentials_empty(self):
        with warnings.catch_warnings(record=True) as w:
            result = self.client.set_credentials("")
            self.assertFalse(result)
            self.assertEqual(len(w), 1)

    def test_set_credentials_none(self):
        with warnings.catch_warnings(record=True) as w:
            result = self.client.set_credentials(None)
            self.assertFalse(result)
            self.assertEqual(len(w), 1)

    def test_parse_login_string_valid(self):
        token, user_id, uri = self.client._parse_login_string("token123:12345@example.com")
        self.assertEqual(token, "token123")
        self.assertEqual(user_id, "12345")
        self.assertEqual(uri, "example.com")

    def test_parse_login_string_invalid(self):
        with self.assertRaises(ValueError):
            self.client._parse_login_string("invalid_format")

    def test_parse_login_string_missing_user_id(self):
        with self.assertRaises(ValueError):
            self.client._parse_login_string("token123@example.com")

    def test_validate_credentials_valid(self):
        uri = self.client._validate_credentials("token123", "12345", "example.com")
        self.assertEqual(uri, "https://example.com")

    def test_validate_credentials_non_numeric_user_id(self):
        with warnings.catch_warnings(record=True) as w:
            uri = self.client._validate_credentials("token123", "user123", "example.com")
            self.assertEqual(uri, "https://example.com")
            self.assertEqual(len(w), 1)

    def test_validate_credentials_http_scheme(self):
        uri = self.client._validate_credentials("token123", "12345", "http://example.com")
        self.assertEqual(uri, "http://example.com")

    def test_validate_credentials_https_scheme(self):
        uri = self.client._validate_credentials("token123", "12345", "https://example.com")
        self.assertEqual(uri, "https://example.com")

    def test_validate_credentials_invalid_scheme(self):
        with warnings.catch_warnings(record=True) as w:
            uri = self.client._validate_credentials("token123", "12345", "ftp://example.com")
            self.assertEqual(uri, "https://example.com")
            self.assertEqual(len(w), 1)

    def test_set_credentials_internal(self):
        self.client._set_credentials("token123", "12345", "https://example.com")
        self.assertEqual(self.client.token, "token123")
        self.assertEqual(self.client.user_id, "12345")
        self.assertEqual(self.client.url, "https://example.com")

if __name__ == '__main__':
    unittest.main()
