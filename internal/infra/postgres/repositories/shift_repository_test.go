package repositories

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database/mocks"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestMakeQueryBuilder(t *testing.T) {
	logger := logrus.Logger{}
	shiftRepository := ShiftRepository{
		NewBaseRepository(nil, &logger, "test"),
	}
	profession := entities.CNA
	const (
		offset int64 = 0
		limit  int64 = 10
	)

	var (
		workerId int64     = 1
		start    time.Time = time.Now()
		end      time.Time = time.Now()
	)

	t.Parallel()
	t.Run("create query without any optional param", func(t *testing.T) {
		var sb strings.Builder

		sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, facilities."name" , facilities.is_active
		FROM (SELECT f.id, f."name" , f.is_active from public."Facility" f where f.is_active = true) as facilities
		INNER JOIN public."Shift" s ON facilities.id = s.facility_id
		WHERE s.is_deleted = false
		`)
		sb.WriteString(fmt.Sprintf(`AND s.profession = $%d `, 1))
		sb.WriteString(`AND s."worker_id" IS NULL `)
		sb.WriteString(`order by s."start" `)
		sb.WriteString(fmt.Sprintf(`OFFSET $%d `, 2))
		sb.WriteString(fmt.Sprintf(`LIMIT $%d `, 3))
		expectedQuery := sb.String()
		query, params := shiftRepository.makeQueryBuilder(
			nil,
			profession,
			nil,
			nil,
			nil,
			offset,
			limit,
			nil,
		)
		assert.Equal(t, query, expectedQuery)
		assert.Equal(t, len(params), 3)
	})

	t.Run("create query only with worker id", func(t *testing.T) {
		var sb strings.Builder

		sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, facilities."name" , facilities.is_active
		FROM (SELECT f.id, f."name" , f.is_active from public."Facility" f where f.is_active = true) as facilities
		INNER JOIN public."Shift" s ON facilities.id = s.facility_id
		WHERE s.is_deleted = false
		`)
		sb.WriteString(fmt.Sprintf(`AND s.profession = $%d `, 1))
		sb.WriteString(fmt.Sprintf(`AND s."worker_id" = $%d `, 2))
		sb.WriteString(`order by s."start" `)
		sb.WriteString(fmt.Sprintf(`OFFSET $%d `, 3))
		sb.WriteString(fmt.Sprintf(`LIMIT $%d `, 4))
		expectedQuery := sb.String()
		query, params := shiftRepository.makeQueryBuilder(
			nil,
			profession,
			nil,
			nil,
			&workerId,
			offset,
			limit,
			nil,
		)
		assert.Equal(t, query, expectedQuery)
		assert.Equal(t, len(params), 4)
	})

	t.Run("create query only with start and end time and limit equal 0", func(t *testing.T) {
		var sb strings.Builder

		sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, facilities."name" , facilities.is_active
		FROM (SELECT f.id, f."name" , f.is_active from public."Facility" f where f.is_active = true) as facilities
		INNER JOIN public."Shift" s ON facilities.id = s.facility_id
		WHERE s.is_deleted = false
		`)
		sb.WriteString(fmt.Sprintf(`AND s.profession = $%d `, 1))
		sb.WriteString(`AND s."worker_id" IS NULL `)
		sb.WriteString(fmt.Sprintf(`AND s."start" >= $%d `, 2))
		sb.WriteString(fmt.Sprintf(`AND s."end" <= $%d `, 3))
		sb.WriteString(`order by s."start" `)
		sb.WriteString(fmt.Sprintf(`OFFSET $%d `, 4))
		expectedQuery := sb.String()
		query, params := shiftRepository.makeQueryBuilder(
			nil,
			profession,
			&start,
			&end,
			nil,
			offset,
			0, // limit 0 won't be included in query
			nil,
		)
		assert.Equal(t, query, expectedQuery)
		assert.Equal(t, len(params), 4)
	})
}

func TestShiftRepository(t *testing.T) {
	logger := logrus.New()
	profession := entities.CNA
	const (
		offset int64 = 0
		limit  int64 = 10
	)

	var (
		start time.Time = time.Now()
		end   time.Time = time.Now()
	)

	var sb strings.Builder

	sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, facilities."name" , facilities.is_active
		FROM (SELECT f.id, f."name" , f.is_active from public."Facility" f where f.is_active = true) as facilities
		INNER JOIN public."Shift" s ON facilities.id = s.facility_id
		WHERE s.is_deleted = false
		`)
	sb.WriteString(fmt.Sprintf(`AND s.profession = $%d `, 1))
	sb.WriteString(`AND s."worker_id" IS NULL `)
	sb.WriteString(fmt.Sprintf(`AND s."start" >= $%d `, 2))
	sb.WriteString(fmt.Sprintf(`AND s."end" <= $%d `, 3))
	sb.WriteString(`order by s."start" `)
	sb.WriteString(fmt.Sprintf(`OFFSET $%d `, 4))
	sb.WriteString(fmt.Sprintf(`LIMIT $%d `, 5))
	query := sb.String()

	shiftColumns := []string{
		"id", "start", "end", "profession", "is_deleted",
		"worker_id", "facility_id", "name",
		"is_active",
	}

	t.Parallel()

	t.Run("success retrieve data", func(t *testing.T) {
		db, mock := mocks.NewMockEqualMatcher(t)
		defer db.Close()

		expectedShift := entities.Shift{
			ID:         1,
			Start:      start,
			End:        end,
			Profession: profession,
			IsDeleted:  false,
			WorkerId:   sql.NullInt64{},
			Facility: entities.Facility{
				ID:       1,
				Name:     "name",
				IsActive: true,
			},
		}

		repo := NewShiftRepository(db, logger)
		queryMocker := mock.ExpectQuery(query).WithArgs(profession.ToString(), AnyTime{}, AnyTime{}, offset, limit)
		queryMocker.WillReturnRows(
			sqlmock.NewRows(shiftColumns).
				AddRow(
					expectedShift.ID,
					expectedShift.Start,
					expectedShift.End,
					expectedShift.Profession,
					expectedShift.IsDeleted,
					expectedShift.WorkerId,
					expectedShift.Facility.ID,
					expectedShift.Facility.Name,
					expectedShift.Facility.IsActive,
				),
		)

		ctx := context.Background()
		shift, err := repo.FindShifts(ctx, profession, &start, &end, nil, offset, limit, nil)

		assert.NoError(t, err)
		assert.Equal(t, []entities.Shift{expectedShift}, shift)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("success retrieve data but end less than start date", func(t *testing.T) {
		db, mock := mocks.NewMockRegexMatcher(t)
		defer db.Close()

		end = time.Now()
		start = time.Now().Add(10 * time.Second)

		expectedShift := entities.Shift{
			ID:         1,
			Start:      end,
			End:        start,
			Profession: profession,
			IsDeleted:  false,
			WorkerId:   sql.NullInt64{},
			Facility: entities.Facility{
				ID:       1,
				Name:     "name",
				IsActive: true,
			},
		}

		repo := NewShiftRepository(db, logger)
		queryMocker := mock.ExpectQuery(".*").WithArgs(profession.ToString(), AnyTime{}, AnyTime{}, offset, limit) // start and end will be shifted
		queryMocker.WillReturnRows(
			sqlmock.NewRows(shiftColumns).
				AddRow(
					expectedShift.ID,
					end,
					start,
					profession,
					expectedShift.IsDeleted,
					expectedShift.WorkerId,
					expectedShift.Facility.ID,
					expectedShift.Facility.Name,
					expectedShift.Facility.IsActive,
				),
		)

		ctx := context.Background()
		_, err := repo.FindShifts(ctx, profession, &start, &end, nil, offset, limit, nil)

		assert.NoError(t, err)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("returns error when DB fails", func(t *testing.T) {
		db, mock := mocks.NewMockRegexMatcher(t)
		defer db.Close()

		start = time.Now()
		end = time.Now().Add(10 * time.Second)
		repo := NewShiftRepository(db, logger)
		queryMocker := mock.ExpectQuery(".*").WithArgs(profession.ToString(), AnyTime{}, AnyTime{}, offset, limit)
		queryMocker.WillReturnError(errors.New("error"))

		facility, err := repo.FindShifts(context.Background(), profession, &start, &end, nil, offset, limit, nil)
		assert.Error(t, err)
		assert.Empty(t, facility)

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestEmptyConnCreator(t *testing.T) {
	_ = NewShiftRepository(nil, logrus.New())
}
