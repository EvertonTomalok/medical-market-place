package helpers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	SetResponseMessageError(c, "Error message")
	assert.Equal(t, c.Writer.Status(), 400)
}
