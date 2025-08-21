package paste

import (
	"github.com/gabrielolivrp/pastebin-api/pkg/cache"
	"github.com/gabrielolivrp/pastebin-api/pkg/database"
	"github.com/gabrielolivrp/pastebin-api/pkg/http/midleware"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.RouterGroup,
	logger logging.Logger,
	dbClient database.Client,
	cacheClient cache.Client,
) {
	pasteRepository := NewPasteRepository(dbClient)
	pasteService := NewPasteService(pasteRepository, cacheClient)

	handler := NewPasteHandler(logger, pasteService)
	r.POST("/pastes", midleware.ErrorMiddleware(handler.CreateHandler))
	r.GET("/pastes/:id", midleware.ErrorMiddleware(handler.GetByIdHandler))
}
