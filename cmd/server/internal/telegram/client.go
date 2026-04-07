package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"investor-notifications/pkg/models"
	"log"
	"net/http"
)

type Client struct {
	token   string
	chatID  int64
	baseURL string
}

func NewClient(token string, chatID int64) *Client {
	return &Client{
		token:   token,
		chatID:  chatID,
		baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
	}
}

func (c *Client) SendAlert(alert models.Alert, price float64) error {
	text := fmt.Sprintf("🔔 Alert Triggered!\n\n%s %s %.2f\nCurrent price: %.2f",
		alert.Symbol, alert.Condition, alert.TargetPrice, price)

	return c.SendMessage(text)
}

func (c *Client) SendMessage(text string) error {
	if c.token == "" || c.chatID == 0 {
		log.Println("Telegram not configured, skipping message")
		return nil
	}

	msg := map[string]interface{}{
		"chat_id": c.chatID,
		"text":    text,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	url := c.baseURL + "/sendMessage"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("telegram API returned %d", resp.StatusCode)
	}

	log.Printf("Telegram message sent: %s", text)
	return nil
}
