package info_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/info"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInfoApi(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set known app info variables
	app.Build = "test-build"
	app.Version = "1.2.3"
	app.CommitHash = "abcdef1"
	app.BuildTimestamp = "2026-09-08T12:00:00"

	// Configure DI container for testing
	app.Test()
	app.Provide(zap.S)
	controller := info.NewInfoCnt()
	assert.NotNil(t, controller)

	router := gin.New()
	group := router.Group("/api")
	controller.RegisterEndpoints(group)

	t.Run("GET /api/info returns server info dto", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/info", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp info.ServerInfoDto
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "test-build", resp.Build)
		assert.Equal(t, "1.2.3", resp.Version)
		assert.Equal(t, "abcdef1", resp.CommitHash)
		assert.Equal(t, "2026-09-08T12:00:00", resp.BuildTimestamp)
	})
}
