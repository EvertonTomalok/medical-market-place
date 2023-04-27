package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/database"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type FacilityRequirementRepository struct {
	BaseRepository
}

func NewFacilityRequirementRepository(db *database.Postgres, logger *logrus.Logger) *FacilityRequirementRepository {
	path := "infra.postgres.repositories.facility_requirement_repository"
	if db == nil {
		// Using shared Conn instead db passed as param
		db = Conn
	}
	return &FacilityRequirementRepository{
		NewBaseRepository(db, logger, path),
	}
}

func (frr *FacilityRequirementRepository) FindRequirementsByFacilitiesId(
	ctx context.Context,
	facilitiesId []int64,
) (map[int64]entities.FacilityRequirementAggregated, error) {
	sqlStmt := `
	SELECT fr.facility_id, array_agg(distinct fr.document_id) as documents_id  FROM  public."FacilityRequirement" fr
	WHERE fr.facility_id = ANY($1::int[])
	GROUP BY fr.facility_id;
	`
	seen := make(map[int64]bool) // use this structure to not put duplicated values in list
	facilitiesIdAsArrayString := []string{}
	for _, id := range facilitiesId {
		if _, ok := seen[id]; !ok {
			facilitiesIdAsArrayString = append(facilitiesIdAsArrayString, fmt.Sprintf("%d", id))
			seen[id] = true
		}
	}
	param := "{" + strings.Join(facilitiesIdAsArrayString, ",") + "}"

	frr.logger.Debugf("[%s] Will run query: %s \n\nWith params %+v \n\n\n\n", frr.path, sqlStmt, param)

	rows, err := frr.db.QueryContext(ctx, sqlStmt, param)

	returnMap := make(map[int64]entities.FacilityRequirementAggregated)
	if err != nil {
		frr.logger.Errorf("[%s] Error making query: %+v \n", frr.path, err)
		return returnMap, err
	}
	defer rows.Close()

	for rows.Next() {
		facilityRequirement := entities.FacilityRequirementAggregated{}
		err := rows.Scan(
			&facilityRequirement.FacilityId,
			pq.Array(&facilityRequirement.DocumentsId),
		)
		if err != nil {
			frr.logger.Errorf("[%s] Error scanning query: %+v \n", frr.path, err)
			return returnMap, err
		}
		returnMap[facilityRequirement.FacilityId] = facilityRequirement
	}
	return returnMap, nil
}
