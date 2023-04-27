package gateway

import (
	"context"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
	entities "github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
)

// mockgen -source=repositories_gateway.go -destination=repositories_gateway_mock.go -package=gateway
type WorkerGateway interface {
	FindById(ctx context.Context, workerId int64) (entities.Worker, error)
}

// mockgen -source=repositories_gateway.go -destination=repositories_gateway_mock.go -package=gateway
type FacilityRequirementsGateway interface {
	FindRequirementsByFacilitiesId(
		ctx context.Context,
		facilitiesId []int64,
	) (map[int64]entities.FacilityRequirementAggregated, error)
}

// mockgen -source=repositories_gateway.go -destination=repositories_gateway_mock.go -package=gateway
type FacilityGateway interface {
	FindById(ctx context.Context, facilityId int64) (entities.Facility, error)
	FindActive(ctx context.Context) ([]entities.Facility, error)
}

// mockgen -source=repositories_gateway.go -destination=repositories_gateway_mock.go -package=gateway
type ShiftGateway interface {
	FindShifts(
		ctx context.Context,
		profession entities.Profession,
		startTime *time.Time,
		endTime *time.Time,
		workerID *int64,
		offset int64,
		limit int64,
		betweenDates *[]dto.BetweenDates,
	) ([]entities.Shift, error)
	FindShiftsByFacilities(
		ctx context.Context,
		facilitiesIds []int64,
		profession entities.Profession,
		startTime *time.Time,
		endTime *time.Time,
		offset int64,
		limit int64,
		betweenDates *[]dto.BetweenDates,
	) ([]entities.Shift, error)
}

// mockgen -source=repositories_gateway.go -destination=repositories_gateway_mock.go -package=gateway
type DocumentWorkerGateway interface {
	FindDocumentsIds(ctx context.Context, workerId int64) ([]int64, error)
}
