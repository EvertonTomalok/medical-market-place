package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentWorkerEntity(t *testing.T) {
	t.Run("empty entity", func(t *testing.T) {
		document_worker := DocumentWorker{}
		assert.Equal(t, document_worker.ID, int64(0))
	})

	t.Run("valid entity", func(t *testing.T) {
		oneInt64 := int64(1)
		document_worker := DocumentWorker{
			ID:         oneInt64,
			WorkerId:   oneInt64,
			DocumentId: oneInt64,
		}
		assert.Equal(t, document_worker.ID, oneInt64)
		assert.Equal(t, document_worker.WorkerId, oneInt64)
		assert.Equal(t, document_worker.DocumentId, oneInt64)
	})
}
