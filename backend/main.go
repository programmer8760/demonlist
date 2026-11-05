package main

import (
	"errors"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
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

	if err = db.Ping(); err != nil {
		log.Error("database didn't respond to ping", "err", err)
		os.Exit(1)
	}
	log.Info("database connection established")

	err = database.MigrateUp(db)
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("database migrations: no changes")
		} else {
			log.Error("database migration failed", "err", err)
			os.Exit(1)
		}
	} else {
		log.Info("database migrations: success")
	}
}
