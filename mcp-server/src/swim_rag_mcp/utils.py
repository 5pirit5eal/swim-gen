import os

import httpx
from google.auth.transport.requests import Request
from google.oauth2 import id_token as token


def get_auth_token(url: str) -> httpx.BasicAuth | None:
    """Get the authorization token for the Swim RAG API.

    Args:
        url (str): The URL of the Swim RAG API.

    Returns:
        httpx.BasicAuth | None: The authorization token or None if not applicable.
    """
    if os.getenv("K_SERVICE"):
        # Add authorization headers from the service account in env as the service runs in google cloud run
        auth_req = Request()
        id_token = token.fetch_id_token(auth_req, url)
        auth = httpx.BasicAuth(username="Bearer Token", password=id_token)  # type: ignore
    else:
        # Expect local proxy of cloud run service
        auth = None
    return auth
