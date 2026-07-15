package ats_util

import (
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

// Real user-agent rotation pool.
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:127.0) Gecko/20100101 Firefox/127.0",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
}

func randomUA() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func randomDelay(minMs, maxMs int) {
	d := time.Duration(rand.Intn(maxMs-minMs+1)+minMs) * time.Millisecond
	time.Sleep(d)
}

// humanType types a string character-by-character with realistic pauses.
func humanType(el *rod.Element, text string) {
	for _, ch := range text {
		el.MustInput(string(ch))
		randomDelay(40, 120)
	}
}

// stealthHarden patches browser fingerprint properties.
func stealthHarden(page *rod.Page) {
	page.MustEval(`
		Object.defineProperty(navigator, 'webdriver', { get: () => false });
		Object.defineProperty(navigator, 'plugins', { get: () => [1, 2, 3, 4, 5] });
		Object.defineProperty(navigator, 'languages', { get: () => ['en-US', 'en'] });
		window.chrome = { runtime: {} };
	`)
}

// setCommonHeaders mimics a real browser request header set.
func setCommonHeaders(page *rod.Page) {
	page.SetExtraHeaders([]string{
		"Accept-Language", "en-US,en;q=0.9",
		"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
		"DNT", "1",
		"Connection", "keep-alive",
		"Upgrade-Insecure-Requests", "1",
		"Sec-Fetch-Dest", "document",
		"Sec-Fetch-Mode", "navigate",
		"Sec-Fetch-Site", "none",
		"Sec-Fetch-User", "?1",
	})
}

func JobDescWebScrape(company, jobTitle string) string {
	searchString := fmt.Sprintf("New jobs from %s with %s role", company, jobTitle)

	// Launcher with automation-evasion flags
	l := launcher.New().
		Headless(true).
		Set("disable-blink-features=AutomationControlled").
		Set("disable-automation").
		Set("no-sandbox").
		Set("disable-setuid-sandbox").
		Set("disable-infobars").
		Set("window-size=1920,1080").
		Set("start-maximized")

	// Uncomment to route through a residential proxy:
	// l = l.Proxy("http://user:pass@host:port")

	browser := rod.New().ControlURL(l.MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	setCommonHeaders(page)
	stealthHarden(page)

	// Human-like pre-navigation pause
	randomDelay(1000, 3000)

	page.MustNavigate("https://www.google.com")
	page.MustWaitLoad()

	randomDelay(500, 1500)

	// Find search box and type naturally
	searchBox := page.MustElement("textarea[name='q']")
	humanType(searchBox, searchString)

	randomDelay(200, 500)

	// Press Enter
	searchBox.MustKeyActions().Press(input.Enter)
	page.MustWaitLoad()

	randomDelay(2000, 4000)

	// Extract result links
	results := page.MustElements("div.g")

	links := []string{}
	for i, result := range results {
		if len(links) >= 5 || i > 50 {
			break
		}

		linkElem, err := result.Element("a")
		if err != nil {
			continue
		}

		href, err := linkElem.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		if u, err := url.ParseRequestURI(*href); err == nil {
			links = append(links, u.String())
		}
	}

	var allText strings.Builder

	for _, item := range links {
		randomDelay(2000, 5000)

		directPage := browser.MustPage()
		setCommonHeaders(directPage)
		stealthHarden(directPage)

		// Set referer so it looks like the user clicked through from Google
		directPage.SetExtraHeaders([]string{
			"Referer", "https://www.google.com/",
		})

		directPage.MustNavigate(item)
		directPage.MustWaitLoad()

		randomDelay(1000, 2000)

		// Scroll like a real reader
		directPage.MustEval(`window.scrollTo({ top: document.body.scrollHeight * 0.3, behavior: 'smooth' })`)
		randomDelay(500, 1200)
		directPage.MustEval(`window.scrollTo({ top: document.body.scrollHeight * 0.6, behavior: 'smooth' })`)
		randomDelay(300, 800)

		bodyText := directPage.MustElement("body").MustText()
		allText.WriteString(bodyText)
		allText.WriteString("\n")
	}

	return allText.String()
}
