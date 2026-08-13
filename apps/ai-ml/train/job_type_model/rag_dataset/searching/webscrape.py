import os
from bs4 import BeautifulSoup

from dotenv import load_dotenv, find_dotenv
from pymilvus.client.global_topology import requests

from ..targets.content_elements import *

load_dotenv(find_dotenv())


def scrape_links(query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"):
    print(f"Getting links for query: {query}")
    params = {"q": query, "categories": search_type, "limit": 20}
    response = requests.get(f"{scraper_url}/search", params=params)
    soup = BeautifulSoup(response.text, "html.parser")

    all_links = []
    for link in soup.find_all(LINK_ELEMENT, href=True):
        href = link.get("href")

        if isinstance(href, list):
            href = href[0] if href else ""
        else:
            href = str(href) if href else ""

        if not href or href.startswith("#") or href.startswith("javascript:"):
            continue

        if href.startswith("http"):
            all_links.append(href)
        else:
            all_links.append(f"{scraper_url}{href}")

    print(f"Found {len(all_links)} links")
    return all_links


def scrape_dataset_contents(
    query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"
):
    all_links = scrape_links(query, scraper_url, search_type)

    all_contents = []

    for index, link in enumerate(all_links, start=1):
        print(f"Searching link {index}/{len(all_links)}: {link}")
        response = requests.get(link)

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

    return all_contents
