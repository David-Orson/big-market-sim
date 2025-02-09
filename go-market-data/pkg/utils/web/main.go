package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Handler interface {
	SetupRoutes(mux *http.ServeMux)
}

func Serve[H Handler](h H, port string, ready chan<- bool, cancel chan<- bool) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ROUTE DOES NOT EXIST")
		RespondError(w, "route does not exist", 404)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I'm healthy!")
		Respond(w, "OK", http.StatusOK)
	})

	h.SetupRoutes(mux)

	corsMux := WithCORS(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("server error: %v", err)
			cancel <- true
		}
	}()
	ready <- true
}

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Respond(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if message != "" {
		w.Write([]byte(message))
	}
}

func RespondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func RespondError(w http.ResponseWriter, err string, status int) {
	fmt.Println(err)
	Respond(w, err, status)
}

func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

type AuthToken struct {
	Token string `json:"token"`
}

func Authenticate(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Authenticate")
		token := r.Header.Get("Authorization")
		if token == "" {
			RespondError(w, "No token provided", 401)
			return
		}
		token = token[7:]

		rows, err := db.Query("SELECT token FROM auth_token WHERE token = ?", token)
		if err != nil {
			fmt.Println("Select token error: ", err)
			RespondError(w, err.Error(), 500)
			return
		}

		defer rows.Close()

		var tokens []string
		for rows.Next() {
			var token string
			if err := rows.Scan(&token); err != nil {
				fmt.Println("Error reading data: ", err)
				RespondError(w, err.Error(), 500)
			}
			tokens = append(tokens, token)
		}

		found := false

		for _, t := range tokens {
			if t == token {
				found = true
			}
		}

		if !found {
			RespondError(w, "Invalid token", 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
