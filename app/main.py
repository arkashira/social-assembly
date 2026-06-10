from fastapi import FastAPI
from .routes import router

def create_app() -> FastAPI:
    """Factory that returns a fully configured FastAPI app."""
    app = FastAPI(
        title="Social Assembly",
        description="Demo of a decentralized social network core.",
        version="0.1.0",
    )
    app.include_router(router)
    return app

app = create_app()