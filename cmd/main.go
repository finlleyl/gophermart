package main

import (
	"github.com/go-chi/chi/v5"
	"gophermart/internal/api"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"gophermart/pkg/gzip"
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

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", gzip.GzipMiddleware(
			api.RegisterHandler(db, cfg)),
		)
	})

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		log.Fatal(err)
	}
}
