import json
import random
import time
from pathlib import Path

from typing import Any, cast

from playwright.sync_api import BrowserContext, BrowserType, Page, Playwright, ViewportSize


_REAL_USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
]

_VIEWPORTS: list[ViewportSize] = [
    {"width": 1920, "height": 1080},
    {"width": 1366, "height": 768},
    {"width": 1536, "height": 864},
    {"width": 1440, "height": 900},
    {"width": 1280, "height": 720},
]

_LOCALES = ["en-US", "en-GB", "en-CA", "en-AU"]
_TIMEZONES = [
    "America/New_York",
    "America/Chicago",
    "America/Los_Angeles",
    "America/Denver",
    "Europe/London",
    "Europe/Berlin",
]

SESSION_DIR = Path(__file__).parent / "browser_sessions"


def random_delay(min_sec=1.0, max_sec=3.0):
    time.sleep(random.uniform(min_sec, max_sec))


def debug_screenshot(page: Page, label: str = "debug"):
    """Save a screenshot for debugging. Only active if DEBUG_SCREENSHOTS env var is set."""
    import os
    if os.environ.get("DEBUG_SCREENSHOTS"):
        path = SESSION_DIR / f"{label}_{int(time.time())}.png"
        page.screenshot(path=str(path))
        print(f"  [debug] Screenshot saved: {path}")


def _context_dir(proxy_label=None):
    label = proxy_label or "default"
    p = SESSION_DIR / label
    p.mkdir(parents=True, exist_ok=True)
    return str(p)


def _load_cookies(context: BrowserContext, proxy_label=None):
    path = Path(_context_dir(proxy_label)) / "cookies.json"
    if path.exists():
        with open(path) as f:
            context.add_cookies(json.load(f))


def _save_cookies(context: BrowserContext, proxy_label=None):
    path = Path(_context_dir(proxy_label)) / "cookies.json"
    path.parent.mkdir(parents=True, exist_ok=True)
    with open(path, "w") as f:
        json.dump(context.cookies(), f, indent=2)


def create_stealth_browser_context(
    p: Playwright,
    *,
    headless: bool = True,
    proxy: str | None = None,
    user_agent: str | None = None,
    session_label: str | None = None,
    slow_mo: int | None = None,
    use_firefox: bool = False,
) -> tuple[BrowserContext, str]:
    """Launch a stealth browser context that mimics a real user.

    Returns (context, used_user_agent). Set use_firefox=True for Bing
    (Firefox is detected much less than Chromium).
    """
    ua = user_agent or random.choice(_REAL_USER_AGENTS)
    viewport = random.choice(_VIEWPORTS)

    proxy_opts = None
    if proxy:
        proxy_opts = {"server": proxy}

    browser_type: BrowserType
    if use_firefox:
        browser_type = p.firefox
    else:
        browser_type = p.chromium

    extra_args: list[str] = [
        "--disable-features=AllowPopupsDuringPageHide",
        "--disable-ipc-flooding-protection",
        "--no-default-browser-check",
    ]

    browser = browser_type.launch(
        headless=headless,
        slow_mo=slow_mo,
        args=extra_args if not use_firefox else None,
        firefox_user_prefs={
            "dom.webdriver.enabled": False,
            "dom.webnotifications.enabled": False,
            "media.autoplay.enabled": False,
        } if use_firefox else None,
    )

    headers: dict[str, str] = {
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.9",
        "Accept-Encoding": "gzip, deflate, br",
        "Sec-Fetch-Dest": "document",
        "Sec-Fetch-Mode": "navigate",
        "Sec-Fetch-Site": "none",
        "Sec-Fetch-User": "?1",
        "Upgrade-Insecure-Requests": "1",
        "DNT": "1",
        "Connection": "keep-alive",
    }

    if not use_firefox:
        sec_ch_version = "125"
        try:
            sec_ch_version = ua.split("Chrome/")[-1].split(" ")[0]
        except Exception:
            pass
        headers["Sec-CH-UA"] = f'"Not/A)Brand";v="99", "Google Chrome";v="{sec_ch_version}"'
        headers["Sec-CH-UA-Mobile"] = "?0"
        sec_ch_platform = '"Windows"'
        if "Mac" in ua and "OS X" in ua:
            sec_ch_platform = '"macOS"'
        elif "Linux" in ua:
            sec_ch_platform = '"Linux"'
        headers["Sec-CH-UA-Platform"] = sec_ch_platform

    context = browser.new_context(
        user_agent=ua,
        viewport=viewport,
        locale=random.choice(_LOCALES),
        timezone_id=random.choice(_TIMEZONES),
        color_scheme="light",
        proxy=cast("Any", proxy_opts),
        extra_http_headers=headers,
    )

    if not use_firefox:
        context.add_init_script("""
            Object.defineProperty(navigator, 'webdriver', { get: () => false });
            Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
            Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
            window.chrome = { runtime: {} };
        """)

    _load_cookies(context, session_label)

    return context, ua


def close_stealth_context(context: BrowserContext, session_label: str | None = None):
    _save_cookies(context, session_label)
    context.close()
    br = context.browser
    if br:
        br.close()


def human_like_scroll(page: Page, steps: int = 3):
    """Scroll the page in random increments like a human reading."""
    for _ in range(steps):
        delta = random.randint(200, 600)
        page.evaluate(f"window.scrollBy(0, {delta})")
        random_delay(0.3, 1.2)


def human_like_mouse_move(page: Page, target_selector: str):
    """Move the mouse to a target element using a natural bezier-like path."""
    try:
        el = page.query_selector(target_selector)
        if not el:
            return
        box = el.bounding_box()
        if not box:
            return
        target_x = box["x"] + box["width"] / 2
        target_y = box["y"] + box["height"] / 2

        vp = page.viewport_size
        if not vp:
            return
        width = vp.get("width", 1920)
        height = vp.get("height", 1080)
        start_x = random.randint(0, width)
        start_y = random.randint(0, height)

        steps_count = random.randint(8, 16)
        for i in range(steps_count):
            t = (i + 1) / steps_count
            # Simple bezier-like interpolation with overshoot
            mid_x = start_x + (target_x - start_x) * t
            mid_y = start_y + (target_y - start_y) * t
            jitter_x = random.uniform(-5, 5)
            jitter_y = random.uniform(-3, 3)
            page.mouse.move(mid_x + jitter_x, mid_y + jitter_y)
            time.sleep(random.uniform(0.01, 0.03))

        page.mouse.move(target_x, target_y)
    except Exception:
        pass


def stealth_search(page: Page, query: str):
    """Type a search query character-by-character with human-like delays."""
    for char in query:
        page.keyboard.type(char, delay=random.randint(40, 120))


def create_stealth_browser(
    p: Playwright,
    *,
    headless: bool = True,
    proxy: str | None = None,
    user_agent: str | None = None,
    session_label: str | None = None,
):
    """Convenience wrapper returning (context, page, used_user_agent)."""
    context, ua = create_stealth_browser_context(
        p,
        headless=headless,
        proxy=proxy,
        user_agent=user_agent,
        session_label=session_label,
    )
    page = context.new_page()
    return context, page, ua
