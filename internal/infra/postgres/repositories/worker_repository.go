package repositories

import (
	"context"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
)

type WorkerRepository struct {
	BaseRepository
}

func NewWorkerRepository(db *database.Postgres, logger *logrus.Logger) *WorkerRepository {
	path := "infra.postgres.repositories.worker_repository"
	if db == nil {
		// Using shared Conn instead db passed as param
		db = Conn
	}
	return &WorkerRepository{
		NewBaseRepository(db, logger, path),
	}
}

func (w *WorkerRepository) FindById(ctx context.Context, workerId int64) (entities.Worker, error) {
	sqlStmt := `
		SELECT w.id, w."name", w.is_active, w."profession"
		FROM public."Worker" w WHERE id = $1;
	`
	w.logger.Debugf("[%s] Will run query: %s \n\n\n\nWith param %d \n\n\n\n", w.path, sqlStmt, workerId)

	worker := entities.Worker{}
	err := w.db.QueryRowContext(ctx, sqlStmt, workerId).Scan(&worker.ID, &worker.Name, &worker.IsActive, &worker.Profession)
	if err != nil {
		w.logger.Errorf("[%s] error running query %+v \n", w.path, err)
		return entities.Worker{}, err
	}
	return worker, nil
}
