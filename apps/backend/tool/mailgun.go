package tool

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	systemconfig "resuming/system-config"
)

func SendEmail(to, subject, body string, isHTML bool) error {
	form := url.Values{
		"from":    {fmt.Sprintf("Resuming <noreply@%s>", systemconfig.MailgunDomain)},
		"to":      {to},
		"subject": {subject},
	}
	if isHTML {
		form.Set("html", body)
	} else {
		form.Set("text", body)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", systemconfig.MailgunDomain),
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create Mailgun request: %w", err)
	}

	req.SetBasicAuth("api", systemconfig.MailgunAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send via Mailgun: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("mailgun API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
