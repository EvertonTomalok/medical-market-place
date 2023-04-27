package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerHasAllDocumentsRequired(t *testing.T) {
	t.Parallel()

	t.Run("has all documents required", func(t *testing.T) {
		docsRequired := []int64{1, 2}
		workerDocuments := []int64{1, 2, 3}
		expectedResult := true

		result := WorkerHasAllDocumentsRequired(workerDocuments, docsRequired)
		assert.Equal(t, result, expectedResult)
	})

	t.Run("missing all docs required", func(t *testing.T) {
		docsRequired := []int64{1, 2, 3, 4}
		workerDocuments := []int64{1, 2, 3}
		expectedResult := false

		result := WorkerHasAllDocumentsRequired(workerDocuments, docsRequired)
		assert.Equal(t, result, expectedResult)
	})

	t.Run("no docs required", func(t *testing.T) {
		docsRequired := []int64{}
		workerDocuments := []int64{1, 2, 3}
		expectedResult := true

		result := WorkerHasAllDocumentsRequired(workerDocuments, docsRequired)
		assert.Equal(t, result, expectedResult)
	})

	t.Run("worker doesn`t have any doc required", func(t *testing.T) {
		docsRequired := []int64{1, 2, 3}
		workerDocuments := []int64{}
		expectedResult := false

		result := WorkerHasAllDocumentsRequired(workerDocuments, docsRequired)
		assert.Equal(t, result, expectedResult)
	})

	t.Run("worker doesn`t have any doc required and no document required", func(t *testing.T) {
		docsRequired := []int64{}
		workerDocuments := []int64{}
		expectedResult := true

		result := WorkerHasAllDocumentsRequired(workerDocuments, docsRequired)
		assert.Equal(t, result, expectedResult)
	})
}

func TestRemoveDuplicate(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5}
	result := RemoveDuplicate([]int{1, 1, 1, 2, 2, 2, 3, 4, 5, 5, 5, 5, 5, 5, 5})
	assert.ElementsMatch(t, expected, result)
}
