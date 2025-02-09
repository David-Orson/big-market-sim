package main

import (
	"fmt"
	"go-market-data/m/v2/internal/api"
	"go-market-data/m/v2/pkg/utils/db"
	"go-market-data/m/v2/pkg/utils/web"
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

	<-ready
	fmt.Println("Server is ready! on port " + port)
	<-cancel
	fmt.Println("Server is shutting down...")
}
