package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFacilityEntity(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		facility := Facility{}
		assert.Equal(t, facility.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		facility := Facility{
			ID:       int64(1),
			Name:     "facility",
			IsActive: true,
		}
		assert.Equal(t, facility.ID, int64(1))
		assert.Equal(t, facility.Name, "facility")
		assert.Equal(t, facility.IsActive, true)
	})

	t.Run("valid entity to dto", func(t *testing.T) {
		facility := Facility{
			ID:       int64(1),
			Name:     "facility",
			IsActive: true,
		}
		dto := facility.ToDTO()
		assert.Equal(t, dto.ID, int64(1))
		assert.Equal(t, dto.Name, "facility")
		assert.Equal(t, dto.IsActive, true)
	})
}
