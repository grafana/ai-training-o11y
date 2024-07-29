import pytest
from unittest.mock import patch
from o11y.exported.log import log

@pytest.fixture
def mock_client():
    with patch('o11y.exported.log.client') as mock:
        yield mock

@pytest.fixture
def mock_logger():
    with patch('o11y.exported.log.logger') as mock:
        yield mock

def test_log_valid_input(mock_client):
    valid_log = {"metric1": 10, "metric2": 20.5}
    assert log(valid_log) is True
    mock_client.send_model_metrics.assert_called_once_with(valid_log)

def test_log_invalid_input_not_dict(mock_logger):
    assert log("not a dict") is False
    mock_logger.error.assert_called_once_with("Log must be a dict")

def test_log_invalid_input_non_string_keys(mock_logger):
    invalid_log = {1: 10, "metric2": 20}
    assert log(invalid_log) is False
    mock_logger.error.assert_called_once_with("Log must contain only string keys and numeric values")

def test_log_invalid_input_non_numeric_values(mock_logger):
    invalid_log = {"metric1": "10", "metric2": 20}
    assert log(invalid_log) is False
    mock_logger.error.assert_called_once_with("Log must contain only string keys and numeric values")

def test_log_with_valid_x_axis(mock_client):
    valid_log = {"metric1": 10, "metric2": 20.5}
    valid_x_axis = {"step": 1}
    assert log(valid_log, x_axis=valid_x_axis) is True
    mock_client.send_model_metrics.assert_called_once_with(valid_log, x_axis=valid_x_axis)

def test_log_with_invalid_x_axis_not_dict(mock_logger):
    valid_log = {"metric1": 10}
    invalid_x_axis = "not a dict"
    assert log(valid_log, x_axis=invalid_x_axis) is False
    mock_logger.error.assert_called_once_with("x_axis must be a dict with one key")

def test_log_with_invalid_x_axis_multiple_keys(mock_logger):
    valid_log = {"metric1": 10}
    invalid_x_axis = {"step": 1, "extra": 2}
    assert log(valid_log, x_axis=invalid_x_axis) is False
    mock_logger.error.assert_called_once_with("x_axis must be a dict with one key")

def test_log_with_invalid_x_axis_non_string_key(mock_logger):
    valid_log = {"metric1": 10}
    invalid_x_axis = {1: 1}
    assert log(valid_log, x_axis=invalid_x_axis) is False
    mock_logger.error.assert_called_once_with("x_axis must have a string key and a numeric value")

def test_log_with_invalid_x_axis_non_numeric_value(mock_logger):
    valid_log = {"metric1": 10}
    invalid_x_axis = {"step": "1"}
    assert log(valid_log, x_axis=invalid_x_axis) is False
    mock_logger.error.assert_called_once_with("x_axis must have a string key and a numeric value")

def test_log_with_x_axis_key_in_log_different_value(mock_logger):
    valid_log = {"metric1": 10, "step": 2}
    invalid_x_axis = {"step": 1}
    assert log(valid_log, x_axis=invalid_x_axis) is False
    mock_logger.error.assert_called_once_with("x_axis key must not be in your metrics, or must have the same value")

def test_log_with_x_axis_key_in_log_same_value(mock_client):
    valid_log = {"metric1": 10, "step": 1}
    valid_x_axis = {"step": 1}
    assert log(valid_log, x_axis=valid_x_axis) is True
    mock_client.send_model_metrics.assert_called_once_with(valid_log, x_axis=valid_x_axis)
