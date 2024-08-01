# SPDX-License-Identifier: Apache-2.0
from .. import client

def update(metadata):
    """
    Updates the user_metadata for the current process
    :return: None
    """
    # TODO: Validate the metadata
    client.update_metadata(metadata)
    return True
