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

        for item in data:
            comparison_data = [text, item]

            embeddings = similarity_model.encode(comparison_data)

            score = util.cos_sim(embeddings[0], embeddings[1])

            if not score >= 0.92:
                client.insert(collection_name="data", data=query)
