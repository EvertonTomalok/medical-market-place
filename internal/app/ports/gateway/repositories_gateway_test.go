package gateway

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRepositoriesGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const workerID int64 = 1
	ctx := context.Background()

	t.Run("worker gateway error", func(t *testing.T) {
		workerGateway := NewMockWorkerGateway(ctrl)
		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Worker{}, errors.New("error"))

		worker, err := workerGateway.FindById(ctx, workerID)
		assert.Error(t, err)
		assert.Empty(t, worker)
	})

	t.Run("return valid worker", func(t *testing.T) {
		workerEntityMock := entities.Worker{
			ID:         1,
			Name:       "worker 1",
			IsActive:   true,
			Profession: entities.CNA,
		}
		workerGateway := NewMockWorkerGateway(ctrl)
		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(workerEntityMock, nil)

		worker, err := workerGateway.FindById(ctx, workerID)
		assert.NoError(t, err)
		assert.Equal(t, worker, workerEntityMock)
	})
}

func TestFacilityGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const facilityID int64 = 1
	ctx := context.Background()

	t.Run("facility gateway error", func(t *testing.T) {
		facilityGateway := NewMockFacilityGateway(ctrl)
		facilityGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Facility{}, errors.New("error"))

		facility, err := facilityGateway.FindById(ctx, facilityID)
		assert.Error(t, err)
		assert.Empty(t, facility)
	})

	t.Run("return valid facility", func(t *testing.T) {
		facilityEntityMock := entities.Facility{
			ID:       1,
			Name:     "worker 1",
			IsActive: true,
		}
		facilityGateway := NewMockFacilityGateway(ctrl)
		facilityGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(facilityEntityMock, nil)

		facility, err := facilityGateway.FindById(ctx, facilityID)
		assert.NoError(t, err)
		assert.Equal(t, facility, facilityEntityMock)
	})
}

func TestShiftGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("shift gateway error", func(t *testing.T) {
		ShiftGateway := NewMockShiftGateway(ctrl)
		ShiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return([]entities.Shift{}, errors.New("error"))

		shifts, err := ShiftGateway.FindShifts(ctx, entities.CNA, nil, nil, nil, 0, 100, nil)
		assert.Error(t, err)
		assert.Empty(t, shifts)
	})

	t.Run("return valid shift", func(t *testing.T) {
		now := time.Now()
		shiftsMock := []entities.Shift{
			{
				ID:         1,
				Start:      now,
				End:        now,
				Profession: entities.CNA,
				IsDeleted:  false,
				Facility: entities.Facility{
					ID:       1,
					Name:     "some name",
					IsActive: true,
				},
			},
		}
		ShiftGateway := NewMockShiftGateway(ctrl)
		ShiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(shiftsMock, nil)

		shifts, err := ShiftGateway.FindShifts(ctx, entities.CNA, nil, nil, nil, 0, 100, nil)
		assert.NoError(t, err)
		assert.Equal(t, shiftsMock, shifts)
	})
}

func TestDocumentWorkerGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("document worker gateway error", func(t *testing.T) {
		documentWorkerGateway := NewMockDocumentWorkerGateway(ctrl)
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{}, errors.New("error"))

		shifts, err := documentWorkerGateway.FindDocumentsIds(ctx, 1)
		assert.Error(t, err)
		assert.Empty(t, shifts)
	})

	t.Run("return valid shift", func(t *testing.T) {
		expectedResult := []int64{1, 2, 3}
		documentWorkerGateway := NewMockDocumentWorkerGateway(ctrl)
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{1, 2, 3}, nil)

		documentWorker, err := documentWorkerGateway.FindDocumentsIds(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, documentWorker)
	})
}

func TestFacilityRequirementsGateway(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("document worker gateway error", func(t *testing.T) {
		facilityRequirementsGateway := NewMockFacilityRequirementsGateway(ctrl)
		facilityRequirementsGateway.EXPECT().
			FindRequirementsByFacilitiesId(gomock.Any(), gomock.Any()).
			Return(map[int64]entities.FacilityRequirementAggregated{}, errors.New("error"))

		facilitiesRequirement, err := facilityRequirementsGateway.FindRequirementsByFacilitiesId(ctx, []int64{1})
		assert.Error(t, err)
		assert.Empty(t, facilitiesRequirement)
	})

	t.Run("return valid shift", func(t *testing.T) {
		facilityRequirementsGateway := NewMockFacilityRequirementsGateway(ctrl)

		facilityReuirementReturn := map[int64]entities.FacilityRequirementAggregated{
			1: {
				FacilityId:  1,
				DocumentsId: []int64{1, 2},
			},
		}
		facilityRequirementsGateway.EXPECT().
			FindRequirementsByFacilitiesId(gomock.Any(), gomock.Any()).
			Return(facilityReuirementReturn, nil)

		facilitiesRequirement, err := facilityRequirementsGateway.FindRequirementsByFacilitiesId(ctx, []int64{1})
		assert.NoError(t, err)
		assert.Equal(t, facilityReuirementReturn, facilitiesRequirement)
	})
}
