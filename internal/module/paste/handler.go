package paste

import (
	"errors"

	"github.com/gabrielolivrp/pastebin-api/pkg/http/response"
	"github.com/gabrielolivrp/pastebin-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type pasteHandler struct {
	service PasteService
	logger  logging.Logger
}

func NewPasteHandler(logger logging.Logger, service PasteService) *pasteHandler {
	return &pasteHandler{
		service,
		logger,
	}
}

type CreatePasteRequest struct {
	Content string `json:"content" form:"content" binding:"required,min=10"`
	Title   string `json:"title" form:"title" binding:"required,max=255"`
	Lang    string `json:"lang" form:"lang" binding:"required,max=50"`
}

func (h *pasteHandler) CreateHandler(c *gin.Context) *response.APIResponse {
	var req CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			h.logger.Error("Validation error", logging.Field{
				Key:   "error",
				Value: err.Error(),
			})
			return response.ValidationError(response.ParseValidationErrors(err))
		}
		h.logger.Error("Binding error", logging.Field{
			Key:   "error",
			Value: err.Error(),
		})
		return response.InternalServerError(err.Error())
	}

	result, err := h.service.Create(c.Request.Context(), CreatePasteParams(req))
	if err != nil {
		h.logger.Error("Failed to create paste", logging.Field{
			Key:   "error",
			Value: err.Error(),
		})
		return response.InternalServerError(err.Error())
	}

	return response.Created(result)
}

func (h *pasteHandler) GetByIdHandler(c *gin.Context) *response.APIResponse {
	id := c.Param("id")
	if id == "" {
		h.logger.Error("Invalid paste ID", logging.Field{
			Key:   "id",
			Value: id,
		})
		return response.BadRequest("Invalid paste ID")
	}

	paste, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrPasteNotFound) {
			h.logger.Error("Paste not found", logging.Field{
				Key:   "id",
				Value: id,
			})
			return response.NotFound("Paste")
		}
		h.logger.Error("Failed to get paste by ID", logging.Field{
			Key:   "error",
			Value: err.Error(),
		})
		return response.InternalServerError(err.Error())
	}

	return response.OK(paste)
}
