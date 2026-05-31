package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/investor/internal/portfolio"
	"github.com/alekparkhomenko/investor/platform/pkg/logger"
)

// ---------------------------------------------------------------------------
// mockStore implements portfolio.Store for testing.
// ---------------------------------------------------------------------------

type mockStore struct {
	portfolios map[string]*model.Portfolio
	tickers    []model.AvailableTicker
	nextID     int
}

func newMockStore() *mockStore {
	return &mockStore{
		portfolios: make(map[string]*model.Portfolio),
		tickers: []model.AvailableTicker{
			{Symbol: "SBER", Name: "Sberbank", Market: "stock", Board: "TQBR"},
			{Symbol: "GAZP", Name: "Gazprom", Market: "stock", Board: "TQBR"},
			{Symbol: "TATN", Name: "Tatneft", Market: "stock", Board: "TQBR"},
		},
		nextID: 1,
	}
}

func (m *mockStore) GetPortfolio(_ context.Context, userID string) (*model.Portfolio, error) {
	p, ok := m.portfolios[userID]
	if !ok {
		return nil, portfolio.ErrPortfolioNotFound
	}
	return p, nil
}

func (m *mockStore) CreatePortfolio(_ context.Context, p *model.Portfolio) (*model.Portfolio, error) {
	id := m.nextID
	m.nextID++
	now := time.Now()
	created := &model.Portfolio{
		ID:        id,
		UserID:    p.UserID,
		Name:      p.Name,
		Tickers:   []model.Ticker{},
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.portfolios[p.UserID] = created
	return created, nil
}

func (m *mockStore) AddTickers(_ context.Context, portfolioID int, symbols []string) ([]model.Ticker, error) {
	var added []model.Ticker
	for _, sym := range symbols {
		t := model.Ticker{
			ID:          m.nextID,
			PortfolioID: portfolioID,
			Symbol:      sym,
			AddedAt:     time.Now(),
		}
		m.nextID++
		added = append(added, t)
	}
	// Update the portfolio's tickers in the mock
	for _, p := range m.portfolios {
		if p.ID == portfolioID {
			p.Tickers = append(p.Tickers, added...)
			p.UpdatedAt = time.Now()
			break
		}
	}
	return added, nil
}

func (m *mockStore) RemoveTicker(_ context.Context, portfolioID int, symbol string) error {
	for _, p := range m.portfolios {
		if p.ID == portfolioID {
			for i, t := range p.Tickers {
				if t.Symbol == symbol {
					p.Tickers = append(p.Tickers[:i], p.Tickers[i+1:]...)
					p.UpdatedAt = time.Now()
					return nil
				}
			}
			return portfolio.ErrTickerNotFound
		}
	}
	return portfolio.ErrPortfolioNotFound
}

func (m *mockStore) GetTickers(_ context.Context, portfolioID int) ([]model.Ticker, error) {
	for _, p := range m.portfolios {
		if p.ID == portfolioID {
			return p.Tickers, nil
		}
	}
	return nil, portfolio.ErrPortfolioNotFound
}

func (m *mockStore) ListAvailableTickers(_ context.Context) ([]model.AvailableTicker, error) {
	return m.tickers, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func newTestHandler() *Handler {
	log, _ := logger.New(&logger.Config{LokiEnabled: false})
	svc := portfolio.NewService(newMockStore(), log)
	return NewHandler(svc)
}

func newTestRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/tickers", h.ListTickers).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/portfolio", h.GetPortfolio).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/portfolio", h.AddToPortfolio).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/portfolio/{ticker}", h.RemoveFromPortfolio).Methods(http.MethodDelete)
	return r
}

func decodeResponse(t *testing.T, body []byte, v interface{}) {
	t.Helper()
	if err := json.Unmarshal(body, v); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestListTickers_Success(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tickers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp TickersResponse
	decodeResponse(t, rec.Body.Bytes(), &resp)

	if len(resp.Tickers) != 3 {
		t.Fatalf("expected 3 tickers, got %d", len(resp.Tickers))
	}
	if resp.Tickers[0].Symbol != "SBER" {
		t.Fatalf("expected first ticker SBER, got %s", resp.Tickers[0].Symbol)
	}
}

func TestGetPortfolio_NotFound(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/portfolio", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Service should auto-create a portfolio or return error —
	// our handler currently returns 500 for store errors.
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for missing portfolio, got %d", rec.Code)
	}

	var errResp ErrorResponse
	decodeResponse(t, rec.Body.Bytes(), &errResp)
	if errResp.Error != "internal_error" {
		t.Fatalf("expected error code 'internal_error', got '%s'", errResp.Error)
	}
}

func TestGetPortfolio_Success(t *testing.T) {
	h := newTestHandler()

	// Pre-create a portfolio via AddToPortfolio
	router := newTestRouter(h)
	body := `{"tickers":["SBER","GAZP"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 for add, got %d", rec.Code)
	}

	// Now fetch the portfolio
	req = httptest.NewRequest(http.MethodGet, "/api/v1/portfolio", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var portfolioResp model.Portfolio
	decodeResponse(t, rec.Body.Bytes(), &portfolioResp)

	if len(portfolioResp.Tickers) != 2 {
		t.Fatalf("expected 2 tickers, got %d", len(portfolioResp.Tickers))
	}
}

func TestAddToPortfolio_Success(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	body := `{"tickers":["SBER","GAZP","TATN"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp AddTickersResponse
	decodeResponse(t, rec.Body.Bytes(), &resp)

	if resp.Added != 3 {
		t.Fatalf("expected 3 added tickers, got %d", resp.Added)
	}
}

func TestAddToPortfolio_EmptyBody(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestAddToPortfolio_EmptyTickers(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	body := `{"tickers":[]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

func TestAddToPortfolio_InvalidSymbols(t *testing.T) {
	h := newTestHandler()
	router := newTestRouter(h)

	body := `{"tickers":["INVALID1","INVALID2"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status 422, got %d", rec.Code)
	}

	var errResp ErrorResponse
	decodeResponse(t, rec.Body.Bytes(), &errResp)
	if errResp.Error != "invalid_symbols" {
		t.Fatalf("expected error 'invalid_symbols', got '%s'", errResp.Error)
	}
}

func TestRemoveFromPortfolio_Success(t *testing.T) {
	h := newTestHandler()

	// First add tickers
	router := newTestRouter(h)
	body := `{"tickers":["SBER","GAZP"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup: expected 200, got %d", rec.Code)
	}

	// Now remove one
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/portfolio/SBER", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp RemoveTickerResponse
	decodeResponse(t, rec.Body.Bytes(), &resp)
	if !resp.Removed {
		t.Fatal("expected Removed to be true")
	}
	if resp.Ticker != "SBER" {
		t.Fatalf("expected ticker 'SBER', got '%s'", resp.Ticker)
	}
}

func TestRemoveFromPortfolio_NotFound(t *testing.T) {
	h := newTestHandler()

	// First add tickers
	router := newTestRouter(h)
	body := `{"tickers":["SBER"]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("setup: expected 200, got %d", rec.Code)
	}

	// Try to remove a ticker not in portfolio
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/portfolio/INVALID", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}

	var errResp ErrorResponse
	decodeResponse(t, rec.Body.Bytes(), &errResp)
	if errResp.Error != "not_found" {
		t.Fatalf("expected error 'not_found', got '%s'", errResp.Error)
	}
}

func TestResponseHelpers_ErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()
	errorResponse(w, http.StatusNotFound, "not_found", "resource not found")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}

	var resp ErrorResponse
	decodeResponse(t, w.Body.Bytes(), &resp)
	if resp.Error != "not_found" {
		t.Fatalf("expected error 'not_found', got '%s'", resp.Error)
	}
	if resp.Message != "resource not found" {
		t.Fatalf("expected message 'resource not found', got '%s'", resp.Message)
	}
	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected code 404, got %d", resp.Code)
	}
}

func TestResponseHelpers_JSONResponseWithNil(t *testing.T) {
	w := httptest.NewRecorder()
	jsonResponse(w, http.StatusNoContent, nil)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}
