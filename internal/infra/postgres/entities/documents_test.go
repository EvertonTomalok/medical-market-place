package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentEntity(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		document := Document{}
		assert.Equal(t, document.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		document := Document{
			ID:       int64(1),
			Name:     "document",
			IsActive: true,
		}
		assert.Equal(t, document.ID, int64(1))
		assert.Equal(t, document.Name, "document")
		assert.Equal(t, document.IsActive, true)
	})
}
