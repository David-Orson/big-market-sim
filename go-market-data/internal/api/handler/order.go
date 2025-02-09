package handler

import (
	"database/sql"
	"encoding/json"
	"go-market-data/m/v2/internal/models"
	"go-market-data/m/v2/pkg/utils/web"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type OrderHandler struct {
	DB *sql.DB
}

func (oh *OrderHandler) Setup(mux *http.ServeMux) {
	mux.HandleFunc("GET /order-book", oh.getOrderBook)
	mux.HandleFunc("GET /trades", oh.getTrades)
	mux.HandleFunc("POST /place-order", oh.placeOrder)
}

func (oh *OrderHandler) getOrderBook(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing symbol", http.StatusBadRequest)
		return
	}

	// Fetch order book (dummy implementation for now)
	ob := getOrderBook(symbol)
	if ob == nil {
		http.Error(w, "Order book not found", http.StatusNotFound)
		return
	}

	web.RespondJSON(w, ob, http.StatusOK)
}

func getOrderBook(symbol string) *models.OrderBook {
	// Dummy implementation for now
	return &models.OrderBook{
		Symbol: symbol,
		Bids: []*models.Order{
			{ID: "1", Price: 100, Quantity: 10},
			{ID: "2", Price: 99, Quantity: 5},
		},
		Asks: []*models.Order{
			{ID: "3", Price: 101, Quantity: 15},
			{ID: "4", Price: 102, Quantity: 20},
		},
	}
}

func (oh *OrderHandler) getTrades(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Missing symbol", http.StatusBadRequest)
		return
	}

	// Fetch trades for symbol
	trades := getRecentTrades(symbol)
	web.RespondJSON(w, trades, http.StatusOK)
}

func getRecentTrades(symbol string) []*models.Trade {
	// Dummy implementation for now
	return []*models.Trade{
		{ID: "1", Price: 100, Quantity: 5, Timestamp: time.Now()},
		{ID: "2", Price: 101, Quantity: 10, Timestamp: time.Now()},
	}
}

func (oh *OrderHandler) placeOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if order.Symbol == "" || order.Quantity <= 0 {
		http.Error(w, "Invalid order details", http.StatusBadRequest)
		return
	}

	order.PlacedAt = time.Now()
	order.ID = uuid.New().String()
	order.Status = "pending"

	// Store order in database (or in-memory order book)
	// Here, we're simulating with a simple print for now
	log.Printf("Received Order: %+v\n", order)

	web.Respond(w, "Order placed", http.StatusOK)
}
