// Package classification of todolist API
//
// Documentation for todolist API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStatus service status
func GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
