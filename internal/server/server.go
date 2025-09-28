package server

import (
	"github.com/gabrielolivrp/pastebin-api/internal/module/health"
	"github.com/gabrielolivrp/pastebin-api/internal/module/paste"
	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/config"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func registerRoutes(
	logger logging.Logger,
	dbClient database.Client,
	cacheClient cache.Client,
) *gin.Engine {
	e := gin.Default()
	e.Use(cors.Default())
	group := e.Group("/api/v1")
	{
		health.RegisterRoutes(group, logger, dbClient, cacheClient)
		paste.RegisterRoutes(group, logger, dbClient, cacheClient)
	}
	return e
}

func Start(config *config.Config) error {
	logger, err := logging.NewLogger(config.Environment)
	if err != nil {
		return err
	}

	dbClient, err := database.NewClient(database.ClientConfig{
		Host:     config.DB.Host,
		Port:     config.DB.Port,
		Username: config.DB.Username,
		Password: config.DB.Password,
		Database: config.DB.Database,
		SSLMode:  config.DB.SSLMode,
	})

	if err != nil {
		logger.Error("Failed to connect to database", logging.Field{
			Key:   "error",
			Value: err,
		})
		return err
	}

	cacheClient, err := cache.NewClient(cache.ClientConfig{
		Host:     config.Cache.Host,
		Port:     config.Cache.Port,
		Password: config.Cache.Password,
	})
	if err != nil {
		logger.Error("Failed to connect to cache", logging.Field{
			Key:   "error",
			Value: err,
		})
		return err
	}

	if err := cacheClient.Ping(); err != nil {
		logger.Error("Failed to connect to cache", logging.Field{
			Key:   "error",
			Value: err,
		})
	}

	e := registerRoutes(logger, dbClient, cacheClient)
	if err := e.Run(":" + config.Port); err != nil {
		logger.Error("Failed to start server", logging.Field{
			Key:   "error",
			Value: err,
		})
		return err
	}
	return nil
}
