from pydantic import BaseModel


class TextRequest(BaseModel):
    text: str


class DoubleTextRequest(BaseModel):
    text1: str
    text2: str
