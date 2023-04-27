package handlers

import (
	"net/http"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func Readiness(c *gin.Context) {
	if !database.Started {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	c.String(http.StatusOK, "ok")
}
