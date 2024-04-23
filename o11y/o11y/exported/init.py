from typing import Dict, Optional
from pydantic import BaseModel, ValidationError
from .. import client
from .. import logger

# The api itself only takes one argument, a "user_metadata" argument
# This consists entirely of string k:v pairs
class APIRequest(BaseModel):
    project: str
    run: Optional[str]
    user_metadata: Dict[str, str]

def init(*, project=None, run=None, metadata=None):
    """
    Initializes the logging client. This should be called at the beginning of the script.
    :param project: The project name
    :param run: The run name
    :param user_metadata: Any user metadata to attach to the run
    :return: None
    """
    data = {
        "user_metadata": metadata,
    }
    if project:
        data["project"] = project
    else:
        logger.warning("No project name provided, will be logged to project 'Default'.")
    if run:
        data["run"] = run

    try:
        data = APIRequest(**data)
    except ValidationError as e:
        logger.error(f"Invalid metadata: {e}")
        return False

    success = client.register_process(data)
    if not success:
        logger.error("Initialization failed, logs will NOT be sent.")
        return False
    return True