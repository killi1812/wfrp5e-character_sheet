package ginutil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseUUIDParam extracts and parses a UUID route parameter, aborting with 400 Bad Request on failure.
func ParseUUIDParam(c *gin.Context, param string) (uuid.UUID, bool) {
	parsedUUID, err := uuid.Parse(c.Param(param))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return uuid.Nil, false
	}
	return parsedUUID, true
}

// HandleServiceError inspects a service error.
// If err is nil, returns false.
// If err matches notFoundErr, aborts with 404 Not Found and returns true.
// Otherwise, aborts with 500 Internal Server Error and returns true.
func HandleServiceError(c *gin.Context, err error, notFoundErr error) bool {
	if err == nil {
		return false
	}
	if notFoundErr != nil && errors.Is(err, notFoundErr) {
		c.AbortWithError(http.StatusNotFound, err)
		return true
	}
	c.AbortWithError(http.StatusInternalServerError, err)
	return true
}
