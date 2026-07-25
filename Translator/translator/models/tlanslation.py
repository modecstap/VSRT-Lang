from pydantic import BaseModel


class Translation(BaseModel):
    original: str
    translations: list[str]
