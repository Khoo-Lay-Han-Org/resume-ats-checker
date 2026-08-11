from ..service.editor_model import *
from ..targets.keywords import *


def generate_label_ideation(query):
    result = agent_label_ideator.invoke(
        {"messages": [{"role": "user", "content": query}]}
    )


def polish_extracted_sentence():
    pass
