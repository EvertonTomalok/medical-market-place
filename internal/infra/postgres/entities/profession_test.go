package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProfessionEnum(t *testing.T) {
	t.Run("valid enum", func(t *testing.T) {
		profession := Profession("CNA")
		assert.Equal(t, profession.IsValid(), true)
	})

	t.Run("invlid valid enum", func(t *testing.T) {
		profession := Profession("INVALID")
		assert.Equal(t, profession.IsValid(), false)
	})
}
