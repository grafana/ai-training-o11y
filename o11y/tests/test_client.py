import pytest
import json
from unittest.mock import patch, Mock
from o11y._internal.client import Client

@pytest.fixture
def client():
    client = Client()
    client.tenant_id = "test_tenant"
    client.token = "test_token"
    client.url = "https://test.com"
    return client

@pytest.mark.parametrize("credentials, expected_result, expect_warning", [
    ("token123:12345@example.com", True, False),
    ("invalid_format", False, True),
    ("", False, True),
    (None, False, True),
])
def test_set_credentials(client, credentials, expected_result, expect_warning):
    if expect_warning:
        with pytest.warns((Warning)) as warning_info:
            result = client.set_credentials(credentials)
        assert len(warning_info) > 0, f"Expected a warning for input '{credentials}', but none was raised"
        assert warning_info[0].filename.endswith("o11y/_internal/client.py"), f"Warning not from expected module. Got: {warning_info[0].filename}"
        print(f"Warning type: {type(warning_info[0].message).__name__}")
        print(f"Warning message: {str(warning_info[0].message)}")
    else:
        result = client.set_credentials(credentials)

    assert result == expected_result, f"Expected set_credentials to return {expected_result}, but got {result}"

    if expected_result:
        assert client.token == "token123", f"Expected token to be 'token123', but got '{client.token}'"
        assert client.tenant_id == "12345", f"Expected tenant_id to be '12345', but got '{client.tenant_id}'"
        assert client.url == "https://example.com", f"Expected url to be 'https://example.com', got '{client.url}'"

def test_parse_login_string_valid(client):
    token, tenant_id, uri = client._parse_login_string("token123:12345@example.com")
    assert token == "token123", f"Expected token 'token123', got '{token}'"
    assert tenant_id == "12345", f"Expected tenant_id '12345', got '{tenant_id}'"
    assert uri == "example.com", f"Expected uri 'example.com', got '{uri}'"

@pytest.mark.parametrize("invalid_login", [
    "invalid_format",
    "token123@example.com",
])
def test_parse_login_string_invalid(client, invalid_login):
    with pytest.raises(ValueError, match="Invalid (login string|credentials) format"):
        client._parse_login_string(invalid_login)

def test_validate_credentials_valid(client):
    uri = client._validate_credentials("token123", "12345", "example.com")
    assert uri == "https://example.com", f"Expected 'https://example.com', got '{uri}'"

def test_validate_credentials_non_numeric_tenant_id(client):
    with pytest.warns((Warning, DeprecationWarning)) as warning_info:
        uri = client._validate_credentials("token123", "user123", "example.com")
    assert uri == "https://example.com", f"Expected 'https://example.com', got '{uri}'"
    assert len(warning_info) == 1, f"Expected 1 warning, got {len(warning_info)}"
    assert "Invalid tenant_id: must be purely numeric" in str(warning_info[0].message), f"Unexpected warning message: {warning_info[0].message}"
    print(f"Warning type: {type(warning_info[0].message).__name__}")
    print(f"Warning message: {str(warning_info[0].message)}")

@pytest.mark.parametrize("scheme, expected_uri", [
    ("http://", "http://example.com"),
    ("https://", "https://example.com"),
])
def test_validate_credentials_schemes(client, scheme, expected_uri):
    uri = client._validate_credentials("token123", "12345", f"{scheme}example.com")
    assert uri == expected_uri, f"Expected '{expected_uri}', got '{uri}'"

def test_validate_credentials_invalid_scheme(client):
    with pytest.warns((Warning, DeprecationWarning)) as warning_info:
        uri = client._validate_credentials("token123", "12345", "ftp://example.com")
    assert uri == "https://example.com", f"Expected 'https://example.com', got '{uri}'"
    assert len(warning_info) == 1, f"Expected 1 warning, got {len(warning_info)}"
    assert "Invalid URI scheme" in str(warning_info[0].message), f"Unexpected warning message: {warning_info[0].message}"
    print(f"Warning type: {type(warning_info[0].message).__name__}")
    print(f"Warning message: {str(warning_info[0].message)}")

def test_set_credentials_internal(client):
    client._set_credentials("token123", "12345", "https://example.com")
    assert client.token == "token123", f"Expected token 'token123', got '{client.token}'"
    assert client.tenant_id == "12345", f"Expected tenant_id '12345', got '{client.tenant_id}'"
    assert client.url == "https://example.com", f"Expected url 'https://example.com', got '{client.url}'"

def test_register_process_success(client):
    mock_response = Mock()
    mock_response.status_code = 200
    mock_response.json.return_value = {"data": {"process_uuid": "test_uuid"}}

    with patch('requests.post', return_value=mock_response):
        data = {"user_metadata": {"key": "value"}}
        result = client.register_process(data)

    assert result is True
    assert client.process_uuid == "test_uuid"
    assert client.user_metadata == {"key": "value"}

def test_register_process_failure(client):
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.text = "Error"

    with patch('requests.post', return_value=mock_response):
        data = {"user_metadata": {"key": "value"}}
        result = client.register_process(data)

    assert result is False
    assert client.process_uuid is None
    assert client.user_metadata is None

def test_update_metadata_success(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response):
        result = client.update_metadata("test_uuid", {"new_key": "new_value"})

    assert result is True

def test_update_metadata_failure_no_process(client):
    result = client.update_metadata(None, {"new_key": "new_value"})
    assert result is False

def test_update_metadata_failure_api_error(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.text = "Error"

    with patch('requests.post', return_value=mock_response):
        result = client.update_metadata("test_uuid", {"new_key": "new_value"})

    assert result is False

def test_report_state_success(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response):
        result = client.report_state("new_state")

    assert result is True

def test_report_state_failure_no_process(client):
    result = client.report_state("new_state")
    assert result is False

def test_report_state_failure_api_error(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.text = "Error"

    with patch('requests.post', return_value=mock_response):
        result = client.report_state("new_state")

    assert result is False

def test_register_process_clear_existing(client):
    client.process_uuid = "old_uuid"
    client.user_metadata = {"old": "data"}
    mock_response = Mock()
    mock_response.status_code = 200
    mock_response.json.return_value = {"data": {"process_uuid": "new_uuid"}}

    with patch('requests.post', return_value=mock_response):
        data = {"user_metadata": {"key": "value"}}
        result = client.register_process(data)

    assert result is True
    assert client.process_uuid == "new_uuid"
    assert client.user_metadata == {"key": "value"}

def test_register_process_invalid_json_response(client):
    mock_response = Mock()
    mock_response.status_code = 200
    mock_response.json.side_effect = json.JSONDecodeError("Invalid JSON", "", 0)

    with patch('requests.post', return_value=mock_response):
        data = {"user_metadata": {"key": "value"}}
        result = client.register_process(data)

    assert result is False
    assert client.process_uuid is None
    assert client.user_metadata is None

def test_update_metadata_empty_metadata(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response):
        result = client.update_metadata("test_uuid", {})

    assert result is True

def test_report_state_empty_state(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response):
        result = client.report_state("")

    assert result is True
