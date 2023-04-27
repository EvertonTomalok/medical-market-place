package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacilityRequirementEntity(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		facility_requirement := FacilityRequirement{}
		assert.Equal(t, facility_requirement.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		oneInt64 := int64(1)
		facility_requirement := FacilityRequirement{
			ID:         oneInt64,
			FacilityId: oneInt64,
			DocumentId: oneInt64,
		}
		assert.Equal(t, facility_requirement.ID, oneInt64)
		assert.Equal(t, facility_requirement.FacilityId, oneInt64)
		assert.Equal(t, facility_requirement.DocumentId, oneInt64)
	})
}
