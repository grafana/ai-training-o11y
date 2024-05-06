from .. import client
from .. import logger

def log(log):
    """
    Sends a log to the Loki server
    :param log: The log message
    :return: None
    """
    # Check that log is a dict with all keys strings not containing .
    if not isinstance(log, dict):
        logger.error("Log must be a dict")
        return False
    for key in log.keys():
        if not isinstance(key, str):
            logger.error("Keys in log must be strings")
            return False
        if chr(31) in key:
            logger.error("Keys in log must not contain the unit separator character")
            return False
    # Check that all values are numbers
    for key in log.keys():
        if not isinstance(log[key], (int, float)):
            logger.error("Values in log must be numbers")
            return False

    client.send_model_metrics(log)
