package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/maxik12233/task-junior/internal/config"
	"github.com/maxik12233/task-junior/internal/repository"
	"github.com/maxik12233/task-junior/internal/service"
	"github.com/maxik12233/task-junior/internal/transport"
	"github.com/maxik12233/task-junior/pkg/api/logging"
	"github.com/maxik12233/task-junior/pkg/api/paginate"
	"github.com/maxik12233/task-junior/pkg/api/sort"
	"github.com/maxik12233/task-junior/pkg/cors"
	"github.com/maxik12233/task-junior/pkg/logger"
	"github.com/maxik12233/task-junior/pkg/metrics"
	"github.com/maxik12233/task-junior/pkg/name_info_sdk"
)

func main() {
	log := logger.Init()

	if err := godotenv.Load(".env"); err != nil {
		log.Info("Couldn't initialize env variables via .env file")
	}

	config.BindConfig("config.yaml")
	cfg := config.GetConfig()

	dbSession, err := repository.EstablishDatabaseConnection(cfg.PostgresConnectionString)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fatal error database connection: %s \n", err))
	}

	// Configure router
	var router *gin.Engine
	if !cfg.IsDebug {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	} else {
		router = gin.Default()
	}
	router.Use(cors.CORSMiddleware())
	router.Use(logging.ResponseLogger(log), logging.RequestLogger(log))
	router.Use(paginate.Middleware(cfg.DefaultPage, cfg.DefaultPerPage))
	router.Use(sort.Middleware(cfg.DefaultSortField, cfg.DefaultSortOrder))

	// Register general metrics endpoint
	metric := metrics.Metric{Logger: log}
	metric.Register(router)

	// Logic
	repo := repository.NewRepository(dbSession, log)
	svc := service.NewService(repo, log, name_info_sdk.NewNameInfo(""))
	trans := transport.NewTransport(svc, log)
	trans.RegisterRoutes(router)

	port := fmt.Sprintf(":%d", cfg.Port)
	log.Info(fmt.Sprintf("Running server on port %s...", port))
	router.Run(port)
}
