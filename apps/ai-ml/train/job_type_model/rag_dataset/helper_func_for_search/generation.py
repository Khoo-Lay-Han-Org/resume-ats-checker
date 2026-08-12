from ..service.editor_model import *
import json
import time


def generate_label_ideation(query):
    while True:
        try:
            result = agent_label_ideator.invoke(
                {"messages": [{"role": "user", "content": query}]}
            )
            break
        except:
            print("Groq need rest, waiting 5 minutes...")
            time.sleep(300)

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
    polished_sentences = []
    for query in queries:
        while True:
            try:
                result = agent_sentence_polisher.invoke(
                    {"messages": [{"role": "user", "content": query}]}
                )
                break
            except:
                print("Groq need rest, waiting 5 minutes...")
                time.sleep(300)

        unparsed_sentence = None
        for block in result["messages"][-1].content_blocks:
            if block.get("type") == "text":
                unparsed_sentence = block["text"]
                break

        if unparsed_sentence is None:
            unparsed_sentence = (
                result["messages"][-1].content_blocks[-1].get("text", "")
            )

        polished_sentence = unparsed_sentence.strip()
        if polished_sentence:
            polished_sentences.append(polished_sentence)

    return polished_sentences
