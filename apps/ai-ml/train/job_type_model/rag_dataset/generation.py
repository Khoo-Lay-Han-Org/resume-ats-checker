import json

from labels import ALL_JOB_ROLES
from util import *


def generate_dataset():
    completed = load_checkpoint()
    skipped_items = []

    for _, item in enumerate(ALL_JOB_ROLES):
        if item in completed:
            print(f"\n\n\nSKIPPING (already checkpointed): {item}\n\n\n")
            continue

        print(f"\n\n\nITEM: {item}\n\n\n")

        common_search_labels = slight_change_label_search(item)

        if not common_search_labels:
            print(
                f"\n\nSkipping '{item}' due to API failure or empty search labels.\n\n"
            )
            skipped_items.append(item)
            continue

        all_content = scrape_content(common_search_labels)
        polished_content = polish_scraped_content(all_content)
        embedded_data = embed_scraped_content(polished_content)
        storing_status = insert_to_vector_store(embedded_data)

        if storing_status != True:
            raise Exception("Failed to store data")

        completed.add(item)
        save_checkpoint(completed)

    data = []
    for _, item in enumerate(ALL_JOB_ROLES):
        closest_sentences = vector_search(item)
        single_data = {"job_type": item, "sentence": closest_sentences}
        data.append(single_data)

    with open("dataset/dataset.json", "w") as file:
        json.dump(data, file, indent=2)

    if skipped_items:
        print(f"\n\nSkipped {len(skipped_items)} items due to API errors:")
        for skipped in skipped_items:
            print(f"  - {skipped}")


generate_dataset()
