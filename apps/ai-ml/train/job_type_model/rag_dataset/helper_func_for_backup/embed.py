from ..service.embedding_model import model
from ..storage.route_config import STORED_DATA
import json


def embed_recorded_content():
    with open(STORED_DATA, "r") as f:
        data = json.load(f)

    all_content = []
    for item in data:
        result = model.encode(item, truncate_dim=50, normalize_embeddings=True)

        all_content.append({"text": item, "vector": result})

    return all_content
