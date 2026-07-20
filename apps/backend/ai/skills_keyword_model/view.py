from ..typing import TextRequest
from .config import SkillsKeywordPipeline


def predict(data: dict) -> list:
    request = TextRequest(**data)
    text = request.text

    return SkillsKeywordPipeline(text)
