package main

import (
	"log"
	
	"github.com/senyabanana/avito-shop-service/internal/handler"
	"github.com/senyabanana/avito-shop-service/internal/server"
)

func main() {
	handlers := new(handler.Handler)

	srv := new(server.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
