# __init__.py
import logging

from exported.init import init
from exported.log import log
from exported.finish import finish
from _internal import MetadataClient

# Create a logger
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

metadata_client = MetadataClient()

__all__ = [
    'init',
    'log',
    'finish'
    ]
