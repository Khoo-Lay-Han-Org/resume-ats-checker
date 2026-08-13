import os
from urllib.parse import urlparse

from bs4 import BeautifulSoup
from dotenv import find_dotenv, load_dotenv
from requests import RequestException, get

from ..config.scraping import MAX_CONTENT_ITEMS, MAX_LINKS, REQUEST_TIMEOUT
from ..targets.content_elements import CONTENT_ELEMENTS, RESULT_LINK_SELECTOR

load_dotenv(find_dotenv())


def scrape_links(query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"):
    print(f"Getting links for query: {query}")
    params = {"q": query, "categories": search_type, "limit": 20}
    response = get(f"{scraper_url}/search", params=params, timeout=REQUEST_TIMEOUT)
    soup = BeautifulSoup(response.text, "html.parser")

    scraper_host = str(urlparse(scraper_url or "").netloc)

    all_links = []
    for link in soup.select(RESULT_LINK_SELECTOR):
        href = link.get("href")

        if isinstance(href, list):
            href = href[0] if href else ""
        else:
            href = str(href) if href else ""

        if not href or href.startswith("#") or href.startswith("javascript:"):
            continue

        if not href.startswith("http"):
            continue

        if scraper_host in href:
            continue

        if href not in all_links:
            all_links.append(href)

    all_links = all_links[:MAX_LINKS]
    print(f"Found {len(all_links)} unique result links")
    return all_links


def scrape_dataset_contents(
    query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"
):
    all_links = scrape_links(query, scraper_url, search_type)

    all_contents = []

    for index, link in enumerate(all_links, start=1):
        try:
            response = get(link, timeout=REQUEST_TIMEOUT)
        except RequestException as e:
            print(f"  Failed to fetch link: {e}")
            continue

        if response.status_code != 200:
            print(f"  Failed to fetch link: {response.status_code}")
            continue

        soup = BeautifulSoup(response.text, "html.parser")

        page_contents = []

        for target in CONTENT_ELEMENTS:
            for content in soup.find_all(target):
                text = content.get_text(strip=True)

                if len(text) >= 100:
                    page_contents.append(text)

        print(f"  Got {len(page_contents)} content item(s) from {link}")
        for item in page_contents:
            print(f"    - {item[:200]}")

        all_contents.extend(page_contents)

        if len(all_contents) >= MAX_CONTENT_ITEMS:
            print(f"  Reached cap of {MAX_CONTENT_ITEMS} content items, stopping")
            break

    return all_contents[:MAX_CONTENT_ITEMS]
