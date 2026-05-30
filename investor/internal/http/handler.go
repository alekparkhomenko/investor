package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alekparkhomenko/investor/investor/internal/model"
	"github.com/alekparkhomenko/investor/investor/internal/portfolio"
	"github.com/gorilla/mux"
)

// TickersResponse is the response for GET /api/v1/tickers.
type TickersResponse struct {
	Tickers []model.AvailableTicker `json:"tickers"`
}

// AddTickersRequest is the request for POST /api/v1/portfolio.
type AddTickersRequest struct {
	Tickers []string `json:"tickers"`
}

// AddTickersResponse is the response for POST /api/v1/portfolio.
type AddTickersResponse struct {
	Added   int      `json:"added"`
	Tickers []string `json:"tickers"`
}

// RemoveTickerResponse is the response for DELETE /api/v1/portfolio/{ticker}.
type RemoveTickerResponse struct {
	Removed bool   `json:"removed"`
	Ticker  string `json:"ticker"`
}

// ErrorResponse is the standard error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Handler handles HTTP requests for portfolio.
type Handler struct {
	svc *portfolio.Service
}

// NewHandler creates a new Handler.
func NewHandler(svc *portfolio.Service) *Handler {
	return &Handler{
		svc: svc,
	}
}

// ListTickers returns list of available MOEX tickers.
// @Summary List available MOEX tickers
// @Description Get all available MOEX stock tickers
// @Tags tickers
// @Accept json
// @Produce json
// @Success 200 {object} TickersResponse
// @Router /api/v1/tickers [get]
func (h *Handler) ListTickers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tickers, err := h.svc.ListAvailableTickers(ctx)
	if err != nil {
		writeError(w, "internal_error", "Failed to fetch tickers", http.StatusInternalServerError)
		return
	}

	response := TickersResponse{
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

	// TODO: extract userID from auth context in PORT-003
	userID := "default"
	portfolioResult, err := h.svc.GetPortfolio(ctx, userID)
	if err != nil {
		writeError(w, "internal_error", "Failed to fetch portfolio", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portfolioResult)
}

// AddToPortfolio adds tickers to portfolio.
// @Summary Add tickers to portfolio
// @Description Add tickers to user's portfolio
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} AddTickersResponse
// @Router /api/v1/portfolio [post]
func (h *Handler) AddToPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req AddTickersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid_request", "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Tickers) == 0 {
		writeError(w, "invalid_request", "No tickers provided", http.StatusBadRequest)
		return
	}

	// TODO: extract userID from auth context in PORT-003
	userID := "default"
	_, err := h.svc.AddTickers(ctx, userID, req.Tickers)
	if err != nil {
		writeError(w, "internal_error", "Failed to add tickers", http.StatusInternalServerError)
		return
	}

	response := AddTickersResponse{
		Added:   len(req.Tickers),
		Tickers: req.Tickers,
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
// @Success 200 {object} RemoveTickerResponse
// @Param ticker path string true "Ticker symbol"
// @Router /api/v1/portfolio/{ticker} [delete]
func (h *Handler) RemoveFromPortfolio(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	ticker := vars["ticker"]

	// TODO: extract userID from auth context in PORT-003
	userID := "default"
	err := h.svc.RemoveTicker(ctx, userID, ticker)
	if err != nil {
		if errors.Is(err, portfolio.ErrTickerNotFound) {
			writeError(w, "not_found", "Ticker not in portfolio", http.StatusNotFound)
			return
		}
		writeError(w, "internal_error", "Failed to remove ticker", http.StatusInternalServerError)
		return
	}

	response := RemoveTickerResponse{
		Removed: true,
		Ticker:  ticker,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, err, msg string, code int) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err,
		Message: msg,
		Code:    code,
	})
}
