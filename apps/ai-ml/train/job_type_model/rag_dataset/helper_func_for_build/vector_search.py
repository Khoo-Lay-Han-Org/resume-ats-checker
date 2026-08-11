from ..storage.route_config import CHECKPOINT_FILE
from ..targets.keywords import ALL_JOB_ROLES
from ..service.embedding_model import model
from ..service.vector_store import client
import json
import random
import math


def vector_search_based_on_labels():
    with open(CHECKPOINT_FILE, "r") as f:
        data = json.load(f)

    dataset_data = []
    for i, label in enumerate(ALL_JOB_ROLES):
        if not isinstance(data[i], list):
            raise Exception("data is not list")

        label_ideation = []
        for item in data[i]["label"]:
            label_ideation.append(item)

        random.shuffle(label_ideation)

        half_length = math.ceil(len(label_ideation) / 2)

        shuffled_data = data[:half_length]

        vectors = []
        for item in shuffled_data:
            vector = model.encode(item, truncate_dim=50, normalize_embeddings=True)
            vectors.append(vector)

        results = client.search(
            collection_name="data",
            data=[vectors],
            limit=300,
            offset=0,
            output_fields=["text"],
        )

        all_sentences = []
        for result in results:
            for item in result:
                all_sentences.append(item)

        random.shuffle(all_sentences)

        shuffled_sentence = all_sentences[:300]

        dataset_data.append({"class": label, "feature": shuffled_sentence})

    return dataset_data
