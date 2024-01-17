package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/maxik12233/task-junior/internal/repository"
	"github.com/maxik12233/task-junior/internal/service"
	"github.com/maxik12233/task-junior/internal/transport"
	"github.com/maxik12233/task-junior/pkg/cors"
	"github.com/maxik12233/task-junior/pkg/logger"
	"github.com/maxik12233/task-junior/pkg/metrics"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Couldn't load env vars from .env file")
	}

	var dbUrl string
	dbUrl = os.Getenv("DB_CONNECTION_STRING")
	if dbUrl == "" {
		panic("DB_CONNECTION_STRING env is empty string")
	}
	dbSession, err := repository.EstablishDatabaseConnection(dbUrl)
	if err != nil {
		panic(fmt.Errorf("Fatal error database connection: %s \n", err))
	}

	logger := logger.Init()

	// Configure router
	router := gin.Default()
	// All routes using cors middleware
	router.Use(cors.CORSMiddleware())

	// Register general metrics endpoint
	metric := metrics.Metric{Logger: logger}
	metric.Register(router)

	// Logic
	repo := repository.NewRepository(dbSession, logger)
	svc := service.NewService(repo, logger)
	trans := transport.NewTransport(svc, logger)
	trans.RegisterRoutes(router)

	router.Run(":5050")
}
