package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/programmer8760/demonlist/backend/database"
	"github.com/programmer8760/demonlist/backend/internal/logger"
)

func main() {
	logLevel := slog.LevelInfo
	log := logger.New(logLevel)

	if _, ok := os.LookupEnv("HTTP_PORT"); !ok {
		err := godotenv.Load()
		if err != nil {
			log.Error("loading .env file failed", "err", err)
			os.Exit(1)
		}
	}

	db, err := database.Connect()
	defer db.Close()
	if err != nil {
		os.Exit(1)
	}

	log.Info("db connection established")
}
