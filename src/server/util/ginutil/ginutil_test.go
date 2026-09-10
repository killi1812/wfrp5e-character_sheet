package ginutil

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("valid uuid", func(t *testing.T) {
		validID := uuid.New()
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "uuid", Value: validID.String()}}

		parsed, ok := ParseUUIDParam(c, "uuid")
		assert.True(t, ok)
		assert.Equal(t, validID, parsed)
		assert.False(t, c.IsAborted())
	})

	t.Run("invalid uuid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "uuid", Value: "invalid-uuid"}}

		parsed, ok := ParseUUIDParam(c, "uuid")
		assert.False(t, ok)
		assert.Equal(t, uuid.Nil, parsed)
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandleServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	errNotFound := errors.New("not found")

	t.Run("nil error returns false", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handled := HandleServiceError(c, nil, errNotFound)
		assert.False(t, handled)
		assert.False(t, c.IsAborted())
	})

	t.Run("not found error returns true with 404", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handled := HandleServiceError(c, errNotFound, errNotFound)
		assert.True(t, handled)
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("other error returns true with 500", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		handled := HandleServiceError(c, errors.New("db error"), errNotFound)
		assert.True(t, handled)
		assert.True(t, c.IsAborted())
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
