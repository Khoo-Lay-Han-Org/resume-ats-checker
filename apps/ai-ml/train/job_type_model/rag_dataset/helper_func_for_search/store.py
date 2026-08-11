from sentence_transformers import util
from ..service.similarity import model as similarity_model
from ..storage.route_config import STORED_DATA
from ..service.vector_store import client
from sentence_transformers import util
import json


def store_content_into_vector_db(queries):
    for query in queries:
        text = query["text"]
        embedding = query["vector"]

        with open(STORED_DATA, "r") as f:
            data = json.load(f)

        is_duplicate = False
        for item in data:
            comparison_data = [text, item]

            embeddings = similarity_model.encode(comparison_data)

            score = util.cos_sim(embeddings[0], embeddings[1])

            if score >= 0.995:
                is_duplicate = True
                break

        if not is_duplicate:
            client.insert(collection_name="data", data=query)
