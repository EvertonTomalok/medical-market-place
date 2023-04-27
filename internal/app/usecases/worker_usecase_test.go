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

func TestWorkerUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	t.Run("worker gateway error", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Worker{}, errors.New("error"))

		worker, err := WorkerUsecases{}.FindWorker(ctx, 1)
		assert.Error(t, err)
		assert.Empty(t, worker)
	})

	t.Run("worker not found", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(entities.Worker{}, sql.ErrNoRows)

		worker, err := WorkerUsecases{}.FindWorker(ctx, 1)
		assert.NoError(t, err)
		assert.Empty(t, worker)
	})

	t.Run("valid return", func(t *testing.T) {
		workerGateway := gateway.NewMockWorkerGateway(ctrl)
		WorkerRepository = workerGateway

		expected := entities.Worker{
			ID:         1,
			Name:       "name",
			IsActive:   true,
			Profession: entities.CNA,
		}

		workerGateway.EXPECT().
			FindById(gomock.Any(), gomock.Any()).
			Return(expected, nil)

		worker, err := WorkerUsecases{}.FindWorker(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, expected, worker)
	})

}
