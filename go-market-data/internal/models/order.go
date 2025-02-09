package models

import (
	"time"
)

// OrderType represents the type of an order (market or limit)
type OrderType string

const (
	OrderTypeMarket OrderType = "market"
	OrderTypeLimit  OrderType = "limit"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending         OrderStatus = "pending"
	OrderStatusExecuted        OrderStatus = "executed"
	OrderStatusCanceled        OrderStatus = "canceled"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
)

// OrderSide represents whether an order is a buy or sell order
type OrderSide string

const (
	OrderSideBuy  OrderSide = "buy"
	OrderSideSell OrderSide = "sell"
)

// Order represents a buy/sell order
type Order struct {
	ID         string      `json:"order_id"`
	Symbol     string      `json:"symbol"`
	Price      float64     `json:"price,omitempty"` // Optional for market orders
	Quantity   int         `json:"quantity"`
	OrderType  OrderType   `json:"order_type"`
	Side       OrderSide   `json:"side"`
	Status     OrderStatus `json:"status"`
	PlacedAt   time.Time   `json:"placedAt"`              // Order creation time
	ExecutedAt time.Time   `json:"executed_at,omitempty"` // Optional
}
