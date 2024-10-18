from typing import Dict, Union, Optional
from decimal import Decimal
# SPDX-License-Identifier: Apache-2.0
from .. import client
from .. import logger

def log(
        log_data: Dict[str, Union[int, float, Decimal]],
        *,
        x_axis: Optional[Dict[str, int]] = None
        ) -> bool:
    """
    Sends a log to the Loki server.

    Args:
        log (Dict[str, Union[int, float, Decimal]]): The log message as a dictionary with string keys and numeric values.
        x_axis (Optional[Dict[str, Union[int, float]]], optional): A single-item dictionary representing the x-axis. Defaults to None.

    Returns:
        bool: True if the log was sent successfully, False otherwise.
    """

    if not isinstance(log_data, dict):
        logger.error("Log must be a dict")
        return False

    if not all(isinstance(key, str) and isinstance(value, (int, float, Decimal)) for key, value in log_data.items()):
        logger.error("Log must contain only string keys and numeric values")
        return False
    
    if not x_axis:
        return bool(client.send_model_metrics(log_data))

    if not isinstance(x_axis, dict) or len(x_axis) != 1:
        logger.error("x_axis must be a dict with one key")
        return False

    x_key, x_value = next(iter(x_axis.items()))

    if not isinstance(x_key, str) or not isinstance(x_value, int):
        logger.error("x_axis must have a string key and a numeric value")
        return False

    if x_key in log_data and x_value != log_data[x_key]:
        logger.error("x_axis key must not be in your metrics, or must have the same value")
        return False

    return bool(client.send_model_metrics(log_data, x_axis=x_axis))
