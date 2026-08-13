import json
import sys
import time
from ..config.generation import *
import subprocess

from ..service.editor_model import agent_label_ideator, agent_sentence_polisher


def generate_label_ideation(query):
    while True:
        try:
            result = agent_label_ideator.invoke(
                {"messages": [{"role": "user", "content": query}]}
            )
            break
        except Exception as e:
            print(f"Groq need rest, waiting 5 minutes... ({e})")
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
    total = len(queries)

    for index, query in enumerate(queries, start=1):
        print(f"Polishing {index}/{total}...")

        for attempt in range(MAX_POLISH_ATTEMPTS):
            try:
                result = agent_sentence_polisher.invoke(
                    {"messages": [{"role": "user", "content": query}]}
                )
                print(result)
                break
            except Exception as e:
                print(
                    f"  Groq error (attempt {attempt + 1}/{MAX_POLISH_ATTEMPTS}): {e}"
                )
                time.sleep(POLISH_RETRY_DELAY)
        else:
            print("  Skipping item: all attempts failed")
            continue

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
            rate_limit_wording = "Rate limit reached"
            if rate_limit_wording in polished_sentence:
                print("\n\n\nRate limit reached, model is resting for 10 minutes\n\n\n")
                time.sleep(600)
                sys.exit(1)
            polished_sentences.append(polished_sentence)

    return polished_sentences
