import json

from labels import ALL_JOB_ROLES
from util import *


def generate_dataset():
    skipped_items = []

    for _, item in enumerate(ALL_JOB_ROLES):
        print(f"\n\n\nITEM: {item}\n\n\n")

        common_search_labels = slight_change_label_search(item)

        if not common_search_labels:
            print(
                f"\n\nSkipping '{item}' due to API failure or empty search labels.\n\n"
            )
            skipped_items.append(item)
            continue

        # ── Checkpoint: which keywords are done for this role? ──────
        cp = load_checkpoint()

        # Save keyword list immediately so checkpoint exists
        cp[item] = {"keywords": common_search_labels, "completed": -1}
        save_checkpoint(cp)

        role_cp = cp.get(item, {})
        if role_cp.get("keywords") == common_search_labels:
            start_from = role_cp["completed"] + 1
        else:
            start_from = 0

        for kw_idx in range(start_from, len(common_search_labels)):
            kw = common_search_labels[kw_idx]
            print(f"\n  Keyword {kw_idx + 1}/{len(common_search_labels)}: {kw}")

            all_content = scrape_content([kw])
            polished_content = polish_scraped_content(all_content)
            embedded_data = embed_scraped_content(polished_content)
            storing_status = insert_to_vector_store(embedded_data)

            if storing_status != True:
                raise Exception(f"Failed to store data for keyword: {kw}")

            cp[item] = {"keywords": common_search_labels, "completed": kw_idx}
            save_checkpoint(cp)

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
