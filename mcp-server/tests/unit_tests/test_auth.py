import os
import time
from unittest.mock import AsyncMock, MagicMock, patch

import pytest

from swim_rag_mcp import auth


@pytest.mark.asyncio
async def test_id_token_manager_get_id_token_cached(monkeypatch):
    manager = auth.IdTokenManager()
    manager._cached_id_token = "token123"
    manager._token_expiry_time = float("inf")
    # Mock get_context
    ctx = MagicMock()
    ctx.debug = AsyncMock()
    monkeypatch.setattr(auth, "get_context", lambda: ctx)
    # Should return cached token
    token = await manager.get_id_token("audience")
    assert token == "token123"
    ctx.debug.assert_awaited_with("Using cached ID token.")


@pytest.mark.asyncio
async def test_id_token_manager_get_id_token_fetch(monkeypatch):
    manager = auth.IdTokenManager()
    manager._cached_id_token = None
    manager._token_expiry_time = 0
    ctx = MagicMock()
    ctx.debug = AsyncMock()
    monkeypatch.setattr(auth, "get_context", lambda: ctx)
    # Patch fetch_id_token
    with patch(
        "src.swim_rag_mcp.auth.token.fetch_id_token", return_value="newtoken"
    ) as mock_fetch:
        token = await manager.get_id_token("audience")
        assert token == "newtoken"
        mock_fetch.assert_called_once()
        assert manager._cached_id_token == "newtoken"
        assert manager._token_expiry_time > time.time()
        ctx.debug.assert_awaited()


@pytest.mark.asyncio
async def test_id_token_manager_get_id_token_error(monkeypatch):
    manager = auth.IdTokenManager()
    manager._cached_id_token = None
    manager._token_expiry_time = 0
    ctx = MagicMock()
    ctx.debug = AsyncMock()
    ctx.error = AsyncMock()
    monkeypatch.setattr(auth, "get_context", lambda: ctx)
    # Patch fetch_id_token to raise
    with patch(
        "src.swim_rag_mcp.auth.token.fetch_id_token",
        side_effect=Exception("fail"),
    ):
        with pytest.raises(Exception):
            await manager.get_id_token("audience")
        ctx.error.assert_awaited()


@pytest.mark.asyncio
async def test_id_token_manager_invalidate_token(monkeypatch):
    manager = auth.IdTokenManager()
    manager._cached_id_token = "token"
    manager._token_expiry_time = 12345
    ctx = MagicMock()
    ctx.info = AsyncMock()
    monkeypatch.setattr(auth, "get_context", lambda: ctx)
    await manager.invalidate_token()
    assert manager._cached_id_token is None
    assert manager._token_expiry_time == 0
    ctx.info.assert_awaited_with("Cached token invalidated.")
