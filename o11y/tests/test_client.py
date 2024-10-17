import pytest
import json
from unittest.mock import patch, Mock
from o11y._internal.client import Client

@pytest.fixture
def client():
    client = Client()
    client.set_credentials("https://test_tenant:test_token@test.com")
    return client

@pytest.mark.parametrize("credentials, expected_result, expect_warning", [
    ("https://12345:token123@example.com", True, False),
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
        assert client.url.geturl() == "https://12345:token123@example.com", f"Expected url to be 'https://example.com', got '{client.url.geturl()}'"

def test_parse_login_string_valid(client):
    url = client._parse_login_string("http://12345:token123@example.com")
    assert url.password == "token123"
    assert url.username == "12345"
    assert url.hostname == "example.com"

@pytest.mark.parametrize("invalid_login", [
    "invalid_format",
])
def test_parse_login_string_invalid(client, invalid_login):
    with pytest.raises(ValueError, match="Invalid (login string|credentials) format"):
        client._parse_login_string(invalid_login)

@pytest.mark.parametrize("scheme, expected_uri", [
    ("http", "http://example.com"),
    ("https", "https://example.com"),
])
def test_validate_credentials_schemes(client, scheme, expected_uri):
    url = client._parse_login_string(f"{scheme}://example.com")
    assert url.geturl() == expected_uri, f"Expected '{expected_uri}', got '{url}'"

def test_validate_credentials_invalid_scheme(client):
    with pytest.raises(ValueError, match="Invalid login string format. Scheme must be http or https"):
        client._parse_login_string("ftp://12345:token123@example.com")

def test_set_credentials_internal(client):
    client.set_credentials("https://12345:token123@example.com")
    assert client.url.password == "token123", f"Expected token 'token123', got '{client.url.password}'"
    assert client.url.username == "12345", f"Expected tenant_id '12345', got '{client.url.username}'"
    assert client.url.hostname == "example.com", f"Expected hostname 'example.com', got '{client.url.hostname}'"
    assert client.url.scheme == "https", f"Expected scheme 'https', got '{client.url.scheme}'"

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
