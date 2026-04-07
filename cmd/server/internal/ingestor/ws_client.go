package ingestor

import (
	"encoding/json"
	"fmt"
	"investor-notifications/pkg/models"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	conn    *websocket.Conn
	apiKey  string
	symbols []string
}

func NewWSClient(apiKey string, symbols []string) *WSClient {
	return &WSClient{
		apiKey:  apiKey,
		symbols: symbols,
	}
}

func (c *WSClient) Connect() error {
	url := fmt.Sprintf("wss://ws.finnhub.io?token=%s", c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	for _, symbol := range c.symbols {
		msg := map[string]string{"type": "subscribe", "symbol": symbol}
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error subscribing to %s: %v", symbol, err)
		}
	}

	return nil
}

func (c *WSClient) ReadQuotes() <-chan models.Quote {
	quotes := make(chan models.Quote, 100)

	go func() {
		defer close(quotes)
		for {
			if c.conn == nil {
				time.Sleep(5 * time.Second)
				if err := c.Connect(); err != nil {
					log.Printf("Reconnect error: %v", err)
					continue
				}
			}

			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				c.conn = nil
				continue
			}

			var data map[string]interface{}
			if err := json.Unmarshal(msg, &data); err != nil {
				continue
			}

			if data["type"] == "trade" {
				if trades, ok := data["data"].([]interface{}); ok {
					for _, t := range trades {
						trade := t.(map[string]interface{})
						quote := models.Quote{
							Symbol:    trade["s"].(string),
							Price:     trade["p"].(float64),
							Timestamp: int64(trade["t"].(float64)),
						}
						quotes <- quote
					}
				}
			}
		}
	}()

	return quotes
}
