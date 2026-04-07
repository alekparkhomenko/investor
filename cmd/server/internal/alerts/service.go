package alerts

import (
	"investor-notifications/pkg/models"
	"log"
	"sync"
	"time"
)

type Service struct {
	alerts []models.Alert
	mu     sync.RWMutex
}

func NewService() *Service {
	return &Service{
		alerts: make([]models.Alert, 0),
	}
}

func (s *Service) Add(alert models.Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts = append(s.alerts, alert)
	log.Printf("Added alert: %s %s %.2f", alert.Symbol, alert.Condition, alert.TargetPrice)
}

func (s *Service) Check(quote models.Quote) []models.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var triggered []models.Alert

	for i := range s.alerts {
		alert := &s.alerts[i]

		if alert.Symbol != quote.Symbol {
			continue
		}

		if alert.Triggered {
			continue
		}

		shouldTrigger := false
		if alert.Condition == "above" && quote.Price >= alert.TargetPrice {
			shouldTrigger = true
		} else if alert.Condition == "below" && quote.Price <= alert.TargetPrice {
			shouldTrigger = true
		}

		if shouldTrigger {
			alert.Triggered = true
			alert.CreatedAt = time.Now().Unix()
			triggered = append(triggered, *alert)
			log.Printf("Alert triggered: %s %s %.2f at %.2f", alert.Symbol, alert.Condition, alert.TargetPrice, quote.Price)
		}
	}

	return triggered
}

func (s *Service) GetAll() []models.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]models.Alert, len(s.alerts))
	copy(result, s.alerts)
	return result
}
