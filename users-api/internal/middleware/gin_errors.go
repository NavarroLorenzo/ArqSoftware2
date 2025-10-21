package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func JSONError(c *gin.Context, status int, code, msg string) {
	c.AbortWithStatusJSON(status, gin.H{"error": ErrorResponse{Code: code, Message: msg}})
}

// Captura panic y responde JSON 500 (al final de la cadena)
func RecoverJSON() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		JSONError(c, http.StatusInternalServerError, "panic", "internal error")
	})
}
