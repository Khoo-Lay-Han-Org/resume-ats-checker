from pymilvus.client.global_topology import requests

from .vector_store import client
from .embedding_model import model

from playwright.sync_api import sync_playwright
from bs4 import BeautifulSoup


def scrape_content(user_input):
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        page.goto("https://duckduckgo.com/")
        page.fill("#searchbox_input", user_input),
        page.keyboard.press("Enter")
        page.wait_for_selector(".result", timeout=10000)

        results = page.query_selector_all(".result")

        links = []
        for result in results[:5]:
            title_element = result.query_selector(".result__title a")
            if title_element:
                link = title_element.get_attribute("href")
                links.append(link)

        all_content = []
        for link in links:
            response = requests.get(
                link,
                headers={
                    "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
                },
            )

            if response.status_code == 200:
                soup = BeautifulSoup(response.content, "html.parser")
                paragraphs = soup.find_all("p")
                text = " ".join(p.text for p in paragraphs)

                all_content.append(text)

        browser.close()

        return all_content


def embed_scraped_content(contents):
    data = []
    for content in contents:
        embedded_content = model.encode(content, truncate_dim=1000)
        data.append((embedded_content, content))

    return data


def insert_to_vector_store(contents):
    try:
        for content in contents:
            embedded_content, text = content
            client.insert("data", {"vector": embedded_content, "text": text})
        return True
    except:
        return False


def vector_search(user_input):
    data = model.encode(user_input, truncate_dim=1000)
    pass
