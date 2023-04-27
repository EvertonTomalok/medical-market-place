package usecases

import (
	"context"
	"database/sql"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
)

type WorkerUsecases struct {
	Logger *logrus.Logger
}

func (wu WorkerUsecases) FindWorker(ctx context.Context, workerId int64) (entities.Worker, error) {
	worker, err := WorkerRepository.FindById(ctx, workerId)

	switch err {
	case nil:
		return worker, nil
	case sql.ErrNoRows:
		return entities.Worker{}, nil
	default:
		return worker, err
	}

}
