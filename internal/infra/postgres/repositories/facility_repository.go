package repositories

import (
	"context"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
)

type FacilityRepository struct {
	BaseRepository
}

func NewFacilityRepository(db *database.Postgres, logger *logrus.Logger) *FacilityRepository {
	path := "infra.postgres.repositories.facility_repository"
	if db == nil {
		// Using shared Conn instead db passed as param
		db = Conn
	}
	return &FacilityRepository{
		NewBaseRepository(db, logger, path),
	}
}

func (w *FacilityRepository) FindById(ctx context.Context, facilityId int64) (entities.Facility, error) {
	sqlStmt := `
		SELECT w.id, w."name", w.is_active
		FROM public."Facility" w WHERE id = $1;
	`
	w.logger.Debugf("[%s] Will run query: %s \n", w.path, sqlStmt)

	facility := entities.Facility{}
	err := w.db.QueryRowContext(ctx, sqlStmt, facilityId).Scan(&facility.ID, &facility.Name, &facility.IsActive)
	if err != nil {
		w.logger.Debugf("[%s] error running query %+v \n", w.path, err)
		return entities.Facility{}, err
	}
	return facility, nil
}

func (w *FacilityRepository) FindActive(ctx context.Context) ([]entities.Facility, error) {
	sqlStmt := `SELECT f.id, f."name" , f.is_active from public."Facility" f WHERE f.is_active = true;`
	rows, err := w.db.QueryContext(ctx, sqlStmt)
	if err != nil {
		w.logger.Errorf("[%s] Error making query: %+v \n", w.path, err)
		return []entities.Facility{}, err
	}
	defer rows.Close()

	var facilities []entities.Facility
	for rows.Next() {
		facility := entities.Facility{}
		if err := rows.Scan(&facility.ID, &facility.Name, &facility.IsActive); err != nil {
			w.logger.Errorf("[%s] Error scanning query: %+v \n", w.path, err)
			return []entities.Facility{}, err
		}

		facilities = append(facilities, facility)
	}
	return facilities, nil
}
