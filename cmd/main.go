package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/senyabanana/avito-shop-service/internal/config"
	"github.com/senyabanana/avito-shop-service/internal/database"
	"github.com/senyabanana/avito-shop-service/internal/handler"
	"github.com/senyabanana/avito-shop-service/internal/logger"
	"github.com/senyabanana/avito-shop-service/internal/repository"
	"github.com/senyabanana/avito-shop-service/internal/server"
	"github.com/senyabanana/avito-shop-service/internal/service"
	
	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
)

func main() {
	ctx := context.Background()
	log := logger.NewLogger()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := database.NewPostgresDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	repos := repository.NewRepository(db)
	services := service.NewService(repos, trManager, log)
	handlers := handler.NewHandler(services, log)

	srv := new(server.Server)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		if err := srv.Run(cfg.ServerPort, handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Info("Server started on port ", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Server shutting down...")
	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Errorf("error occured on db connection close: %s", err.Error())
	}

	log.Info("Server stopped gracefully")
}
