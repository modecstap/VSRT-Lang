import asyncio

from server.fastapi_server import FastAPIServer
from server.server_config import ServerConfig


def main():
    server_config = ServerConfig.from_env()
    server = FastAPIServer(server_config)

    asyncio.run(
        server.start()
    )


if __name__ == "__main__":
    main()
