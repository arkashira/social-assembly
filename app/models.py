from pydantic import BaseModel, Field
from typing import Dict, List

class Post(BaseModel):
    id: int = Field(..., example=1)
    author: str = Field(..., example="alice")
    content: str = Field(..., example="Hello, world!")

class PostStore:
    """Very small in‑memory store – perfect for a demo."""
    def __init__(self):
        self._store: Dict[int, Post] = {}

    def add(self, post: Post) -> None:
        self._store[post.id] = post

    def all(self) -> List[Post]:
        return list(self._store.values())

    def exists(self, post_id: int) -> bool:
        return post_id in self._store