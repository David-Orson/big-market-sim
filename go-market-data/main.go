package main

import (
	"fmt"
	"go-market-data/m/v2/internal/api"
	"go-market-data/m/v2/internal/models"
	"go-market-data/m/v2/pkg/utils/db"
	"go-market-data/m/v2/pkg/utils/web"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"golang.org/x/exp/rand"
)

func main() {
	fmt.Println("Starting server...")

	db, err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database: ", err)
		return
	}
	defer db.Sql.Close()

	handler := api.NewHandler(db.Sql)
	ready := make(chan bool, 1)
	cancel := make(chan bool, 1)
	port := "8087"

	web.Serve(handler, port, ready, cancel)

	ob := &models.OrderBook{
		Symbol:    "BTCUSD",
		Bids:      []*models.Order{},
		Asks:      []*models.Order{},
		TradeLog:  []*models.Trade{},
		OrderLock: sync.RWMutex{},
	}
	go generateRandomOrders(ob)

	<-ready
	fmt.Println("Server is ready! on port " + port)
	<-cancel
	fmt.Println("Server is shutting down...")
}

func generateRandomOrders(ob *models.OrderBook) {
	// let's measure the time performance of the order book
	start := time.Now()
	rand.Seed(uint64(time.Now().UnixNano()))
	i := 0
	for {
		i++
		randInt := rand.Intn(1000)
		order := &models.Order{
			ID:        uuid.New().String(),
			Symbol:    ob.Symbol,
			Price:     int64(10000 + randInt), // Random price around 100
			Quantity:  rand.Intn(20) + 1,      // Random quantity (1 to 20)
			OrderType: models.OrderTypeLimit,
			Side:      []models.OrderSide{models.OrderSideBuy, models.OrderSideSell}[rand.Intn(2)],
			PlacedAt:  time.Now(),
			Status:    models.OrderStatusPending,
		}
		ob.MatchOrders(order)
		if i%100000 == 0 {
			fmt.Println("Order generated")
			fmt.Println("bids: ", ob.Bids)
			fmt.Println("asks: ", ob.Asks)
			fmt.Println("trades: ", ob.TradeLog)
			buy, sell, ok := ob.GetBuyAndSellPrices()
			if ok {
				fmt.Println("Best listed buy price: ", buy)
				fmt.Println("Best listed sell price: ", sell)
			} else {
				fmt.Println("Order book is empty")
			}
			if buy != 0 && sell != 0 && buy > sell {
				panic("Buy price is greater than sell price")
			}

			lastTradePrice, ok := ob.GetLastTradePrice()
			if ok {
				fmt.Println("Last traded price: ", lastTradePrice)
			} else {
				fmt.Println("No trades executed yet")
			}

			midPoint, ok := ob.GetMidPoint()
			if ok {
				fmt.Println("Mid point: ", midPoint)
			} else {
				fmt.Println("Order book is empty")
			}
			fmt.Println("Orders processed: ", i)
			size := calculateOrderBookSize(ob)
			fmt.Printf("Estimated OrderBook size: %v bytes (%v MiB)\n", size, size/1024/1024)
			fmt.Println("Lengths: ", len(ob.Bids), len(ob.Asks), len(ob.TradeLog))
			fmt.Println("Time elapsed: ", time.Since(start))
			fmt.Println("OrdersPerMillisecond: ", int64(i)/time.Since(start).Milliseconds())
		}
	}
}

func calculateOrderBookSize(ob *models.OrderBook) uintptr {
	baseSize := unsafe.Sizeof(*ob)                                           // Size of the OrderBook struct
	bidsSize := uintptr(len(ob.Bids)) * unsafe.Sizeof(&models.Order{})       // Size of bid pointers
	asksSize := uintptr(len(ob.Asks)) * unsafe.Sizeof(&models.Order{})       // Size of ask pointers
	tradesSize := uintptr(len(ob.TradeLog)) * unsafe.Sizeof(&models.Trade{}) // Size of trade pointers
	return baseSize + bidsSize + asksSize + tradesSize
}
