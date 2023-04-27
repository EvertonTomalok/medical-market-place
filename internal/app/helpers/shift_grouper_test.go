package helpers

import (
	"testing"
	"time"

	"github.com/EvertonTomalok/marketplace-health/internal/domain/dto"
	"github.com/EvertonTomalok/marketplace-health/internal/infra/postgres/entities"
	"github.com/stretchr/testify/assert"
)

func TestGroupShiftByDate(t *testing.T) {
	shifts := []dto.ShiftDTO{
		{
			ID:         1,
			Start:      time.Date(2023, 4, 22, 1, 10, 0, 0, time.UTC),
			End:        time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
			Profession: entities.CNA.ToString(),
			IsDeleted:  false,
			FacilityId: 1,
		},
		{
			ID:         2,
			Start:      time.Date(2023, 4, 22, 1, 10, 0, 0, time.UTC),
			End:        time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
			Profession: entities.CNA.ToString(),
			IsDeleted:  false,
			FacilityId: 1,
		},
		{
			ID:         3,
			Start:      time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
			End:        time.Date(2023, 4, 24, 1, 10, 0, 0, time.UTC),
			Profession: entities.CNA.ToString(),
			IsDeleted:  false,
			FacilityId: 1,
		},
	}

	grouped := GroupShiftByDate(shifts)
	expectedGroup := dto.GroupedByDateShift{
		"2023-04-22": {
			{
				ID:         1,
				Start:      time.Date(2023, 4, 22, 1, 10, 0, 0, time.UTC),
				End:        time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
				Profession: entities.CNA.ToString(),
				IsDeleted:  false,
				FacilityId: 1,
			},
			{
				ID:         2,
				Start:      time.Date(2023, 4, 22, 1, 10, 0, 0, time.UTC),
				End:        time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
				Profession: entities.CNA.ToString(),
				IsDeleted:  false,
				FacilityId: 1,
			},
		},
		"2023-04-23": {
			{
				ID:         3,
				Start:      time.Date(2023, 4, 23, 1, 10, 0, 0, time.UTC),
				End:        time.Date(2023, 4, 24, 1, 10, 0, 0, time.UTC),
				Profession: entities.CNA.ToString(),
				IsDeleted:  false,
				FacilityId: 1,
			},
		},
	}

	assert.Equal(t, expectedGroup, grouped)
}
