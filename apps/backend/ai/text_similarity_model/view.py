from sentence_transformers import util

from ..typing import DoubleTextRequest
from .config import TextSimilarityModel


def predict(data: dict) -> float:
    request = DoubleTextRequest(**data)

    emb1 = TextSimilarityModel.encode(request.text1, convert_to_tensor=True)
    emb2 = TextSimilarityModel.encode(request.text2, convert_to_tensor=True)
    score = util.pytorch_cos_sim(emb1, emb2).item()

    return score
