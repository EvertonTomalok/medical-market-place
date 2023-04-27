package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database/mocks"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDocumentWorkerRepository(t *testing.T) {
	db, mock := mocks.NewMockEqualMatcher(t)
	defer db.Close()

	logger := logrus.New()
	repo := NewDocumentWorkerRepository(db, logger)
	query := `SELECT DISTINCT dw.document_id FROM public."DocumentWorker" dw WHERE dw.worker_id = $1;`

	facilityColumns := []string{"document_id"}
	const workerID int64 = 1

	t.Parallel()

	t.Run("returns error when DB fails", func(t *testing.T) {
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnError(errors.New("error"))

		workerDocument, err := repo.FindDocumentsIds(context.Background(), workerID)
		assert.Error(t, err)
		assert.Empty(t, workerDocument)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("returns error when DB return wrong value", func(t *testing.T) {
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnRows(
			sqlmock.NewRows(facilityColumns).
				AddRow("invalid return"),
		)

		workerDocument, err := repo.FindDocumentsIds(context.Background(), workerID)
		assert.Error(t, err)
		assert.Empty(t, workerDocument)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("document worker found", func(t *testing.T) {
		documentWorker := entities.DocumentWorker{
			ID:         1,
			WorkerId:   1,
			DocumentId: 1,
		}
		expectIds := []int64{1}
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnRows(
			sqlmock.NewRows(facilityColumns).
				AddRow(documentWorker.DocumentId),
		)

		ids, err := repo.FindDocumentsIds(context.Background(), documentWorker.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectIds, ids)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("create default repository", func(t *testing.T) {
		_ = NewDocumentWorkerRepository(nil, logger)
	})
}
