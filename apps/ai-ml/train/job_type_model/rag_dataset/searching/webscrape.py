import os
from bs4 import BeautifulSoup

from dotenv import load_dotenv, find_dotenv
from pymilvus.client.global_topology import requests

from ..targets.content_elements import *

load_dotenv(find_dotenv())


def scrape_links(query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"):
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

    return all_links


def scrape_dataset_contents(
    query, scraper_url=os.getenv("SEARXNG_URI"), search_type="general"
):
    params = {"q": query, "categories": search_type, "limit": 100}
    response = requests.get(f"{scraper_url}/search", params=params)
    soup = BeautifulSoup(response.text, "html.parser")

    all_contents = []

    for target in CONTENT_ELEMENTS:
        for content in soup.find_all(target):
            text = content.get_text(strip=True)

            if len(text) >= 100:
                all_contents.append(text)

    return all_contents
