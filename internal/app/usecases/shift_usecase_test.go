package usecases

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/app/ports/gateway"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestShiftUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	const (
		profession entities.Profession = entities.CNA
		workerID   int64               = 1
		offset     int64               = 1
		limit      int64               = 1
	)

	t.Parallel()
	t.Run("worker not found", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Worker{}, sql.ErrNoRows)

		shift, errs := ShiftUsecases{}.FindAvailableShiftsSlowerVersion(
			ctx,
			profession,
			workerID,
			nil,
			nil,
			nil,
			offset,
			limit,
		)

		if len(errs) != 1 {
			assert.Fail(t, "errrs should return just one error")
		}
		assert.Equal(t, sql.ErrNoRows, errs[0])
		assert.Empty(t, shift)
	})

	t.Run("shift gateway returns error", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(
				entities.Worker{ID: 1, Name: "name", IsActive: true, Profession: entities.CNA},
				nil,
			)

		shiftGateway := gateway.NewMockShiftGateway(ctrl)
		ShiftRepository = shiftGateway
		shiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(
				[]entities.Shift{}, errors.New("error"),
			)

		documentWorkerGateway := gateway.NewMockDocumentWorkerGateway(ctrl)
		DocumentWorkerRepository = documentWorkerGateway
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{1, 2, 3}, nil)

		shift, errs := ShiftUsecases{}.FindAvailableShiftsSlowerVersion(
			ctx,
			profession,
			workerID,
			nil,
			nil,
			nil,
			offset,
			limit,
		)

		if len(errs) == 0 {
			assert.Fail(t, "errrs should return more than one error")
		}
		assert.Empty(t, shift)
	})

	t.Run("document worker gateway error", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(
				entities.Worker{ID: 1, Name: "name", IsActive: true, Profession: entities.CNA},
				nil,
			)

		shiftGateway := gateway.NewMockShiftGateway(ctrl)
		ShiftRepository = shiftGateway
		shiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(
				[]entities.Shift{
					{
						ID:         1,
						Start:      time.Now(),
						End:        time.Now().Add(60 * time.Second),
						Profession: entities.CNA,
						IsDeleted:  false,
						Facility: entities.Facility{
							ID:       1,
							Name:     "facility",
							IsActive: true,
						},
						WorkerId: sql.NullInt64{},
					},
				},
				nil,
			)

		documentWorkerGateway := gateway.NewMockDocumentWorkerGateway(ctrl)
		DocumentWorkerRepository = documentWorkerGateway
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{}, errors.New("error"))

		shift, errs := ShiftUsecases{}.FindAvailableShiftsSlowerVersion(
			ctx,
			profession,
			workerID,
			nil,
			nil,
			nil,
			offset,
			limit,
		)

		if len(errs) == 0 {
			assert.Fail(t, "errrs should return more than one error")
		}
		assert.Empty(t, shift)
	})

	t.Run("facility requirements gateway error", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(
				entities.Worker{ID: 1, Name: "name", IsActive: true, Profession: entities.CNA},
				nil,
			)

		shiftGateway := gateway.NewMockShiftGateway(ctrl)
		ShiftRepository = shiftGateway
		shiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(
				[]entities.Shift{
					{
						ID:         1,
						Start:      time.Now(),
						End:        time.Now().Add(60 * time.Second),
						Profession: entities.CNA,
						IsDeleted:  false,
						Facility: entities.Facility{
							ID:       1,
							Name:     "facility",
							IsActive: true,
						},
						WorkerId: sql.NullInt64{},
					},
				},
				nil,
			)

		documentWorkerGateway := gateway.NewMockDocumentWorkerGateway(ctrl)
		DocumentWorkerRepository = documentWorkerGateway
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{1, 2}, nil)

		facilityRequirementsGateway := gateway.NewMockFacilityRequirementsGateway(ctrl)
		FacilityRequirementsRepository = facilityRequirementsGateway

		facilityRequirementsGateway.EXPECT().
			FindRequirementsByFacilitiesId(gomock.Any(), gomock.Any()).
			Return(
				map[int64]entities.FacilityRequirementAggregated{}, errors.New("error"),
			)

		shift, errs := ShiftUsecases{}.FindAvailableShiftsSlowerVersion(
			ctx,
			profession,
			workerID,
			nil,
			nil,
			nil,
			offset,
			limit,
		)

		if len(errs) != 1 {
			assert.Fail(t, "errrs should return just one error")
		}
		assert.Empty(t, shift)
	})

	t.Run("success returning shifts", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(
				entities.Worker{ID: 1, Name: "name", IsActive: true, Profession: entities.CNA},
				nil,
			)

		shiftGateway := gateway.NewMockShiftGateway(ctrl)
		ShiftRepository = shiftGateway
		shiftGateway.EXPECT().
			FindShifts(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			AnyTimes().
			Return(
				[]entities.Shift{
					{
						ID:         1,
						Start:      time.Now(),
						End:        time.Now().Add(60 * time.Second),
						Profession: entities.CNA,
						IsDeleted:  false,
						Facility: entities.Facility{
							ID:       1,
							Name:     "facility",
							IsActive: true,
						},
						WorkerId: sql.NullInt64{},
					},
				},
				nil,
			)

		documentWorkerGateway := gateway.NewMockDocumentWorkerGateway(ctrl)
		DocumentWorkerRepository = documentWorkerGateway
		documentWorkerGateway.EXPECT().
			FindDocumentsIds(gomock.Any(), gomock.Any()).
			Return([]int64{1, 2, 3}, nil)

		facilityRequirementsGateway := gateway.NewMockFacilityRequirementsGateway(ctrl)
		FacilityRequirementsRepository = facilityRequirementsGateway
		facilityReuirementReturn := map[int64]entities.FacilityRequirementAggregated{
			1: {
				FacilityId:  1,
				DocumentsId: []int64{1, 2},
			},
		}
		facilityRequirementsGateway.EXPECT().
			FindRequirementsByFacilitiesId(gomock.Any(), gomock.Any()).
			Return(facilityReuirementReturn, nil)

		shift, errs := ShiftUsecases{}.FindAvailableShiftsSlowerVersion(
			ctx,
			profession,
			workerID,
			nil,
			nil,
			nil,
			offset,
			limit,
		)

		if len(errs) > 0 {
			t.Fail()
		}

		assert.Equal(t, 2, len(shift))
	})
}

func TestShiftConflict(t *testing.T) {

	t.Run("shift wrap the shifTaken", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(4 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(2 * time.Hour),
				End:   now.Add(3 * time.Hour),
			},
		}

		expectedResult := true
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift start before shift taken, end in mid of the shift taken", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(4 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(2 * time.Hour),
				End:   now.Add(10 * time.Hour),
			},
		}

		expectedResult := true
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift start midle of shift taken, and end after shift taken", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(4 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(-2 * time.Hour),
				End:   now.Add(2 * time.Hour),
			},
		}

		expectedResult := true
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift has colision in start", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(20 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(-10 * time.Hour),
				End:   now.Add(25 * time.Hour),
			},
		}

		expectedResult := true
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift has colision in end", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(24 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(3 * time.Hour),
				End:   now.Add(25 * time.Hour),
			},
		}

		expectedResult := true
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift has no colision, starts after shiftTaken ends", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(24 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(30 * time.Hour),
				End:   now.Add(35 * time.Hour),
			},
		}

		expectedResult := false
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})

	t.Run("shift has no colision, ends before shiftTaken starts", func(t *testing.T) {
		now := time.Now()
		shift := entities.Shift{
			Start: now,
			End:   now.Add(1 * time.Hour),
		}

		shiftsTaken := []entities.Shift{
			{
				Start: now.Add(-2 * time.Hour),
				End:   now.Add(-1 * time.Hour),
			},
		}

		expectedResult := false
		result := ShiftUsecases{logger: logrus.New()}.shiftConflict(shift, shiftsTaken)

		assert.Equal(t, expectedResult, result)
	})
}
