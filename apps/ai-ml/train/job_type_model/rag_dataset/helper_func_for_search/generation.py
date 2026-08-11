from ..service.editor_model import *
import json


def generate_label_ideation(query):
    result = agent_label_ideator.invoke(
        {"messages": [{"role": "user", "content": query}]}
    )

    unparsed_labels = None
    for block in result["messages"][-1].content_blocks:
        if block.get("type") == "text":
            unparsed_labels = block["text"]
            break

    if unparsed_labels is None:
        unparsed_labels = result["messages"][-1].content_blocks[-1].get("text", "")

    ideated_labels = json.loads(unparsed_labels)

    search_labels = []
    for item in ideated_labels:
        search_labels.append(item)

    return search_labels


def polish_extracted_sentence(queries):
    try:
        result = agent_sentence_polisher.invoke(
            {"messages": [{"role": "user", "content": queries}]}
        )
    except Exception:
        return []

    unparsed_sentences = None
    for block in result["messages"][-1].content_blocks:
        if block.get("type") == "text":
            unparsed_sentences = block["text"]
            break

    if unparsed_sentences is None:
        unparsed_sentences = result["messages"][-1].content_blocks[-1].get("text", "")

    if isinstance(unparsed_sentences, str):
        return [
            sentence.strip()
            for sentence in unparsed_sentences.splitlines()
            if sentence.strip()
        ]

    try:
        polished_sentences = json.loads(unparsed_sentences)
    except (TypeError, ValueError):
        return []

    if isinstance(polished_sentences, str):
        return [
            sentence.strip()
            for sentence in polished_sentences.splitlines()
            if sentence.strip()
        ]

    return polished_sentences
