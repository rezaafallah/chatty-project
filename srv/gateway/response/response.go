package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"success": false, "error": msg})
}