package main

import (
	"fmt"
	"net/http"

	"github.com/MotiurRahmanSany/student-management-api/internal/config"
	"github.com/MotiurRahmanSany/student-management-api/internal/database"
)

func main() {
	config := config.GetConfig()
	fmt.Printf("Starting %s version %s on port %d\n", config.AppName, config.Version, config.HttpPort)
	fmt.Printf("Database host: %s, port: %d, user: %s\n", config.Db.Host, config.Db.Port, config.Db.User)

	pool, err := database.NewConnection(config.Db)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer pool.Close()
	// queries := db.New(pool)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})
	fmt.Printf("Server is running on port %d\n", config.HttpPort)
	http.ListenAndServe(fmt.Sprintf(":%d", config.HttpPort), mux)
}
