from .. import metadata_client
from .. import logger

def init(*, project="Default", run=None, user_metadata=None):
    """
    Initializes the logging client. This should be called at the beginning of the script.
    :param project: The project name
    :param run: The run name
    :param user_metadata: Any user metadata to attach to the run
    :return: None
    """
    # Register the process
    data = {
        "project": project,
        "user_metadata": user_metadata,
    }
    if run:
        data["run"] = run

    success = metadata_client.register_process(user_metadata)
    if not success:
        logger.error("Failed to register process")
        return False
    return True