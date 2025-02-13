package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"shop/internal/controller"
	"shop/internal/repository"
	"shop/internal/usecase"
	"shop/pkg/database"
	"shop/pkg/logger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file")
	}

	port := os.Getenv("HTTP_PORT")
	db := database.InitializeDBPostgres(3, 10)
	//db.Seed()
	logger.InitLogger()

	repository := repository.NewRepository(db.GetDB())
	usecase := usecase.NewUsecase(repository)
	handlers := controller.NewHandler(usecase)
	router := handlers.Handle()

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("connection failed: %s\n", err.Error())
	}

	log.Infof("server is running on port %s\n", port)
}
