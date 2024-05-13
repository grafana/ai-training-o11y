from .. import client
from .. import logger

def log(log, *, x_axis=None):
    """
    Sends a log to the Loki server
    :param log: The log message
    :return: None
    """
    # Check that log is a dict with all keys strings
    if not isinstance(log, dict):
        logger.error("Log must be a dict")
        return False
    for key in log.keys():
        if not isinstance(key, str):
            logger.error("Keys in log must be strings")
            return False

    # Check that all values are numbers
    for key in log.keys():
        if not isinstance(log[key], (int, float)):
            logger.error("Values in log must be numbers")
            return False


    if x_axis is None:
        client.send_model_metrics(log, section=section)
        return

    # Check if x_axis exists
    if not isinstance(x_axis, dict) or len(x_axis) != 1:
        logger.error("x_axis must be a dict with one key")
        return False

    # Check that x_axis' single key is a string, and its value is a number
    x_key = list(x_axis.keys())[0]
    x_value = x_axis[x_key]
    if not isinstance(x_key, str):
        logger.error("x_axis key must be a string")
        return False
    if not isinstance(x_value, (int, float)):
        logger.error("x_axis value must be a number")
        return False
    
    # Check that this key is not already in the log line
    if x_key in log.keys() and x_value != log[x_key]:
        logger.error("x_axis key must not be in your metrics, or must have the same value")
        return False

    client.send_model_metrics(log, x_axis=x_axis)
