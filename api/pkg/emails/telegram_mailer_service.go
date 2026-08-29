package emails

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/NdoleStudio/stacktrace"
)

const telegramMessageMaxRunes = 4096

// TelegramConfig configures delivery of notification copies to a Telegram chat.
type TelegramConfig struct {
	BotToken string
	ChatID   string
}

type telegramMailer struct {
	botToken string
	chatID   string
	client   *http.Client
}

// NewTelegramMailer creates a mailer that sends a plain-text copy of each email
// notification to the configured Telegram chat.
func NewTelegramMailer(config TelegramConfig) Mailer {
	return &telegramMailer{
		botToken: config.BotToken,
		chatID:   config.ChatID,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (mailer *telegramMailer) Send(ctx context.Context, email *Email) error {
	text := strings.TrimSpace(fmt.Sprintf("%s\n\n%s", email.Subject, email.Text))
	if len([]rune(text)) > telegramMessageMaxRunes {
		text = string([]rune(text)[:telegramMessageMaxRunes-1]) + "…"
	}

	payload, err := json.Marshal(struct {
		ChatID string `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: mailer.chatID,
		Text:   text,
	})
	if err != nil {
		return stacktrace.Propagatef(err, "cannot encode Telegram notification")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", mailer.botToken), bytes.NewReader(payload))
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create Telegram notification request")
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := mailer.client.Do(request)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot send Telegram notification")
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1024))
		return fmt.Errorf("Telegram notification failed with status %s: %s", response.Status, strings.TrimSpace(string(body)))
	}

	return nil
}
