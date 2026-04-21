package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/investor/internal/storage"
	"github.com/gorilla/mux"
)

// Handler handles HTTP requests for portfolio.
type Handler struct {
	store *storage.PortfolioStore
}

// NewHandler creates a new Handler.
func NewHandler(store *storage.PortfolioStore) *Handler {
	return &Handler{
		store: store,
	}
}

// ListTickers returns list of available MOEX tickers.
// @Summary List available MOEX tickers
// @Description Get all available MOEX stock tickers
// @Tags tickers
// @Accept json
// @Produce json
// @Success 200 {object} model.TickersResponse
// @Router /api/v1/tickers [get]
func (h *Handler) ListTickers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tickers, err := h.store.GetAllTickers(ctx)
	if err != nil {
		writeError(w, "internal_error", "Failed to fetch tickers", http.StatusInternalServerError)
		return
	}

	response := model.TickersResponse{
		Tickers: tickers,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPortfolio returns user's portfolio.
// @Summary Get user's portfolio
// @Description Get user's saved ticker portfolio
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} model.Portfolio
// @Router /api/v1/portfolio [get]
func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	portfolio, err := h.store.GetPortfolio(ctx)
	if err != nil {
		writeError(w, "internal_error", "Failed to fetch portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portfolio)
}

// AddToPortfolio adds tickers to portfolio.
// @Summary Add tickers to portfolio
// @Description Add tickers to user's portfolio
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} model.AddTickersResponse
// @Router /api/v1/portfolio [post]
func (h *Handler) AddToPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.AddTickersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Tickers) == 0 {
		writeError(w, "invalid_request", "No tickers provided", http.StatusBadRequest)
		return
	}

	added, err := h.store.AddTickers(ctx, req.Tickers)
	if err != nil {
		writeError(w, "internal_error", "Failed to add tickers", http.StatusInternalServerError)
		return
	}

	response := model.AddTickersResponse{
		Added:   len(added),
		Tickers: added,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RemoveFromPortfolio removes ticker from portfolio.
// @Summary Remove ticker from portfolio
// @Description Remove a ticker from user's portfolio
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} model.RemoveTickerResponse
// @Param ticker path string true "Ticker symbol"
// @Router /api/v1/portfolio/{ticker} [delete]
func (h *Handler) RemoveFromPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	ticker := vars["ticker"]

	err := h.store.RemoveTicker(ctx, ticker)
	if err != nil {
		if errors.Is(err, storage.ErrTickerNotFound) {
			writeError(w, "not_found", "Ticker not in portfolio", http.StatusNotFound)
			return
		}
		writeError(w, "internal_error", "Failed to remove ticker", http.StatusInternalServerError)
		return
	}

	response := model.RemoveTickerResponse{
		Removed: true,
		Ticker:  ticker,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, err, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error:   err,
		Message: msg,
		Code:    code,
	})
}