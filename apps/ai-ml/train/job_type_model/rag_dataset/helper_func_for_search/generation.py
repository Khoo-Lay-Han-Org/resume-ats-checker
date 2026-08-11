from ..service.editor_model import *
import json


def generate_label_ideation(query):
    result = agent_label_ideator.invoke(
        {"messages": [{"role": "user", "content": query}]}
    )

    for block in result["messages"][-1].content_blocks:
        if block.get("type") == "text":
            return block["text"]

    unparsed_labels = result["messages"][-1].content_blocks[-1].get("text", "")

    ideated_labels = json.loads(unparsed_labels)

    search_labels = []
    for item in ideated_labels:
        search_labels.append(item)

    return search_labels


def polish_extracted_sentence():
    pass
