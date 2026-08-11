import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))

from rag_dataset.helper_func_for_build.dataset_build import (
    store_extracted_content_to_dataset,
)
from rag_dataset.helper_func_for_build.vector_search import (
    vector_search_based_on_labels,
)
from rag_dataset.helper_func_for_restoring.embed import embed_recorded_content
from rag_dataset.helper_func_for_search.store import store_content_into_vector_db


def restore():
    recorded_content = embed_recorded_content()
    store_content_into_vector_db(recorded_content)

    dataset = vector_search_based_on_labels()
    store_extracted_content_to_dataset(dataset)


if __name__ == "__main__":
    restore()
