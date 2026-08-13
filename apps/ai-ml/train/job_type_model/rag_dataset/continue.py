import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from rag_dataset.searching.embed import embed_scraped_content
from rag_dataset.searching.generation import (
    generate_label_ideation,
    polish_extracted_sentence,
)
from rag_dataset.searching.progress_record import (
    record_each_polished_content,
    record_ideated_labels,
    record_label_indexing,
)
from rag_dataset.searching.store import store_content_into_vector_db
from rag_dataset.searching.webscrape import scrape_dataset_contents
from rag_dataset.storage.route_config import CHECKPOINT_FILE
from rag_dataset.targets.keywords import ALL_JOB_ROLES


def continue_pipeline():
    with open(CHECKPOINT_FILE) as f:
        checkpoint = json.load(f)

    starting_role = len(checkpoint)
    starting_ideation = 0

    for index, entry in enumerate(checkpoint):
        if entry["index"] < len(entry["label"]):
            starting_role = index
            starting_ideation = entry["index"]
            break

    for role_index, role in enumerate(ALL_JOB_ROLES):
        if role_index < starting_role:
            continue

        if role_index < len(checkpoint):
            ideated_labels = checkpoint[role_index]["label"]
        else:
            ideated_labels = generate_label_ideation(role)
            record_ideated_labels(ideated_labels)

        skip = starting_ideation if role_index == starting_role else 0

        for ideation in ideated_labels[skip:]:
            scraped_content = scrape_dataset_contents(ideation)

            if scraped_content:
                polished_content = polish_extracted_sentence(scraped_content)

                if polished_content:
                    embedded_content = embed_scraped_content(polished_content)

                    store_content_into_vector_db(embedded_content)
                    record_each_polished_content(polished_content)

            record_label_indexing()


if __name__ == "__main__":
    continue_pipeline()
