package health

import (
	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/gabrielolivrp/pastebin-api/pkg/http/response"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger      logging.Logger
	dbClient    database.Client
	cacheClient cache.Client
}

func NewHealthHandler(
	logger logging.Logger,
	dbClient database.Client,
	cacheClient cache.Client,
) *HealthHandler {
	return &HealthHandler{logger, dbClient, cacheClient}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) *response.APIResponse {
	dbErr := h.dbClient.Ping()
	cacheErr := h.cacheClient.Ping()

	if dbErr != nil || cacheErr != nil {
		h.logger.Error("Health check failed", logging.Field{
			Key:   "database_error",
			Value: dbErr,
		}, logging.Field{
			Key:   "cache_error",
			Value: cacheErr,
		})
		return response.ServiceUnavailable("Some resources are unavailable")
	}

	return response.OK(map[string]interface{}{
		"status": "ok",
	})
}
