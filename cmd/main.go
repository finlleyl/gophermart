package main

import (
	"github.com/go-chi/chi/v5"
	"gophermart/pkg/config"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	cfg := config.LoadConfig()

	if err := http.ListenAndServe(cfg.RunAddress, r); err != nil {
		log.Fatal(err)
	}
}
