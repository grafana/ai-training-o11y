from .. import client

def log(log):
    """
    Sends a log to the Loki server
    :param log: The log message
    :return: None
    """
    # In principle we should validate that this is json in the future, but for right now let's not
    client.send_log(log)