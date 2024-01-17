package main

import (
	"fmt"
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
	log := logger.Init()

	if err := godotenv.Load(".env"); err != nil {
		log.Info("Couldn't initialize env variables via .env file")
	}

	envs := ParseEnvs()

	if envs.dbConnectionString == nil {
		log.Fatal("Database connection string env is empty string")
	}
	dbSession, err := repository.EstablishDatabaseConnection(*envs.dbConnectionString)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fatal error database connection: %s \n", err))
	}

	var localPort string
	if envs.Port != nil {
		localPort = *envs.Port
	} else {
		log.Info("Local port env is empty string. Using default port.")
		localPort = ":5050"
	}

	// Configure router
	var router *gin.Engine
	if envs.Mode != nil && *envs.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	} else {
		router = gin.Default()
	}
	// All routes using cors middleware
	router.Use(cors.CORSMiddleware())
	router.Use(logger.ResponseLogger(log), logger.RequestLogger(log))

	// Register general metrics endpoint
	metric := metrics.Metric{Logger: log}
	metric.Register(router)

	// Logic
	repo := repository.NewRepository(dbSession, log)
	svc := service.NewService(repo, log)
	trans := transport.NewTransport(svc, log)
	trans.RegisterRoutes(router)

	log.Info(fmt.Sprintf("Running server on port %s...", localPort))
	router.Run(localPort)
}

type Envs struct {
	Mode               *string
	dbConnectionString *string
	Port               *string
}

func ParseEnvs() Envs {
	envs := Envs{}

	if port, exists := os.LookupEnv("LOCAL_PORT"); exists {
		envs.Port = &port
	} else {
		envs.Port = nil
	}

	if dbUrl, exists := os.LookupEnv("DB_CONNECTION_STRING"); exists {
		envs.dbConnectionString = &dbUrl
	} else {
		envs.dbConnectionString = nil
	}

	if mode, exists := os.LookupEnv("MODE"); exists {
		envs.Mode = &mode
	} else {
		envs.Mode = nil
	}

	return envs
}
