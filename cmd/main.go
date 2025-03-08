package main

import (
	"gophermart/pkg/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	if err := http.ListenAndServe(cfg.RunAddress, nil); err != nil {
		log.Fatal(err)
	}
}
