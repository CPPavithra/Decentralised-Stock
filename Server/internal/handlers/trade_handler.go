package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"trading-platform-backend/internal/models"
	"trading-platform-backend/internal/repositories"
	"trading-platform-backend/internal/services"

	"github.com/google/uuid"
)

type TradeHandler struct {
	TradeRepo    *repositories.TradeRepository
	AptosService *services.AptosService
}

// BuyStock handles buying stocks
func (h *TradeHandler) BuyStock(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID") // Retrieved from AuthMiddleware

	var req struct {
		StockSymbol string  `json:"stock_symbol"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		Action      string  `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create new trade object
	trade := models.Trade{
		ID:          uuid.New().String(),
		UserID:      userID,
		StockSymbol: req.StockSymbol,
		Quantity:    req.Quantity,
		Price:       req.Price,
	}

	// Store trade in database
	err := h.TradeRepo.CreateTrade(&trade, req.Action)
	if err != nil {
		http.Error(w, "Failed to create trade", http.StatusInternalServerError)
		return
	}

	// Mint stock tokens on the blockchain
	err = h.AptosService.MintStockToken(context.Background(), req.StockSymbol, req.Quantity, userID, "4c3a1a2e5f6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stock purchased successfully"})
}

// SellStock handles selling stocks
func (h *TradeHandler) SellStock(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID") // Retrieved from AuthMiddleware

	var req struct {
		StockSymbol string  `json:"stock_symbol"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		Action      string  `json:"action"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Create new trade object
	trade := models.Trade{
		ID:          uuid.New().String(),
		UserID:      userID,
		StockSymbol: req.StockSymbol,
		Quantity:    req.Quantity,
		Price:       req.Price,
	}

	// Store trade in database
	err := h.TradeRepo.CreateTrade(&trade, req.Action)
	if err != nil {
		http.Error(w, "Failed to create trade", http.StatusInternalServerError)
		return
	}

	// Transfer stock tokens on the blockchain
	err = h.AptosService.TransferStockToken(context.Background(), userID, "receiverAddress", req.StockSymbol, req.Quantity, trade.ID) // Replace "receiverAddress" with actual receiver
	if err != nil {
		http.Error(w, "Failed to transfer stock tokens", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stock sold successfully"})
}
