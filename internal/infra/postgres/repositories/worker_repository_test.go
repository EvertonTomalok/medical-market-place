package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database/mocks"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWorkerRepository(t *testing.T) {
	db, mock := mocks.NewMockEqualMatcher(t)
	defer db.Close()

	logger := logrus.New()
	repo := NewWorkerRepository(db, logger)
	query := `
		SELECT w.id, w."name", w.is_active, w."profession"
		FROM public."Worker" w WHERE id = $1;
	`

	workerColumns := []string{"id", "name", "is_active", "profession"}

	t.Parallel()

	t.Run("returns error when DB fails", func(t *testing.T) {
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnError(errors.New("error"))

		worker, err := repo.FindById(context.Background(), 1)
		assert.Error(t, err)
		assert.Empty(t, worker)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("worker found", func(t *testing.T) {
		expectedWorker := entities.Worker{
			ID:         1,
			Name:       "1",
			IsActive:   true,
			Profession: entities.CNA,
		}
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnRows(
			sqlmock.NewRows(workerColumns).
				AddRow(expectedWorker.ID, expectedWorker.Name, expectedWorker.IsActive, expectedWorker.Profession),
		)

		worker, err := repo.FindById(context.Background(), expectedWorker.ID)
		assert.Nil(t, err)
		assert.Equal(t, worker, expectedWorker)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("worker not found", func(t *testing.T) {
		expectedWorker := entities.Worker{
			ID:         1,
			Name:       "1",
			IsActive:   true,
			Profession: entities.CNA,
		}
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnError(sql.ErrNoRows)

		_, err := repo.FindById(context.Background(), expectedWorker.ID)
		assert.Equal(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("create default repository", func(t *testing.T) {
		_ = NewWorkerRepository(nil, logger)
	})
}
