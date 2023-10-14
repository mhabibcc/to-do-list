package main

import (
	"net/http"
	"to-do-list/internal/config"

	"github.com/go-chi/chi/v5"
)

func startServer(router *chi.Mux, cfg *config.Config) error {
	return http.ListenAndServe(cfg.Server.HTTP.Address, router)

}
