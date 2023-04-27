package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/sirupsen/logrus"
)

type ShiftRepository struct {
	BaseRepository
}

func NewShiftRepository(db *database.Postgres, logger *logrus.Logger) *ShiftRepository {
	path := "infra.postgres.repositories.shift_repository"
	if db == nil {
		// Using shared Conn instead db passed as param
		db = Conn
	}
	return &ShiftRepository{
		NewBaseRepository(db, logger, path),
	}
}

func (sr *ShiftRepository) FindShifts(
	ctx context.Context,
	profession entities.Profession,
	startTime *time.Time,
	endTime *time.Time,
	workerId *int64,
	offset int64,
	limit int64,
	betweenDates *[]dto.BetweenDates,
) ([]entities.Shift, error) {
	// IMPROVEMENT POSSIBLE
	// CREATE paginator

	sqlStmt, params := sr.makeQueryBuilder(
		nil,
		profession,
		startTime,
		endTime,
		workerId,
		offset,
		limit,
		betweenDates,
	)

	sr.logger.Debugf("[%s] Will run query: %s \n\nWith params %+v \n\n\n\n", sr.path, sqlStmt, params)

	rows, err := sr.db.QueryContext(ctx, sqlStmt, params...)
	if err != nil {
		sr.logger.Errorf("[%s] Error making query: %+v \n", sr.path, err)
		return []entities.Shift{}, err
	}
	defer rows.Close()

	var shifts []entities.Shift
	for rows.Next() {
		shift := entities.Shift{}
		err = rows.Scan(
			&shift.ID,
			&shift.Start,
			&shift.End,
			&shift.Profession,
			&shift.IsDeleted,
			&shift.WorkerId,
			&shift.Facility.ID,
			&shift.Facility.Name,
			&shift.Facility.IsActive,
		)
		if err != nil {
			sr.logger.Errorf("[%s] Error scanning query: %+v \n", sr.path, err)
			return shifts, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (sr ShiftRepository) FindShiftsByFacilities(
	ctx context.Context,
	facilitiesIds []int64,
	profession entities.Profession,
	startTime *time.Time,
	endTime *time.Time,
	offset int64,
	limit int64,
	betweenDates *[]dto.BetweenDates,
) ([]entities.Shift, error) {
	sqlStmt, params := sr.makeQueryBuilder(
		&facilitiesIds,
		profession,
		startTime,
		endTime,
		nil,
		offset,
		limit,
		betweenDates,
	)

	sr.logger.Debugf("[%s] Will run query: %s \n\nWith params %+v \n\n\n\n", sr.path, sqlStmt, params)

	rows, err := sr.db.QueryContext(ctx, sqlStmt, params...)
	if err != nil {
		sr.logger.Errorf("[%s] Error making query: %+v \n", sr.path, err)
		return []entities.Shift{}, err
	}
	defer rows.Close()

	var shifts []entities.Shift
	for rows.Next() {
		shift := entities.Shift{}
		err = rows.Scan(
			&shift.ID,
			&shift.Start,
			&shift.End,
			&shift.Profession,
			&shift.IsDeleted,
			&shift.WorkerId,
			&shift.Facility.ID,
			&shift.Facility.Name,
			&shift.Facility.IsActive,
		)
		if err != nil {
			sr.logger.Errorf("[%s] Error scanning query: %+v \n", sr.path, err)
			return shifts, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (sr ShiftRepository) makeQueryBuilder(
	facilitiesIds *[]int64,
	profession entities.Profession,
	startTime *time.Time,
	endTime *time.Time,
	workerId *int64,
	offset int64,
	limit int64,
	betweenDates *[]dto.BetweenDates,
) (string, []interface{}) {

	var sb strings.Builder
	var params = []interface{}{}
	var anchorNum int = 1 // pointer to control anchors in query

	// Swapping dates
	if startTime != nil && endTime != nil {
		if endTime.Before(*startTime) {
			sr.logger.Debugf("[%s] Date start %+v and end %+v will be changed \n", sr.path, startTime, endTime)
			startTime, endTime = endTime, startTime
		}

		if startTime.Day() == endTime.Day() && startTime.Month() == endTime.Month() && startTime.Year() == endTime.Year() {
			// same date
			sr.logger.Debugf("[%s] Date end %+v will be increased in one day \n", sr.path, endTime)
			twentyFourHours := 24 * time.Hour
			*endTime = endTime.Add(twentyFourHours)
		}
	}

	if facilitiesIds == nil {
		sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, facilities."name" , facilities.is_active
		FROM (SELECT f.id, f."name" , f.is_active from public."Facility" f where f.is_active = true) as facilities
		INNER JOIN public."Shift" s ON facilities.id = s.facility_id
		WHERE s.is_deleted = false
		`)
	} else {
		sb.WriteString(`
		SELECT s.id, s."start", s."end", s."profession", s."is_deleted", s.worker_id, s.facility_id, f."name" , f.is_active
		FROM public."Shift" s
		INNER JOIN public."Facility" f on f.id = s.facility_id
		WHERE s.is_deleted = false 
		`)

		if len(*facilitiesIds) > 0 { // Only add this filter if receive some facility ID
			facilitiesIdAsArrayString := []string{}
			for _, id := range *facilitiesIds {
				facilitiesIdAsArrayString = append(facilitiesIdAsArrayString, fmt.Sprintf("%d", id))
			}
			param := "{" + strings.Join(facilitiesIdAsArrayString, ",") + "}"
			sb.WriteString(fmt.Sprintf(`AND f.id = ANY($%d::int[]) `, anchorNum))
			params = append(params, param)
			anchorNum++ // increasing anchor num
		}
	}

	sb.WriteString(fmt.Sprintf(`AND s.profession = $%d `, anchorNum))
	params = append(params, profession.ToString())
	anchorNum++ // increasing anchor num

	if workerId != nil {
		sb.WriteString(fmt.Sprintf(`AND s."worker_id" = $%d `, anchorNum))
		params = append(params, *workerId)
		anchorNum++
	} else {
		sb.WriteString(`AND s."worker_id" IS NULL `)
	}

	if startTime != nil {
		sb.WriteString(fmt.Sprintf(`AND s."start" >= $%d `, anchorNum))
		params = append(params, *startTime)
		anchorNum++
	}

	if endTime != nil {
		sb.WriteString(fmt.Sprintf(`AND s."end" <= $%d `, anchorNum))
		params = append(params, *endTime)
		anchorNum++
	}

	if betweenDates != nil {
		for _, dates := range *betweenDates {
			sb.WriteString(fmt.Sprintf(`OR (s."start" < $%d `, anchorNum))
			params = append(params, dates.Start)
			anchorNum++
			sb.WriteString(fmt.Sprintf(`AND s."end" > $%d) `, anchorNum))
			params = append(params, dates.End)
			anchorNum++
		}
	}

	sb.WriteString(`order by s."start" `)
	sb.WriteString(fmt.Sprintf(`OFFSET $%d `, anchorNum))
	params = append(params, offset)
	anchorNum++
	if limit > 0 {
		sb.WriteString(fmt.Sprintf(`LIMIT $%d `, anchorNum))
		params = append(params, limit)
	}

	return sb.String(), params
}
