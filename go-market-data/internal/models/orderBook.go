package models

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Trade represents an executed trade
type Trade struct {
	ID          string    `json:"trade_id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       int64     `json:"price"`
	Quantity    int       `json:"quantity"`
	Timestamp   time.Time `json:"timestamp"`
}

type OrderBook struct {
	Symbol    string       // Symbol the order book is for
	Bids      []*Order     // Buy orders (sorted by price descending)
	Asks      []*Order     // Sell orders (sorted by price ascending)
	TradeLog  []*Trade     // Completed trades
	OrderLock sync.RWMutex // Mutex to handle concurrent access
}

func (ob *OrderBook) GetBuyAndSellPrices() (buyPrice int64, sellPrice int64, ok bool) {
	ob.OrderLock.RLock()
	defer ob.OrderLock.RUnlock()

	// Check if there are bids and asks
	if len(ob.Bids) > 0 {
		buyPrice = ob.Bids[0].Price // Best bid (highest)
	}
	if len(ob.Asks) > 0 {
		sellPrice = ob.Asks[0].Price // Best ask (lowest)
	}

	// If both are zero, the order book is empty
	if buyPrice == 0 && sellPrice == 0 {
		return 0, 0, false
	}
	return buyPrice, sellPrice, true
}

func (ob *OrderBook) GetMidPoint() (int64, bool) {
	ob.OrderLock.RLock()
	defer ob.OrderLock.RUnlock()
	if len(ob.Bids) > 0 && len(ob.Asks) > 0 {
		bestBid := ob.Bids[0].Price
		bestAsk := ob.Asks[0].Price
		return (bestBid + bestAsk) / 2, true
	}
	return 0, false
}

func (ob *OrderBook) GetLastTradePrice() (int64, bool) {
	ob.OrderLock.RLock()
	defer ob.OrderLock.RUnlock()

	// 1. Use the last traded price if available
	if len(ob.TradeLog) > 0 {
		lastTrade := ob.TradeLog[len(ob.TradeLog)-1]
		return lastTrade.Price, true
	}
	return 0, false
}

func (ob *OrderBook) MatchOrders(order *Order) []Trade {
	ob.OrderLock.Lock()
	defer ob.OrderLock.Unlock()

	var trades []Trade
	if order.Side == OrderSideBuy {
		// Match against asks (sell orders)
		for len(ob.Asks) > 0 {
			bestAsk := ob.Asks[0]
			if order.Price < bestAsk.Price { // No match
				break
			}

			// Determine trade quantity
			tradeQuantity := min(order.Quantity, bestAsk.Quantity)
			trade := Trade{
				ID:          uuid.New().String(),
				BuyOrderID:  order.ID,
				SellOrderID: bestAsk.ID,
				Price:       bestAsk.Price,
				Quantity:    tradeQuantity,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)
			ob.TradeLog = append(ob.TradeLog, &trade)

			// Adjust quantities and remove fully filled orders
			order.Quantity -= tradeQuantity
			bestAsk.Quantity -= tradeQuantity
			if bestAsk.Quantity == 0 {
				ob.Asks = ob.Asks[1:] // Remove fully filled ask
			}
			if order.Quantity == 0 {
				break // Fully filled buy order
			}
		}
	} else if order.Side == OrderSideSell {
		// Match against bids (buy orders)
		for len(ob.Bids) > 0 {
			bestBid := ob.Bids[0]
			if order.Price > bestBid.Price { // No match
				break
			}

			// Determine trade quantity
			tradeQuantity := min(order.Quantity, bestBid.Quantity)
			trade := Trade{
				ID:          uuid.New().String(),
				BuyOrderID:  bestBid.ID,
				SellOrderID: order.ID,
				Price:       bestBid.Price,
				Quantity:    tradeQuantity,
				Timestamp:   time.Now(),
			}
			trades = append(trades, trade)
			ob.TradeLog = append(ob.TradeLog, &trade)

			// Adjust quantities and remove fully filled orders
			order.Quantity -= tradeQuantity
			bestBid.Quantity -= tradeQuantity
			if bestBid.Quantity == 0 {
				ob.Bids = ob.Bids[1:] // Remove fully filled bid
			}
			if order.Quantity == 0 {
				break // Fully filled sell order
			}
		}
	}

	// Add unmatched order to the book
	if order.Quantity > 0 {
		if order.Side == OrderSideBuy {
			ob.Bids = append(ob.Bids, order)
			ob.SortBids() // Keep Bids sorted
		} else {
			ob.Asks = append(ob.Asks, order)
			ob.SortAsks() // Keep Asks sorted
		}
	}

	return trades
}

// SortBids sorts buy orders by descending price
func (ob *OrderBook) SortBids() {
	sort.Slice(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})
}

// SortAsks sorts sell orders by ascending price
func (ob *OrderBook) SortAsks() {
	sort.Slice(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})
}
