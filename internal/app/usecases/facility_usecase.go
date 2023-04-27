package usecases

import (
	"context"
	"database/sql"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
)

type FacilityUsecase struct{}

func (fu FacilityUsecase) FindFacility(ctx context.Context, facilityId int64) (entities.Facility, error) {
	facility, err := FacilityRepository.FindById(ctx, facilityId)

	switch err {
	case nil:
		return facility, nil
	case sql.ErrNoRows:
		return entities.Facility{}, nil
	default:
		return facility, err
	}

}
