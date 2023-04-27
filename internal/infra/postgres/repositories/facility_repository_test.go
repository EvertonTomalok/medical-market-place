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

func TestFacilityRepository(t *testing.T) {
	db, mock := mocks.NewMockEqualMatcher(t)
	defer db.Close()

	logger := logrus.New()
	repo := NewFacilityRepository(db, logger)
	query := `
		SELECT w.id, w."name", w.is_active
		FROM public."Facility" w WHERE id = $1;
	`

	facilityColumns := []string{"id", "name", "is_active"}

	t.Parallel()

	t.Run("returns error when DB fails", func(t *testing.T) {
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnError(errors.New("error"))

		facility, err := repo.FindById(context.Background(), 1)
		assert.Error(t, err)
		assert.Empty(t, facility)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("facility found", func(t *testing.T) {
		expectedFacility := entities.Facility{
			ID:       1,
			Name:     "1",
			IsActive: true,
		}
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnRows(
			sqlmock.NewRows(facilityColumns).
				AddRow(expectedFacility.ID, expectedFacility.Name, expectedFacility.IsActive),
		)

		facility, err := repo.FindById(context.Background(), expectedFacility.ID)
		assert.Nil(t, err)
		assert.Equal(t, facility, expectedFacility)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("worker not found", func(t *testing.T) {
		expectedFacility := entities.Worker{
			ID:         1,
			Name:       "1",
			IsActive:   true,
			Profession: entities.CNA,
		}
		queryMocker := mock.ExpectQuery(query).WithArgs(1)
		queryMocker.WillReturnError(sql.ErrNoRows)

		_, err := repo.FindById(context.Background(), expectedFacility.ID)
		assert.Equal(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("create default repository", func(t *testing.T) {
		_ = NewFacilityRepository(nil, logger)
	})
}
