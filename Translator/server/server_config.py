import os

from pydantic import BaseModel


class ServerConfig(BaseModel):

    address: str
    port: int

    @classmethod
    def from_env(cls) -> "ServerConfig":
        return cls(
            address=os.getenv("TRANS_HOST", "localhost"),
            port=int(os.getenv("TRANS_PORT", 8080)),
        )
