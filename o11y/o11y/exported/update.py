from typing import Dict
from pydantic import BaseModel, ValidationError
from .. import logger
from .. import client

class MetadataModel(BaseModel):
    class Config:
        # Enforce strict model definition
        extra = "forbid"

    # Define a dictionary field with string keys and string values
    data: Dict[str, str]

def update(metadata: Dict[str, str]):
    """
    Updates the user_metadata for the current process
    :return: None
    """
    # Validate the metadata
    try:
        validated_metadata = MetadataModel(data=metadata)
    except ValidationError as e:
        logger.error(f"Invalid metadata: {e}")
        return False
    client.update_metadata(validated_metadata["data"])
    return True