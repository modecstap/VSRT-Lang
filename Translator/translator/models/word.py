from pydantic import BaseModel


class Word(BaseModel):
    content: str


