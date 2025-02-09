package db

import (
	"database/sql"
	"fmt"
	"go-market-data/m/v2/pkg/utils/web"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

type API struct {
	Sql *sql.DB
}

func InitDB() (API, error) {
	host := web.GetEnv("DB_HOST", "0.0.0.0")
	user := web.GetEnv("DB_USER", "root")
	password := web.GetEnv("DB_PASSWORD", "password")
	dbname := web.GetEnv("DB_NAME", "lf")
	port := web.GetEnv("DB_PORT", "13306")

	fmt.Println("port", port)

	config := mysql.Config{
		User:            user,
		Passwd:          password,
		Net:             "tcp",
		Addr:            host + ":" + port,
		DBName:          dbname,
		MultiStatements: true,
	}

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", config.FormatDSN())
		if err == nil {
			pingErr := db.Ping()
			if pingErr == nil {
				break
			}
		}

		log.Printf("Database connection failed. Retrying in 5 seconds... (%d/10)\n", i+1)
		time.Sleep(5 * time.Second)
	}

	p := db.Ping()
	if p != nil {
		fmt.Println("Error pinging db", p)
	}

	/*
		c, err := os.ReadFile("./migration.sql")
		if err != nil {
			fmt.Println("Error reading migrations file", err)
			return API{}, err
		}

		mg := string(c)

		_, err = db.Exec(mg)
		if err != nil {
			fmt.Println("Error executing migrations", err)
			return API{}, err
		}
	*/

	return API{Sql: db}, nil
}
