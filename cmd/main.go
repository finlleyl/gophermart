package main

import (
	"github.com/go-chi/chi/v5"
	"gophermart/internal/api"
	"gophermart/internal/middleware"
	"gophermart/pkg/config"
	"gophermart/pkg/database"
	"gophermart/pkg/gzip"
	"gophermart/pkg/logger"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	cfg := config.LoadConfig()
	logger.InitLogger()
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", gzip.GzipMiddleware(
			api.RegisterHandler(db, cfg)),
		)

		r.Post("/login", gzip.GzipMiddleware(
			api.LoginHandler(db, cfg)),
		)

		r.Post("/orders", gzip.GzipMiddleware(
			middleware.CheckCookies(
				cfg, api.LoadOrderHandler(db),
			),
		),
		)
	})

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		logger.Sugar.Fatal(err)
	} else {
		logger.Sugar.Infow(
			"server started", "address", cfg.RunAddress)
	}
}
