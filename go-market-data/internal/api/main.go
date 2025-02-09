package api

import (
	"database/sql"
	"go-market-data/m/v2/internal/api/handler"
	"net/http"
)

type Handler struct {
	*handler.OrderHandler
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		OrderHandler: &handler.OrderHandler{
			DB: db,
		},
	}
}

func (h *Handler) SetupRoutes(mux *http.ServeMux) {
	h.OrderHandler.Setup(mux)
}
