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
        assert client.user_id == "12345", f"Expected user_id to be '12345', but got '{client.user_id}'"
        assert client.url == "https://example.com", f"Expected url to be 'https://example.com', got '{client.url}'"

def test_parse_login_string_valid(client):
    token, user_id, uri = client._parse_login_string("token123:12345@example.com")
    assert token == "token123", f"Expected token 'token123', got '{token}'"
    assert user_id == "12345", f"Expected user_id '12345', got '{user_id}'"
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

def test_validate_credentials_non_numeric_user_id(client):
    with pytest.warns((Warning, DeprecationWarning)) as warning_info:
        uri = client._validate_credentials("token123", "user123", "example.com")
    assert uri == "https://example.com", f"Expected 'https://example.com', got '{uri}'"
    assert len(warning_info) == 1, f"Expected 1 warning, got {len(warning_info)}"
    assert "Invalid user_id: must be purely numeric" in str(warning_info[0].message), f"Unexpected warning message: {warning_info[0].message}"
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
    assert client.user_id == "12345", f"Expected user_id '12345', got '{client.user_id}'"
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

def test_send_model_metrics_success(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response), \
         patch('time.time_ns', return_value=1000000000):
        result = client.send_model_metrics({"metric": "value"}, x_axis={"epoch": 1})

    assert result is True

def test_send_model_metrics_failure_no_process(client):
    result = client.send_model_metrics({"metric": "value"})
    assert result is False

def test_send_model_metrics_failure_api_error(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 400
    mock_response.text = "Error"

    with patch('requests.post', return_value=mock_response):
        result = client.send_model_metrics({"metric": "value"})

    assert result is False

def test_send_model_metrics_without_x_axis(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response), \
         patch('time.time_ns', return_value=1000000000):
        result = client.send_model_metrics({"metric": "value"})

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

def test_send_model_metrics_large_payload(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    large_log = {"metric": "x" * 1000000}  # 1MB of data

    with patch('requests.post', return_value=mock_response), \
         patch('time.time_ns', return_value=1000000000):
        result = client.send_model_metrics(large_log)

    assert result is True

def test_send_model_metrics_invalid_x_axis(client):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response), \
         patch('time.time_ns', return_value=1000000000):
        result = client.send_model_metrics({"metric": "value"}, x_axis={"key1": "value1", "key2": "value2"})

    assert result is True  # The method should still succeed, using only the first key-value pair

@pytest.mark.parametrize("invalid_log", [
    None,
    "",
    123,
    ["list", "instead", "of", "dict"],
])
def test_send_model_metrics_invalid_log(client, invalid_log):
    client.process_uuid = "test_uuid"
    mock_response = Mock()
    mock_response.status_code = 200

    with patch('requests.post', return_value=mock_response), \
         patch('time.time_ns', return_value=1000000000):
        result = client.send_model_metrics(invalid_log)

    assert result is False  # The method should fail for invalid log formats
