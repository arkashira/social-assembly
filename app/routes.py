from fastapi import APIRouter, HTTPException
from .models import Post, PostStore

router = APIRouter(prefix="/posts", tags=["posts"])
store = PostStore()

@router.get("/", response_model=list[Post])
def list_posts():
    """Return all posts."""
    return store.all()

@router.post("/", response_model=Post, status_code=201)
def create_post(post: Post):
    """Create a new post – simple in‑memory store."""
    if store.exists(post.id):
        raise HTTPException(status_code=400, detail="Post ID already exists")
    store.add(post)
    return post