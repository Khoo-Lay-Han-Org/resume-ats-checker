import json
import random
import traceback
from urllib.parse import quote

import requests

from playwright.sync_api import sync_playwright

from vector_store import client
from embedding_model import model
from orchestrate import *
from browser_stealth import (
    close_stealth_context,
    create_stealth_browser_context,
    human_like_scroll,
    random_delay,
)


# ── Search backend ──────────────────────────────────────────────────
# Set to a SearXNG instance URL to bypass Bing's bot detection entirely.
# You can self-host with Docker or use a public instance (e.g. https://searx.be).
# Leave empty to fall back to Bing + Playwright.
SEARXNG_INSTANCE = "http://localhost:8080"

# ── Proxy rotation list ─────────────────────────────────────────────
# Format: "http://user:pass@host:port" or "socks5://user:pass@host:port"
PROXIES: list[str] = []

# How often (per-N-search-keywords) to rotate proxy + start a fresh session
PROXY_ROTATE_EVERY = 3

# Bing search-result selectors (most to least specific)
_BING_RESULT_SELECTORS = [
    "li.b_algo h2 a",
    "#b_results > li h2 a",
    "ol#b_results h2 a",
    "h2 a",
]

# Skip links that are tracking / non-content URLs
_SKIP_LINK_PATTERNS = (
    "bing.com/ck/a",
    "bing.com/search",
    "go.microsoft.com",
    "/news/",
    "/videos/",
)

MAX_RETRIES = 3

HEADLESS = True
USE_FIREFOX = True


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


def _search_searxng_api(query: str, max_results: int = 20) -> list[str]:
    """Search via SearXNG JSON API (fast path, but may be blocked by limiter)."""
    url = f"{SEARXNG_INSTANCE.rstrip('/')}/search"
    params = {
        "q": query,
        "format": "json",
        "language": "en-US",
        "categories": "general",
        "pageno": 1,
    }

    headers = {
        "User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "DNT": "1",
    }

    resp = requests.get(url, params=params, headers=headers, timeout=30)
    resp.raise_for_status()
    data = resp.json()
    results = data.get("results", [])
    links: list[str] = []
    for r in results:
        href = r.get("url", "")
        if href and not any(p in href for p in _SKIP_LINK_PATTERNS):
            links.append(href)
            if len(links) >= max_results:
                break
    return links


def _search_searxng_browser(context, query: str, max_results: int = 20) -> list[str]:
    """Fallback: scrape SearXNG HTML results via Playwright (works even with limiter on)."""
    page = context.new_page()
    url = f"{SEARXNG_INSTANCE.rstrip('/')}/search?q={quote(query)}&language=en-US&categories=general"
    links: list[str] = []

    try:
        page.goto(url, wait_until="domcontentloaded", timeout=30000)
        random_delay(1.0, 2.0)

        for selector in ["article.result a[href]", ".result a.url", "a[href]"]:
            try:
                elements = page.query_selector_all(selector)
                for el in elements:
                    href = el.get_attribute("href")
                    if (
                        href
                        and href not in links
                        and href.startswith("http")
                        and not any(p in href for p in _SKIP_LINK_PATTERNS)
                    ):
                        links.append(href)
                        if len(links) >= max_results:
                            break
                if links:
                    break
            except Exception:
                continue
    except Exception as e:
        print(f"  SearXNG browser search failed: {e}")
    finally:
        page.close()

    return links


def _search_searxng(context, query: str, max_results: int = 20) -> list[str]:
    """Search via SearXNG. Tries JSON API first, falls back to browser scraping."""
    if not SEARXNG_INSTANCE:
        return []

    try:
        links = _search_searxng_api(query, max_results)
        if links:
            print(f"  SearXNG API returned {len(links)} links")
            return links
    except requests.RequestException as e:
        status = e.response.status_code if e.response is not None else "?"
        print(f"  SearXNG API failed (HTTP {status}), trying browser fallback...")
    except Exception:
        print(f"  SearXNG API failed, trying browser fallback...")

    print(f"  Falling back to SearXNG browser scrape...")
    return _search_searxng_browser(context, query, max_results)


def _extract_bing_links(page, existing_links):
    """Try multiple selectors to extract result links from a Bing SERP."""
    links = list(existing_links)
    for selector in _BING_RESULT_SELECTORS:
        try:
            elements = page.query_selector_all(selector)
            if elements:
                for el in elements:
                    href = el.get_attribute("href")
                    if (
                        href
                        and href not in links
                        and not any(p in href for p in _SKIP_LINK_PATTERNS)
                    ):
                        links.append(href)
                if links:
                    break
        except Exception:
            continue
    return links


def _scrape_page_text(page, link: str) -> list[str]:
    """Navigate to a URL and extract long paragraphs from it."""
    try:
        page.goto(link, wait_until="domcontentloaded", timeout=15000)
        page.wait_for_load_state("networkidle", timeout=10000)
    except Exception:
        print(f"\n\nCannot access: \n{link}")
        return []

    random_delay(0.5, 1.5)
    human_like_scroll(page, steps=random.randint(2, 4))

    try:
        paragraphs = page.locator("p").all_inner_texts()
        valid = [p.strip() for p in paragraphs if len(p.strip()) > 30]
        if valid:
            print(f"\n\nExtracted {len(valid)} paragraphs from: \n{link[:80]}...")
        else:
            print(f"\n\nNo long paragraphs on: \n{link[:80]}...")
        return valid
    except Exception:
        print(f"\n\nFailed to extract from: \n{link[:80]}...")
        return []


def scrape_content(modified_search_label, max_pages=10):
    all_content = []

    with sync_playwright() as p:
        context = None
        ua = "Mozilla/5.0"
        proxy = None
        session_index = 0

        for kw_index, item in enumerate(modified_search_label):
            if context is None or (kw_index > 0 and kw_index % PROXY_ROTATE_EVERY == 0):
                if context is not None:
                    close_stealth_context(context, session_label=f"bing_{session_index}")
                    random_delay(5, 10)

                session_index += 1
                proxy = random.choice(PROXIES) if PROXIES else None
                context, ua = create_stealth_browser_context(
                    p,
                    headless=HEADLESS,
                    proxy=proxy,
                    session_label=f"bing_{session_index}",
                    use_firefox=USE_FIREFOX,
                )

            print(f"\n\nSearching with search keywords (scrape content): \n{item}")
            print(f"  Proxy: {proxy or 'none'}  |  UA: {ua[:60]}...")

            # ── Get search result links ──────────────────────────────
            links: list[str] = []

            if SEARXNG_INSTANCE:
                links = _search_searxng(context, item, max_results=30)
            else:
                # Fallback: scrape Bing via Playwright
                page = context.new_page()
                for page_num in range(1, max_pages + 1):
                    first = (page_num - 1) * 10 + 1
                    url = f"https://www.bing.com/search?q={quote(item)}&first={first}"

                    random_delay(1.5, 4.0)

                    loaded = False
                    for attempt in range(MAX_RETRIES):
                        try:
                            page.goto(url, wait_until="domcontentloaded", timeout=30000)
                            random_delay(0.5, 1.5)
                            page.wait_for_selector("#b_results", timeout=20000)
                            loaded = True
                            break
                        except Exception:
                            print(f"  Bing attempt {attempt + 1} failed for page {page_num}")
                            if attempt < MAX_RETRIES - 1:
                                random_delay(20, 40)
                                page.close()
                                page = context.new_page()

                    if not loaded:
                        print(f"  Skipping Bing page {page_num} for '{item}'.")
                        random_delay(20, 40)
                        continue

                    page_links = _extract_bing_links(page, links)
                    if page_links:
                        links = page_links
                    print(f"  Found {len(links)} unique links so far (page {page_num}/{max_pages})")
                    random_delay(1.0, 3.0)

                page.close()

            # ── Scrape each result page ──────────────────────────────
            page = None
            for link in links:
                if page is None:
                    page = context.new_page()
                print(f"\n\nScraping content from: \n{link}")
                random_delay(2.0, 5.0)

                paragraphs = _scrape_page_text(page, link)
                if paragraphs:
                    all_content.append(paragraphs)

            if page is not None:
                page.close()

        if context is not None:
            close_stealth_context(context, session_label=f"bing_{session_index}")

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

            random_delay(0.5, 2.0)

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
                random_delay(30, 60)

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
                    random_delay(30, 60)

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

    for content in contents:
        embedded_content, text = content

        if deduplicate:
            if text in seen_texts:
                print(
                    f"\n\nSkipping duplicate (seen in this batch): {text[:60]}..."
                )
                continue

            escaped = text.replace('"', '\\"')
            try:
                existing = client.query(
                    collection_name="data",
                    filter=f'text == "{escaped}"',
                    output_fields=["id"],
                    limit=1,
                )
                if existing:
                    print(f"\n\nSkipping duplicate (already in DB): {text[:60]}...")
                    continue
            except Exception:
                print(f"\n\nDedup query failed for text, inserting anyway: {text[:60]}...")

        seen_texts.add(text)
        print(
            f"\n\nStoring the following embedded data (insert to vector store): \n{content}"
        )
        try:
            client.insert("data", {"vector": embedded_content, "text": text})
            print(
                f"\n\nStored the following embedded data (insert to vector store): \n{content}"
            )
        except Exception as e:
            print(f"\n\nFailed to insert: {e}\n{traceback.format_exc()}")
            return False

    return True


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
