import uvicorn
from fastapi import FastAPI

from server.routers.translator_router import TranslatorRouter
from server.server_config import ServerConfig


class FastAPIServer:
    def __init__(self, server_config: ServerConfig):
        self.address = server_config.address
        self.port = server_config.port
        self.app = FastAPI()

        self._setup_routers()

    async def start(self):
        config = uvicorn.Config(
            self.app,
            host=self.address,
            port=self.port,
            loop="asyncio"
        )
        server = uvicorn.Server(config)
        await server.serve()

    def _setup_routers(self):
        self.app.include_router(TranslatorRouter(
            "/api/translator",
            ["translator"]
        ).router)