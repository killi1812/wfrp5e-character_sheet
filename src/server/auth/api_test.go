package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockAuthService struct {
	loginFunc         func(identifier, password string) (string, error)
	refreshTokensFunc func(accessToken string) (string, error)
	logoutFunc        func(userUuid string) error
}

func (m *mockAuthService) Login(identifier, password string) (string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(identifier, password)
	}
	return "mock-access-token", nil
}

func (m *mockAuthService) RefreshTokens(accessToken string) (string, error) {
	if m.refreshTokensFunc != nil {
		return m.refreshTokensFunc(accessToken)
	}
	return "mock-refreshed-token", nil
}

func (m *mockAuthService) Logout(userUuid string) error {
	if m.logoutFunc != nil {
		return m.logoutFunc(userUuid)
	}
	return nil
}

func setupAuthTestRouter(service IAuthService) (*gin.Engine, *AuthCtn) {
	gin.SetMode(gin.TestMode)
	app.AccessKey = "test-auth-secret-access-key-12345"
	app.RefreshKey = "test-auth-secret-refresh-key-12345"

	logger := zap.NewNop().Sugar()
	ctn := &AuthCtn{
		auth:   service,
		logger: logger,
	}

	r := gin.New()
	api := r.Group("/api")
	ctn.RegisterEndpoints(api)
	return r, ctn
}

func generateValidToken(t *testing.T, userUuid uuid.UUID) string {
	token, _, err := authUtil.GenerateTokens("auth@test.com", "authuser", "user", userUuid)
	assert.NoError(t, err)
	return "Bearer " + token
}

func TestAuthCtn_Login(t *testing.T) {
	t.Run("Valid login returns 200 with TokenDto", func(t *testing.T) {
		mock := &mockAuthService{
			loginFunc: func(id, pass string) (string, error) {
				assert.Equal(t, "user1", id)
				assert.Equal(t, "pass1", pass)
				return "jwt-token-abc", nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		payload := LoginDto{Username: "user1", Password: "pass1"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var tokenResp TokenDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokenResp))
		assert.Equal(t, "jwt-token-abc", tokenResp.AccessToken)
	})

	t.Run("Invalid credentials returns 401", func(t *testing.T) {
		mock := &mockAuthService{
			loginFunc: func(id, pass string) (string, error) {
				return "", errors.New("invalid credentials")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		payload := LoginDto{Username: "wrong", Password: "wrong"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Malformed JSON returns without crashing", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("bad json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})
}

func TestAuthCtn_RefreshToken(t *testing.T) {
	userUuid := uuid.New()

	t.Run("Valid refresh returns 200 with new token", func(t *testing.T) {
		mock := &mockAuthService{
			refreshTokensFunc: func(token string) (string, error) {
				return "new-refreshed-token", nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var tokenResp TokenDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokenResp))
		assert.Equal(t, "new-refreshed-token", tokenResp.AccessToken)
	})

	t.Run("Refresh service error returns 500", func(t *testing.T) {
		mock := &mockAuthService{
			refreshTokensFunc: func(token string) (string, error) {
				return "", errors.New("refresh failed")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAuthCtn_Logout(t *testing.T) {
	userUuid := uuid.New()

	t.Run("Valid logout returns 200", func(t *testing.T) {
		mock := &mockAuthService{
			logoutFunc: func(id string) error {
				assert.Equal(t, userUuid.String(), id)
				return nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Logout with invalid token returns 401", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Authorization", "invalid-token-without-bearer")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Direct logout handler with invalid token returns 500", func(t *testing.T) {
		_, ctn := setupAuthTestRouter(&mockAuthService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodPost, "/logout", nil)
		c.Request.Header.Set("Authorization", "invalid")
		ctn.logout(c)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Logout service error returns 500", func(t *testing.T) {
		mock := &mockAuthService{
			logoutFunc: func(id string) error {
				return errors.New("logout error")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNewAuthCtn(t *testing.T) {
	app.Test()
	app.Provide(func() IAuthService {
		return &mockAuthService{}
	})
	app.Provide(zap.S)

	ctn := NewAuthCtn()
	assert.NotNil(t, ctn)
}
