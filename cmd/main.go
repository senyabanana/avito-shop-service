package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	// Основной контекст, управляемый сигналами
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	log := logger.NewLogger()

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := database.NewPostgresDB(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	defer db.Close()

	trManager := manager.Must(trmsqlx.NewDefaultFactory(db))
	repos := repository.NewRepository(db)
	services := service.NewService(repos, trManager, cfg.JwtSecretKey, log)
	handlers := handler.NewHandler(services, log)

	srv := new(server.Server)

	go func() {
		if err := srv.Run(ctx, cfg.ServerPort, handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	log.Info("Server started on port ", cfg.ServerPort)

	<-ctx.Done()

	log.Info("Server shutting down...")

	// Завершаем сервер с таймаутом
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorf("error occurred on server shutdown: %s", err.Error())
	}

	log.Info("Server stopped gracefully")
}
