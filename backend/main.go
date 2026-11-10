package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/programmer8760/demonlist/backend/database"
	"github.com/programmer8760/demonlist/backend/internal/handler"
	"github.com/programmer8760/demonlist/backend/internal/logger"
	"github.com/programmer8760/demonlist/backend/internal/repository"
	"github.com/programmer8760/demonlist/backend/internal/service"
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

	listsRepo := repository.NewListsRepository(db)
	listsSvc := service.NewListsService(listsRepo, log)

	levelsRepo := repository.NewLevelsRepository(db)
	levelsSvc := service.NewListsService(levelsRepo, log)

	usersRepo := repository.NewUsersRepository(db)
	recordsRepo := repository.NewRecordsRepository(db)

	usersSvc := service.NewUsersService(usersRepo, recordsRepo, log)
	recordsSvc := service.NewRecordsService(recordsRepo, usersRepo, log)

	h := handler.NewHandler(log, listsSvc, levelsSvc, usersSvc, recordsSvc)

	addr := ":" + os.Getenv("HTTP_PORT")
	server := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	go func() {
		log.Info("http server started", "addr", addr)
		if err = server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("http server stopped", "err", err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
		_ = server.Close()
	}

	log.Info("server stopped gracefully")
}
