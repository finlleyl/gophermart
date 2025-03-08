package main

import (
	"github.com/go-chi/chi/v5"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	cfg := config.LoadConfig()
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		log.Fatal(err)
	}
}
