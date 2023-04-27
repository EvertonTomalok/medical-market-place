package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EvertonTomalok/marketplace-health/internal/app/helpers"
	"github.com/EvertonTomalok/marketplace-health/internal/app/usecases"
	"github.com/gin-gonic/gin"
)

// Health Marketplace API
// @Summary Get Facility By ID
// @Description When found some facility, it will returned as json object as status 200, otherwise will return the status 404 not found. If some error happens, it will return the status code 400
// @Tags FACILITY
// @Router /facility/{facility_id} [get]
// @Param facility_id  path integer  false  "facility find by id"
// @Produce json
// @Success 200 {object} entities.Facility
// @Failure 400 {object} rest.BadResponse
// @Failure 404 {string} string
func GetFacilityByIdHandler(c *gin.Context) {
	ctx := c.Request.Context()
	path := "ui.rest.handlers.Facility"

	facilityIdParam := c.Param("facility_id")
	facilityId, err := strconv.ParseInt(facilityIdParam, 10, 64)

	if err != nil {
		errorMsg := fmt.Sprintf(
			`[%s] Received %s as param, and it couldn't be used as integer.`, path, facilityIdParam,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	facility, err := usecases.FacilityUsecase{}.FindFacility(ctx, facilityId)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"[%s] Error fetching worker with id %s. Error %+v.", path, facilityIdParam, err,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	if facility.ID == 0 {
		c.JSON(http.StatusNotFound, "Facility not found.")
		return
	}
	c.JSON(http.StatusOK, facility)
}
