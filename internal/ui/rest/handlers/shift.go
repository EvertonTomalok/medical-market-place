package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/app/helpers"
	"github.com/EvertonTomalok/marketplace-health/internal/app/usecases"
	restDTO "github.com/EvertonTomalok/marketplace-health/internal/domain/dto/rest"
	"github.com/gin-gonic/gin"
)

// Health Marketplace API
// @Summary Get Available shifts
// @Description It will list all shifts available for some worker, grouped by date. If some error happens, it will return the status code 400
// @Tags SHIFTS
// @Router /shifts/available/{worker_id}/{profession} [get]
// @Param worker_id  path integer true "worker id to find"
// @Param profession path string true "The worker profession, must be 'CNA', 'LVW' or 'RN'."
// @Param offset query integer false "offset to start search"
// @Param limit query integer false "limit of return values"
// @Param start query string false "start date like '2006-01-02'"
// @Param end query string false "end date like '2006-01-02'"
// @Produce json
// @Success 200 {object} dto.GroupedByDateShift
// @Failure 400 {object} rest.BadResponse
// @Failure 404 {object} rest.BadResponse
func GetAvailableShiftsHandler(c *gin.Context) {
	ctx := c.Request.Context()
	path := "ui.rest.handlers.GetAvailableShiftsHandler"

	params := restDTO.ShiftAvailablePathDTO{}
	err := c.BindUri(&params)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"[%s] Error fetching shifts. Error %+v.",
			path,
			err,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	if !params.Profession.IsValid() {
		errorMsg := fmt.Sprintf(
			"[%s] Profession '%s' not valid. The worker profession, must be 'CNA', 'LVW' or 'RN'.",
			path,
			params.Profession,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	now := time.Now()
	sevenDays := time.Hour * 24 * 7
	fiveMonths := time.Hour * 24 * 30 * 5
	queries := restDTO.ShiftAvailableQueryDTO{
		Offset: 0,
		Limit:  0,
		Start:  now.Add(-fiveMonths * 2),
		End:    now.Add(sevenDays),
	}
	err = c.ShouldBindQuery(&queries) // if bind, will overload the default values
	if err != nil {
		errorMsg := fmt.Sprintf(
			"[%s] Error fetching shifts. Error %+v.",
			path,
			err,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}
	shifts, errList := usecases.ShiftUsecases{}.FindAvailableShiftsFasterVersion(
		ctx,
		params.Profession,
		params.WorkerID,
		nil,
		&queries.Start,
		&queries.End,
		queries.Offset,
		queries.Limit,
	)
	if len(errList) > 0 {
		errorsString := make([]string, 0)

		for _, e := range errList {
			errorsString = append(errorsString, e.Error())
		}

		errorMsg := fmt.Sprintf(
			"[%s] Error fetching shifts. Error: %+v",
			path,
			strings.Join(errorsString, " | "),
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	groupedShifts := helpers.GroupShiftByDate(shifts)

	c.JSON(http.StatusOK, groupedShifts)
}
