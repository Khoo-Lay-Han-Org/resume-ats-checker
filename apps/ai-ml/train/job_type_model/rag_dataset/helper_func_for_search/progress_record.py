import json

def files_creation_for_record():
    with open("../storage/checkpoint.json", "w") as f:
        json.dump({}, f)
    with open("../storage/stored_data.json", "w") as f:
        json.dump({}, f)


def record_newly_generated_labels_ideation():
    pass


def record_each_polished_content():
    pass
