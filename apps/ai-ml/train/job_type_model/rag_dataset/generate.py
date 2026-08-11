import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from rag_dataset.helper_func_for_build.dataset_build import (
    store_extracted_content_to_dataset,
)
from rag_dataset.helper_func_for_build.vector_search import (
    vector_search_based_on_labels,
)
from rag_dataset.helper_func_for_search.embed import embed_scraped_content
from rag_dataset.helper_func_for_search.generation import (
    generate_label_ideation,
    polish_extracted_sentence,
)
from rag_dataset.helper_func_for_search.progress_record import (
    files_creation_for_record,
    record_each_polished_content,
    record_ideated_labels,
    record_label_indexing,
)
from rag_dataset.helper_func_for_search.store import store_content_into_vector_db
from rag_dataset.helper_func_for_search.webscrape import scrape_dataset_contents
from rag_dataset.targets.keywords import ALL_JOB_ROLES


def generate():
    files_creation_for_record()

    for role in ALL_JOB_ROLES:
        ideated_labels = generate_label_ideation(role)
        record_ideated_labels(ideated_labels)

        for ideation in ideated_labels:
            scraped_content = scrape_dataset_contents(ideation)

            if scraped_content:
                polished_content = polish_extracted_sentence("\n".join(scraped_content))
                print(polished_content)
                embedded_content = embed_scraped_content(polished_content)

                store_content_into_vector_db(embedded_content)
                record_each_polished_content(polished_content)

            record_label_indexing()

    dataset = vector_search_based_on_labels()
    store_extracted_content_to_dataset(dataset)


if __name__ == "__main__":
    generate()
