package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftEntity(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		shift := Shift{}
		assert.Equal(t, shift.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		shift := Shift{
			ID:         int64(1),
			Profession: CNA,
		}
		assert.Equal(t, shift.ID, int64(1))
	})
}
