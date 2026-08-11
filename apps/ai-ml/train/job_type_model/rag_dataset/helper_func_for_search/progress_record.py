import json


def files_creation_for_record():
    with open("../storage/checkpoint.json", "w") as f:
        json.dump([], f)
    with open("../storage/stored_data.json", "w") as f:
        json.dump([], f)


def record_ideated_labels(queries):
    with open("../storage/checkpoint.json", "r") as f:
        data = json.load(f)


## Use this everytime done recording polished content
def record_label_indexing():
    pass


def record_each_polished_content(queries):
    with open("../storage/stored_data.json", "r") as f:
        data = json.load(f)

    if not isinstance(data, list):
        raise Exception("data is not list")

    for query in queries:
        data.append(query)

    with open("../storage/stored_data.json", "w") as f:
        json.dump(data, f)
