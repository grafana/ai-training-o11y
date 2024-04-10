# login.py

from hashlib import md5
import logging

# Create a logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

# Create a console handler and set level to INFO
ch = logging.StreamHandler()
ch.setLevel(logging.INFO)

# Create a formatter
formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')

# Add formatter to ch
ch.setFormatter(formatter)

# Add ch to logger
logger.addHandler(ch)

def login(login_string: str):
    """
    Logs in using the provided token.

    Parameters:
    login_string: str: A string of the form token@url, where token is a 40-character hex string and url is a valid URL.

    Returns:
    bool: True if login was successful, False otherwise.
    """
    success = False
    try:
        token, url = login_string.split("@")
        # Check that the token is exactly 40 characters of hex
        if len(token) != 40 or not all(c in "0123456789abcdef" for c in token):
            logger.error("Invalid token format")
            return "Invalid token format"
        # TODO: also validate url
    except ValueError:
        logger.error("Invalid login string format")
        return False
    token_hash = md5(token.encode()).hexdigest()

    logger.info(f"Attempting login with hashed token: {token_hash}")
    # TODO: Api call to metadata service to login
    # We just assume it succeeds here
    success = True

    return success