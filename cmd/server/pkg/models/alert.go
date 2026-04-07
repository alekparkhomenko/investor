package models

type Quote struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type Alert struct {
	ID          string  `json:"id"`
	UserID      int64   `json:"user_id"`
	Symbol      string  `json:"symbol"`
	TargetPrice float64 `json:"target_price"`
	Condition   string  `json:"condition"` // "above" or "below"
	Triggered   bool    `json:"triggered"`
	CreatedAt   int64   `json:"created_at"`
}

type AlertTriggered struct {
	AlertID     string  `json:"alert_id"`
	Symbol      string  `json:"symbol"`
	Price       float64 `json:"price"`
	TargetPrice float64 `json:"target_price"`
	Condition   string  `json:"condition"`
	Timestamp   int64   `json:"timestamp"`
}
