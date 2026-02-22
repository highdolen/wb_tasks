package notification

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Sender interface {
	Send(ctx context.Context, chatID string, message string) error
}

type Telegram struct {
	botToken string
	client   *http.Client
}

// NewTelegram - конструктор Telegram-соединения
func NewTelegram(botToken string) *Telegram {
	return &Telegram{
		botToken: botToken,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Send - отправка уведомления в телеграм
func (t *Telegram) Send(ctx context.Context, chatID string, message string) error {
	if t.botToken == "" || chatID == "" {
		return fmt.Errorf("telegram bot token or chat ID not set")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.botToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL,
		strings.NewReader(url.Values{
			"chat_id": {chatID},
			"text":    {message},
		}.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram send failed: %s", resp.Status)
	}

	return nil
}
