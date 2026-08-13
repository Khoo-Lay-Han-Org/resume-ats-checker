import json
from ..storage.route_config import CHECKPOINT_FILE, STORED_DATA


def files_creation_for_record():
    with open(CHECKPOINT_FILE, "w") as f:
        json.dump([], f)
    with open(STORED_DATA, "w") as f:
        json.dump([], f)


def record_ideated_labels(queries):
    new_data = {"label": queries, "index": 0}

    with open(CHECKPOINT_FILE, "r") as f:
        data = json.load(f)

    if not isinstance(data, list):
        raise Exception("data is not list")

    data.append(new_data)

    with open(CHECKPOINT_FILE, "w") as f:
        json.dump(data, f)


def record_label_indexing():
    with open(CHECKPOINT_FILE, "r") as f:
        data = json.load(f)

    if not isinstance(data, list):
        raise Exception("data is not list")

    data[-1]["index"] += 1

    with open(CHECKPOINT_FILE, "w") as f:
        json.dump(data, f)


def record_each_polished_content(queries):
    with open(STORED_DATA, "r") as f:
        data = json.load(f)

    if not isinstance(data, list):
        raise Exception("data is not list")

    for query in queries:
        data.append(query)

    with open(STORED_DATA, "w") as f:
        json.dump(data, f)
