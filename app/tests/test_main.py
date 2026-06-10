import pytest
from httpx import AsyncClient
from app.main import create_app

@pytest.mark.asyncio
async def test_create_and_list_post():
    app = create_app()
    async with AsyncClient(app=app, base_url="http://test") as ac:
        # list should be empty
        r = await ac.get("/posts/")
        assert r.status_code == 200
        assert r.json() == []

        # create a post
        payload = {"id": 1, "author": "bob", "content": "hi"}
        r = await ac.post("/posts/", json=payload)
        assert r.status_code == 201
        assert r.json() == payload

        # list now contains the post
        r = await ac.get("/posts/")
        assert r.json() == [payload]