import json
import sys
import time
from ..config.generation import *
import subprocess

from ..service.editor_model import (
    AVAILABLE_MODELS,
    AGENT_SYSTEM_PROMPT,
    build_agent_model,
)


def generate_label_ideation(query):
    max_model_num = len(AVAILABLE_MODELS)
    current_model_num = 0

    system_prompt = AGENT_SYSTEM_PROMPT[0]

    while True:
        model = AVAILABLE_MODELS[current_model_num]
        agent = build_agent_model(model, system_prompt)

        try:
            result = agent.invoke({"messages": [{"role": "user", "content": query}]})
            print("Result")
            print(result)
        except Exception as e:
            print(f"Groq need rest, waiting 10 minutes... ({e})")
            current_model_num = (
                current_model_num + 1 if current_model_num < max_model_num else 0
            )
            time.sleep(600)
            continue

        unparsed_labels = None
        for block in result["messages"][-1].content_blocks:
            if block.get("type") == "text":
                unparsed_labels = block["text"]
                break

        if unparsed_labels is None:
            unparsed_labels = result["messages"][-1].content_blocks[-1].get("text", "")

        polished_labels = unparsed_labels.strip()
        if not polished_labels:
            print("Error: Empty response from agent")
            return []

        rate_limit_wording = "Rate limit reached"
        if rate_limit_wording in polished_labels:
            print("\n\n\nRate limit reached, model is resting for 240 minutes\n\n\n")
            current_model_num = (
                current_model_num + 1 if current_model_num < max_model_num else 0
            )
            time.sleep(14400)
            continue

        not_found_wording = "does not exist"
        if not_found_wording in polished_labels:
            print("\n\n\nNot found. Skipped.\n\n\n")
            current_model_num = (
                current_model_num + 1 if current_model_num < max_model_num else 0
            )
            continue

        model_not_found_wording = "model_not_found"
        if model_not_found_wording in polished_labels:
            print("\n\n\nModel not found. Skipped.\n\n\n")
            current_model_num = (
                current_model_num + 1 if current_model_num < max_model_num else 0
            )
            continue

        try:
            ideated_labels = json.loads(polished_labels)
        except json.JSONDecodeError as e:
            print(f"Error: Invalid JSON response: {e}")
            print(f"Response: {polished_labels[:200]}...")
            return []

        search_labels = []
        for item in ideated_labels:
            search_labels.append(item)

        return search_labels


def polish_extracted_sentence(queries):
    max_model_num = len(AVAILABLE_MODELS)
    current_model_num = 0

    system_prompt = AGENT_SYSTEM_PROMPT[1]

    polished_sentences = []
    total = len(queries)

    for index, query in enumerate(queries, start=1):
        model = AVAILABLE_MODELS[current_model_num]
        agent = build_agent_model(model, system_prompt)

        print(f"Polishing {index}/{total}...")

        for attempt in range(MAX_POLISH_ATTEMPTS):
            model = AVAILABLE_MODELS[current_model_num]
            agent = build_agent_model(model, system_prompt)
            try:
                result = agent.invoke(
                    {"messages": [{"role": "user", "content": query}]}
                )
                print(result)
                break
            except Exception as e:
                print(
                    f"  Groq error (attempt {attempt + 1}/{MAX_POLISH_ATTEMPTS}): {e}"
                )
                current_model_num = (
                    current_model_num + 1 if current_model_num < max_model_num else 1
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
                print(
                    "\n\n\nRate limit reached, model is resting for 240 minutes\n\n\n"
                )
                time.sleep(14400)

            not_found_wording = "does not exist"
            if not_found_wording in polished_sentence:
                print("\n\n\nNot found. Skipped.\n\n\n")
                continue

            model_not_found_wording = "model_not_found"
            if not_found_wording in polished_sentence:
                print("\n\n\nModel not found. Skipped.\n\n\n")
                continue

            polished_sentences.append(polished_sentence)

    return polished_sentences
