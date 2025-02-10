package main

import (
	"log"

	"github.com/senyabanana/avito-shop-service/internal/handler"
	"github.com/senyabanana/avito-shop-service/internal/repository"
	"github.com/senyabanana/avito-shop-service/internal/server"
	"github.com/senyabanana/avito-shop-service/internal/service"

	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializating configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
