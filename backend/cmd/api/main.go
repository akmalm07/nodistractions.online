package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"nodistractions-online/backend/internal/auth"
	"nodistractions-online/backend/internal/config"
	"nodistractions-online/backend/internal/db"
	"nodistractions-online/backend/internal/httpapi"
)

func main() {
	cfg, err := config.Load(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(database)

	api := httpapi.NewServer(
		db.NewUserRepository(database),
		auth.NewSessionManager(cfg.SessionSecret, cfg.JWTIssuer, cfg.JWTAudience),
	)

	slog.Info("API listening", "address", cfg.ListenAddress)
	log.Fatal(http.ListenAndServe(cfg.ListenAddress, httpapi.CORS(cfg.FrontendURL, api.Routes())))
}
