package ats_util

import (
	"net/url"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
)

func JobDescWebScrape(company, job_title string) string {
	search_string := "New jobs from " + company + " with " + job_title + " role"

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://www.google.com")
	page.MustWaitLoad()
	search_box := page.MustElement("textarea[name='q']")
	search_box.MustInput(search_string)
	search_box.MustKeyActions().Press(input.Enter)
	page.MustWaitLoad()

	results := page.MustElements("div.p")

	links := []string{}
	for i, result := range results {
		if len(links) >= 5 || i > 50 {
			break
		}

		link_elem, err := result.Element("a")
		if err != nil {
			continue
		}

		href, err := link_elem.Attribute("href")
		if err != nil || href == nil {
			continue
		}

		if url, err := url.ParseRequestURI(*href); err == nil {
			links = append(links, url.String())
		}
	}

	var all_text strings.Builder
	for _, item := range links {
		direct_page := browser.MustPage(item)
		// it does not extract HTML tags
		body_text := direct_page.MustElement("body").MustText()
		all_text.WriteString(body_text)
		all_text.WriteString("\n")
	}

	return all_text.String()
}
