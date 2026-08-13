from ..storage.route_config import DATASET
import json


def store_extracted_content_to_dataset(query):
    with open(DATASET, "w") as f:
        json.dump(query, f)
