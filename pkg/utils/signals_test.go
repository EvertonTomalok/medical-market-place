package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignals(t *testing.T) {
	assert.Equal(t, reflect.TypeOf(MakeDoneSignal()).String(), "chan os.Signal")
}
