import os
import threading
import time

from fastmcp.server.dependencies import get_context
from google.auth.transport.requests import Request
from google.oauth2 import id_token as token


class IdTokenManager:
    def __init__(self):
        self._cached_id_token = None
        self._token_expiry_time = 0
        self._lock = threading.Lock()

    async def get_id_token(self, target_audience: str) -> str:
        """Get a valid ID token for the specified target audience, caching it if valid."""
        ctx = get_context()
        with self._lock:  # Acquire lock for thread-safe access
            if self._cached_id_token and self._token_expiry_time > (
                time.time() + 300
            ):
                await ctx.debug("Using cached ID token.")
                return self._cached_id_token

            await ctx.debug("Fetching new ID token...")

            try:
                self._cached_id_token = token.fetch_id_token(
                    Request(), audience=target_audience
                )
                if self._cached_id_token is None:
                    raise ValueError("Failed to fetch ID token.")
                self._token_expiry_time = time.time() + 3600  # 1 hour

                await ctx.debug(
                    f"New ID token fetched. Expires at: {time.ctime(self._token_expiry_time)}"
                )
                return self._cached_id_token
            except Exception as e:
                await ctx.error(f"Error fetching ID token: {e}")
                raise

    async def invalidate_token(self) -> None:
        """Invalidate the cached ID token."""
        ctx = get_context()
        with self._lock:  # Acquire lock for thread-safe invalidation
            self._cached_id_token = None
            self._token_expiry_time = 0
            await ctx.info("Cached token invalidated.")
