package helpers

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto/rest"
	"github.com/gin-gonic/gin"
)

func SetResponseMessageError(c *gin.Context, errorMsg string) {
	log.Error(errorMsg)
	c.JSON(http.StatusBadRequest, rest.BadResponse{ErrorMessage: errorMsg})
}
