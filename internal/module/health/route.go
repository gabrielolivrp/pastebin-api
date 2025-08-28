package health

import (
	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/gabrielolivrp/pastebin-api/pkg/http/middleware"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	logger logging.Logger,
	dbClient database.Client,
	cacheClient cache.Client,
) {
	healthCheckHandler := NewHealthHandler(logger, dbClient, cacheClient)
	r.GET("/healthz", middleware.ErrorMiddleware(healthCheckHandler.HealthCheck))
}
