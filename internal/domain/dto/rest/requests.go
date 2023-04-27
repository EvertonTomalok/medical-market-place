package rest

import (
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
)

type ShiftAvailablePathDTO struct {
	WorkerID   int64               `uri:"worker_id" binding:"required,gte=1"`
	Profession entities.Profession `uri:"profession" binding:"required"`
}

type ShiftAvailableQueryDTO struct {
	Offset     int64     `form:"offset"`
	Limit      int64     `form:"limit"`
	FacilityID []int64   `form:"facility_id"`
	Start      time.Time `form:"start" time_format:"2006-01-02" time_utc:"1"`
	End        time.Time `form:"end" time_format:"2006-01-02" time_utc:"1"`
}
