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
	mux.HandleFunc("POST /place-order", oh.placeOrder)
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
