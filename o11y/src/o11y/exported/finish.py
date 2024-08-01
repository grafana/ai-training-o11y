# SPDX-License-Identifier: Apache-2.0
from .. import client

def finish():
    """
    Finishes the current process
    :return: None
    """
    client.report_state("successful")
