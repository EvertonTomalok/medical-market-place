package rest

import (
	"net/http"

	"github.com/EvertonTomalok/marketplace-health/internal/ui/rest/handlers"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Path    string
	Method  string
	Handler gin.HandlerFunc
}

var healthCheck = []Route{
	{
		"/health",
		http.MethodGet,
		handlers.Health,
	},
	{
		"/readiness",
		http.MethodGet,
		handlers.Readiness,
	},
}

var routes = []Route{
	{
		"/worker/:worker_id",
		http.MethodGet,
		handlers.GetWorkerByIdHandler,
	},
	{
		"/facility/:facility_id",
		http.MethodGet,
		handlers.GetFacilityByIdHandler,
	},
	{
		"/shift/available/:worker_id/:profession",
		http.MethodGet,
		handlers.GetAvailableShiftsHandler,
	},
}
