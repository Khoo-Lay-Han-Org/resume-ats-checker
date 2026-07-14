import time

from vector_store import client
from embedding_model import model
from orchestrate import *

import json
from urllib.parse import quote
from playwright.sync_api import sync_playwright


def _extract_text(result):
    for block in result["messages"][-1].content_blocks:
        if block.get("type") == "text":
            return block["text"]
    return result["messages"][-1].content_blocks[-1].get("text", "")


## LOOP 1 Before complete storing necessary data into vecctor DB
def slight_change_label_search(search_label):
    search_labels = []

    try:
        result = agent_word_ideater.invoke(
            {"messages": [{"role": "user", "content": search_label}]}
        )
    except Exception:
        print(
            f"\n\nFailed to generate search keywords for '{search_label}' (API error). Skipping."
        )
        return search_labels

    try:
        unparsed_array_keywords = _extract_text(result)
        search_keywords = json.loads(unparsed_array_keywords)

        for item in search_keywords:
            search_labels.append(item)
            print(
                f"\n\nAppending search keyword (slight change label search): \n{item}"
            )
    except Exception:
        print(f"\n\nFailed to parse array keywords (slight change label search)")

    return search_labels


def scrape_content(modified_search_label, max_pages=10):
    all_content = []

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page(
            user_agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"
        )

        for _, item in enumerate(modified_search_label):
            print(f"\n\nSearching with search keywords (scrape content): \n{item}")

            links = []

            for page_num in range(1, max_pages + 1):
                first = (page_num - 1) * 10 + 1
                url = f"https://www.bing.com/search?q={quote(item)}&first={first}"

                try:
                    page.goto(url)
                    page.wait_for_selector("#sb_form_q")
                except Exception:
                    try:
                        time.sleep(60)
                        page.goto(url)
                        page.wait_for_selector("#sb_form_q")
                    except Exception:
                        print(
                            f"\n\nFailed to search keyword page {page_num} (scrape content) {item}. Skipping page."
                        )
                        time.sleep(60)
                        continue

                a_tag = page.query_selector_all("h2 a")

                for tag in a_tag:
                    link = tag.get_attribute("href")
                    if link and link not in links:
                        links.append(link)

                print(
                    f"\n\nFound {len(links)} unique links so far (page {page_num}/{max_pages})"
                )

                time.sleep(2)

            for link in links:
                print(f"\n\nSearching with link (scrape content): \n{link}")

                try:
                    page.goto(link, wait_until="domcontentloaded", timeout=15000)
                    page.wait_for_load_state("networkidle", timeout=10000)
                except Exception:
                    print(f"\n\nCannot access (scrape content): \n{link}")
                    continue

                try:
                    paragraphs = page.locator("p").all_inner_texts()
                    valid_paragraphs = []
                    for item in paragraphs:
                        if len(item) > 30:
                            valid_paragraphs.append(item)
                            print(
                                f"\n\nAppending individual sentence as valid sentence (scrape content): \n{item}"
                            )
                    all_content.append(valid_paragraphs)
                    print(
                        f"\n\nAppending the entire valid sentences of the page (scrape content): \n{item}"
                    )
                except Exception:
                    print(f"\n\nFailed to extract (scrape content): \n{link}")
                    pass

        browser.close()

    return all_content


def polish_scraped_content(scraped_content):
    polished_content = []

    for item_arr in scraped_content:
        for item in item_arr:
            item = item.strip()

            if len(item) < 30:
                continue

            print(
                f"\n\nPolishing the following sentence (polish scraped content): \n{item}"
            )

            try:
                result = agent_sentence_polisher.invoke(
                    {"messages": [{"role": "user", "content": item}]}
                )

                try:
                    text = _extract_text(result)
                    polished_content.append(text)

                    print(
                        f"\n\nAppending the following polished sentence (polish scraped content): \n{text}"
                    )
                except Exception:
                    pass
            except Exception:
                time.sleep(60)

                try:
                    result = agent_sentence_polisher.invoke(
                        {"messages": [{"role": "user", "content": item}]}
                    )

                    try:
                        text = _extract_text(result)
                        polished_content.append(text)

                        print(
                            f"\n\nAppending the following polished sentence (polish scraped content): \n{text}"
                        )
                    except Exception:
                        print(
                            f"\n\nFailed to polish the following sentence (polish scraped content): \n{item}"
                        )
                        pass

                except Exception:
                    print(f"\n\nModel Needs Rest (polish scraped content)")
                    time.sleep(60)

    return polished_content


def embed_scraped_content(contents):
    data = []
    for content in contents:
        print(
            f"\n\nEmbedding the following sentence (embed scraped content): \n{content}"
        )

        embedded_content = model.encode(
            content, truncate_dim=50, normalize_embeddings=True
        )
        data.append((embedded_content, content))

        print(
            f"\n\nEmbedded the following sentence (embed scraped content): \n{embedded_content}"
        )

    return data


def insert_to_vector_store(contents, deduplicate=True):
    seen_texts = set()

    try:
        for content in contents:
            embedded_content, text = content

            if deduplicate:
                if text in seen_texts:
                    print(
                        f"\n\nSkipping duplicate (seen in this batch): {text[:60]}..."
                    )
                    continue

                existing = client.query(
                    collection_name="data",
                    filter=f'text == "{text}"',
                    output_fields=["id"],
                    limit=1,
                )
                if existing:
                    print(f"\n\nSkipping duplicate (already in DB): {text[:60]}...")
                    continue

            seen_texts.add(text)
            print(
                f"\n\nStoring the following embedded data (insert to vector store): \n{content}"
            )
            client.insert("data", {"vector": embedded_content, "text": text})
            print(
                f"\n\nStored the following embedded data (insert to vector store): \n{content}"
            )
        return True
    except:
        return False


def deduplicate_vector_store():
    print("\n\nDeduplicating vector store...")

    all_entries = client.query(
        collection_name="data",
        output_fields=["id", "text"],
        limit=100000,
    )

    seen_texts = set()
    ids_to_delete = []

    for entry in all_entries:
        text = entry["text"]
        if text in seen_texts:
            ids_to_delete.append(entry["id"])
        else:
            seen_texts.add(text)

    if ids_to_delete:
        client.delete(collection_name="data", ids=ids_to_delete)
        print(f"\n\nDeleted {len(ids_to_delete)} duplicate entries")
    else:
        print("\n\nNo duplicates found")

    return len(ids_to_delete)


## LOOP 1


## LOOP 2 After already have necessary data in vector db (for dataset)
def vector_search(label, confidence_threshold=0.5):
    print(f"\n\nSearching {label} from vector store (threshold={confidence_threshold})")

    data = model.encode(label, truncate_dim=50, normalize_embeddings=True)

    search_params = {"metric_type": "IP", "params": {}}
    if confidence_threshold is not None:
        search_params["params"]["radius"] = confidence_threshold

    search_results = client.search(
        collection_name="data",
        data=[data],
        limit=300,
        output_fields=["text"],
        search_params=search_params,
    )

    closest_sentences = []
    for result in search_results[0]:
        print(f"\n\nClosest sentence (distance={result['distance']:.4f}): \n{result}")
        closest_sentences.append(result["entity"]["text"])

    return closest_sentences


## LOOP 2
