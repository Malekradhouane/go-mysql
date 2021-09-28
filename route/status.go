package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStatus service status
func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
