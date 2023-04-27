package repositories

import (
	"context"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/sirupsen/logrus"
)

type DocumentWorkerRepository struct {
	BaseRepository
}

func NewDocumentWorkerRepository(db *database.Postgres, logger *logrus.Logger) *DocumentWorkerRepository {
	path := "infra.postgres.repositories.document_worker_repository"
	if db == nil {
		// Using shared Conn instead db passed as param
		db = Conn
	}
	return &DocumentWorkerRepository{
		NewBaseRepository(db, logger, path),
	}
}

func (dwr *DocumentWorkerRepository) FindDocumentsIds(ctx context.Context, workerId int64) ([]int64, error) {
	sqlStmt := `SELECT DISTINCT dw.document_id FROM public."DocumentWorker" dw WHERE dw.worker_id = $1;`
	dwr.logger.Debugf("[%s] Will run query: %s \n", dwr.path, sqlStmt)

	rows, err := dwr.db.QueryContext(ctx, sqlStmt, workerId)
	if err != nil {
		dwr.logger.Errorf("[%s] Error making query: %+v \n", dwr.path, err)
		return []int64{}, err
	}

	documentsIds := make([]int64, 0)
	defer rows.Close()

	for rows.Next() {
		var documentId int64
		err = rows.Scan(&documentId)

		if err != nil {
			dwr.logger.Errorf("[%s] Error scanning query: %+v \n", dwr.path, err)
			return []int64{}, err
		}

		documentsIds = append(documentsIds, documentId)
	}

	return documentsIds, nil
}
