from .. import client
from .. import logger


def init(*, project=None, run=None, metadata=None):
    """
    Initializes the logging client. This should be called at the beginning of the script.
    :param project: The project name
    :param run: The run name
    :param user_metadata: Any user metadata to attach to the run
    :return: None
    """
    data = {}
    if metadata:
        for key in metadata.keys():
            if key in ["project", "run"]:
                logger.warning(f"Metadata key '{key}' is redundant with '{key}' in higher scope. It may not behave as expected")
            if "." in key:
                logger.error(f"Metadata key '{key}' contains a period, which is not allowed.")
                return False
        data['user_metadata'] = metadata
    if project:
        data["project"] = project
    else:
        logger.warning("No project name provided, will be logged to project 'Default'.")
        data["project"] = "Default"
    if run:
        data["run"] = run

    success = client.register_process(data)
    if not success:
        logger.error("Initialization failed, logs will NOT be sent.")
        return False
    return True
