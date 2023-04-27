package usecases

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/EvertonTomalok/marketplace-health/internal/app/ports/gateway"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFacilityUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("facility gateway error", func(t *testing.T) {
		facilityGateway := gateway.NewMockFacilityGateway(ctrl)
		FacilityRepository = facilityGateway

		facilityGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Facility{}, errors.New("error"))

		facility, err := FacilityUsecase{}.FindFacility(ctx, 1)
		assert.Error(t, err)
		assert.Empty(t, facility)
	})

	t.Run("facility found", func(t *testing.T) {
		facilityGateway := gateway.NewMockFacilityGateway(ctrl)
		FacilityRepository = facilityGateway

		expectedFacility := entities.Facility{
			ID:       1,
			Name:     "name",
			IsActive: true,
		}

		facilityGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(expectedFacility, nil)

		facility, err := FacilityUsecase{}.FindFacility(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedFacility, facility)
	})

	t.Run("facility not found", func(t *testing.T) {
		facilityGateway := gateway.NewMockFacilityGateway(ctrl)
		FacilityRepository = facilityGateway

		facilityGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Facility{}, sql.ErrNoRows)

		facility, err := FacilityUsecase{}.FindFacility(ctx, 1)

		assert.NoError(t, err)
		assert.Empty(t, facility)
	})

}
