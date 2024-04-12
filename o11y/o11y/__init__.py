# __init__.py
import logging

# Create logger first, all other modules will use this logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Create a formatter
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Create a console handler and set level to INFO
ch = logging.StreamHandler()
ch.setLevel(logging.INFO)
ch.setFormatter(formatter)

# Add ch to logger
logger.addHandler(ch)

# Everything else depends on the metadata client, so we import it here
from ._internal.metadata_client import MetadataClient
metadata_client = MetadataClient()

from .exported.init import init
# from .exported.log import log
# from .exported.finish import finish

__all__ = [
    'init',
    'log',
    'finish'
    ]