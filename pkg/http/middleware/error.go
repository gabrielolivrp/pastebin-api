package middleware

import (
	"github.com/gabrielolivrp/pastebin-api/pkg/http/response"
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) *response.APIResponse

func ErrorMiddleware(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			c.JSON(err.Code, err)
		}
		c.Next()
	}
}
