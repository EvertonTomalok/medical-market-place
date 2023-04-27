package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/EvertonTomalok/marketplace-health/internal/app/helpers"
	"github.com/EvertonTomalok/marketplace-health/internal/app/usecases"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Health Marketplace API
// @Summary Get Worker By ID
// @Description When found some worker, it will returned as json object as status 200, otherwise will return the status 404 not found. If some error happens, it will return the status code 400
// @Tags WORKER
// @Router /worker/{worker_id} [get]
// @Param worker_id  path integer  false  "worker find by id"
// @Produce json
// @Success 200 {object} entities.Worker
// @Failure 404 {string} string
// @Failure 400 {object} rest.BadResponse
func GetWorkerByIdHandler(c *gin.Context) {
	ctx := c.Request.Context()
	workerIdParam := c.Param("worker_id")
	workerId, err := strconv.ParseInt(workerIdParam, 10, 64)

	if err != nil {
		errorMsg := fmt.Sprintf(
			`[ui.rest.handlers.GetWorkerById]
			Received %s as param, and it couldn't be used as integer.`,
			workerIdParam,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	log.Debugf("Fetching worker id: %d", workerId)
	worker, err := usecases.WorkerUsecases{}.FindWorker(ctx, workerId)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"[ui.rest.handlers.GetWorkerById] Error fetching worker with id %s. Error %+v.",
			workerIdParam,
			err,
		)
		helpers.SetResponseMessageError(c, errorMsg)
		return
	}

	if worker.ID == 0 {
		c.JSON(http.StatusNotFound, "Worker Not Found")
		return
	}
	c.JSON(http.StatusOK, worker)
}
