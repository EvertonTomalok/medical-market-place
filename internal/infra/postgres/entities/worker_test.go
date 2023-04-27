package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerToDTO(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		worker := Worker{}
		workerDTO := worker.ToDTO()
		assert.Equal(t, worker.ID, int64(0))
		assert.Equal(t, workerDTO.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		worker := Worker{
			ID:         int64(1),
			Name:       "Some Name",
			IsActive:   true,
			Profession: CNA,
		}
		workerDTO := worker.ToDTO()
		assert.Equal(t, workerDTO.ID, int64(1))
		assert.Equal(t, workerDTO.Name, "Some Name")
		assert.Equal(t, workerDTO.IsActive, true)
		assert.Equal(t, workerDTO.Profession, "CNA")
	})
}
