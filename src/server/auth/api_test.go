package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockAuthService struct {
	loginFunc         func(username, password string, requestedSession *uuid.UUID, userAgent, clientIP string) (string, string, uuid.UUID, error)
	refreshTokensFunc func(rawRefreshToken string, userAgent, clientIP string) (string, string, uuid.UUID, error)
	logoutCurrentFunc func(userUuid, sessionUuid uuid.UUID) error
	logoutAllFunc     func(userUuid uuid.UUID) error
	listSessionsFunc  func(userUuid, currentSessionUuid uuid.UUID) ([]SessionDto, error)
	revokeSessionFunc func(userUuid, sessionUuid uuid.UUID) error
	validateSessionFunc func(sessionUuid uuid.UUID) bool
}

func (m *mockAuthService) Login(username, password string, requestedSession *uuid.UUID, userAgent, clientIP string) (string, string, uuid.UUID, error) {
	if m.loginFunc != nil {
		return m.loginFunc(username, password, requestedSession, userAgent, clientIP)
	}
	return "mock-access-token", "mock-refresh-token", uuid.New(), nil
}

func (m *mockAuthService) RefreshTokens(rawRefreshToken string, userAgent, clientIP string) (string, string, uuid.UUID, error) {
	if m.refreshTokensFunc != nil {
		return m.refreshTokensFunc(rawRefreshToken, userAgent, clientIP)
	}
	return "mock-new-access-token", "mock-new-refresh-token", uuid.New(), nil
}

func (m *mockAuthService) LogoutCurrent(userUuid, sessionUuid uuid.UUID) error {
	if m.logoutCurrentFunc != nil {
		return m.logoutCurrentFunc(userUuid, sessionUuid)
	}
	return nil
}

func (m *mockAuthService) LogoutAll(userUuid uuid.UUID) error {
	if m.logoutAllFunc != nil {
		return m.logoutAllFunc(userUuid)
	}
	return nil
}

func (m *mockAuthService) ListSessions(userUuid, currentSessionUuid uuid.UUID) ([]SessionDto, error) {
	if m.listSessionsFunc != nil {
		return m.listSessionsFunc(userUuid, currentSessionUuid)
	}
	return []SessionDto{
		{
			SessionUuid: currentSessionUuid,
			UserAgent:   "TestBrowser/1.0",
			ClientIP:    "127.0.0.1",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now(),
			IsCurrent:   true,
		},
	}, nil
}

func (m *mockAuthService) RevokeSession(userUuid, sessionUuid uuid.UUID) error {
	if m.revokeSessionFunc != nil {
		return m.revokeSessionFunc(userUuid, sessionUuid)
	}
	return nil
}

func (m *mockAuthService) ValidateSession(sessionUuid uuid.UUID) bool {
	if m.validateSessionFunc != nil {
		return m.validateSessionFunc(sessionUuid)
	}
	return true
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

func generateValidToken(t *testing.T, userUuid, sessionUuid uuid.UUID) string {
	token, _, err := authUtil.GenerateTokens("auth@test.com", "authuser", "user", userUuid, sessionUuid)
	assert.NoError(t, err)
	return "Bearer " + token
}

func TestAuthCtn_Login(t *testing.T) {
	t.Run("Valid login returns 200 with TokenDto and sets cookie", func(t *testing.T) {
		sessionID := uuid.New()
		mock := &mockAuthService{
			loginFunc: func(username, pass string, requestedSession *uuid.UUID, userAgent, clientIP string) (string, string, uuid.UUID, error) {
				assert.Equal(t, "user1", username)
				assert.Equal(t, "pass1", pass)
				assert.NotNil(t, requestedSession)
				assert.Equal(t, sessionID, *requestedSession)
				return "jwt-access-abc", "jwt-refresh-xyz", sessionID, nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		payload := LoginDto{Username: "user1", Password: "pass1", SessionId: &sessionID}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var tokenResp TokenDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokenResp))
		assert.Equal(t, "jwt-access-abc", tokenResp.AccessToken)
		assert.Equal(t, sessionID, tokenResp.SessionId)

		cookies := w.Result().Cookies()
		var refreshCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "refresh_token" {
				refreshCookie = c
				break
			}
		}
		assert.NotNil(t, refreshCookie)
		assert.Equal(t, "jwt-refresh-xyz", refreshCookie.Value)
		assert.Equal(t, "/api/auth", refreshCookie.Path)
		assert.True(t, refreshCookie.HttpOnly)
	})

	t.Run("Invalid credentials returns 401 with standard error format", func(t *testing.T) {
		mock := &mockAuthService{
			loginFunc: func(username, pass string, requestedSession *uuid.UUID, userAgent, clientIP string) (string, string, uuid.UUID, error) {
				return "", "", uuid.Nil, errors.New("invalid credentials")
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
		var errMap map[string]string
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &errMap))
		assert.Equal(t, "Invalid login credentials", errMap["error"])
	})

	t.Run("Malformed JSON returns 400", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("bad json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAuthCtn_RefreshToken(t *testing.T) {
	t.Run("Valid refresh via cookie returns 200 with new token and rotates cookie", func(t *testing.T) {
		sessionID := uuid.New()
		mock := &mockAuthService{
			refreshTokensFunc: func(rawRefreshToken, userAgent, clientIP string) (string, string, uuid.UUID, error) {
				assert.Equal(t, "old-refresh-cookie", rawRefreshToken)
				return "new-access-token", "new-refresh-cookie", sessionID, nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "old-refresh-cookie",
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var tokenResp TokenDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &tokenResp))
		assert.Equal(t, "new-access-token", tokenResp.AccessToken)
		assert.Equal(t, sessionID, tokenResp.SessionId)

		// Check rotated cookie
		var rotatedCookie *http.Cookie
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				rotatedCookie = c
				break
			}
		}
		assert.NotNil(t, rotatedCookie)
		assert.Equal(t, "new-refresh-cookie", rotatedCookie.Value)
	})

	t.Run("Valid refresh via Authorization header fallback returns 200", func(t *testing.T) {
		sessionID := uuid.New()
		mock := &mockAuthService{
			refreshTokensFunc: func(rawRefreshToken, userAgent, clientIP string) (string, string, uuid.UUID, error) {
				assert.Equal(t, "header-token", rawRefreshToken)
				return "new-access-token", "new-refresh-header", sessionID, nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.Header.Set("Authorization", "header-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Missing refresh token returns 401", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Refresh service error clears cookie and returns 401", func(t *testing.T) {
		mock := &mockAuthService{
			refreshTokensFunc: func(rawRefreshToken, userAgent, clientIP string) (string, string, uuid.UUID, error) {
				return "", "", uuid.Nil, errors.New("replay attack detected")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/refresh", nil)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "compromised-token",
		})
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		// Check that cookie was cleared (MaxAge < 0)
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				assert.True(t, c.MaxAge < 0)
			}
		}
	})
}

func TestAuthCtn_Logout(t *testing.T) {
	userUuid := uuid.New()
	sessionUuid := uuid.New()

	t.Run("Valid logout returns 204 No Content and clears cookie", func(t *testing.T) {
		loggedOutCalled := false
		mock := &mockAuthService{
			logoutCurrentFunc: func(u, s uuid.UUID) error {
				assert.Equal(t, userUuid, u)
				assert.Equal(t, sessionUuid, s)
				loggedOutCalled = true
				return nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, sessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, loggedOutCalled)
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				assert.True(t, c.MaxAge < 0)
			}
		}
	})

	t.Run("Logout without auth token returns 401", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthCtn_LogoutAll(t *testing.T) {
	userUuid := uuid.New()
	sessionUuid := uuid.New()

	t.Run("Valid logout-all returns 204 No Content and calls LogoutAll", func(t *testing.T) {
		called := false
		mock := &mockAuthService{
			logoutAllFunc: func(u uuid.UUID) error {
				assert.Equal(t, userUuid, u)
				called = true
				return nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPost, "/api/auth/logout-all", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, sessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, called)
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				assert.True(t, c.MaxAge < 0)
			}
		}
	})
}

func TestAuthCtn_ListSessions(t *testing.T) {
	userUuid := uuid.New()
	sessionUuid := uuid.New()

	t.Run("Returns 200 with list of sessions", func(t *testing.T) {
		mock := &mockAuthService{
			listSessionsFunc: func(u, s uuid.UUID) ([]SessionDto, error) {
				assert.Equal(t, userUuid, u)
				assert.Equal(t, sessionUuid, s)
				return []SessionDto{
					{
						SessionUuid: sessionUuid,
						UserAgent:   "Chrome/120.0",
						ClientIP:    "192.168.1.1",
						IsCurrent:   true,
					},
					{
						SessionUuid: uuid.New(),
						UserAgent:   "Firefox/120.0",
						ClientIP:    "192.168.1.2",
						IsCurrent:   false,
					},
				}, nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/auth/sessions", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, sessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var sessions []SessionDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &sessions))
		assert.Len(t, sessions, 2)
		assert.True(t, sessions[0].IsCurrent)
		assert.False(t, sessions[1].IsCurrent)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockAuthService{
			listSessionsFunc: func(u, s uuid.UUID) ([]SessionDto, error) {
				return nil, errors.New("db error")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/auth/sessions", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, sessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Direct handler call with missing userUuid returns 401", func(t *testing.T) {
		_, ctn := setupAuthTestRouter(&mockAuthService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/auth/sessions", nil)
		ctn.listSessions(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestAuthCtn_RevokeSession(t *testing.T) {
	userUuid := uuid.New()
	currentSessionUuid := uuid.New()
	targetSessionUuid := uuid.New()

	t.Run("Revoking remote session returns 204 No Content", func(t *testing.T) {
		called := false
		mock := &mockAuthService{
			revokeSessionFunc: func(u, s uuid.UUID) error {
				assert.Equal(t, userUuid, u)
				assert.Equal(t, targetSessionUuid, s)
				called = true
				return nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/auth/sessions/"+targetSessionUuid.String(), nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, currentSessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.True(t, called)
	})

	t.Run("Revoking current session clears cookie and returns 204", func(t *testing.T) {
		mock := &mockAuthService{
			revokeSessionFunc: func(u, s uuid.UUID) error {
				assert.Equal(t, currentSessionUuid, s)
				return nil
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/auth/sessions/"+currentSessionUuid.String(), nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, currentSessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		for _, c := range w.Result().Cookies() {
			if c.Name == "refresh_token" {
				assert.True(t, c.MaxAge < 0)
			}
		}
	})

	t.Run("Invalid session UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupAuthTestRouter(&mockAuthService{})

		req, _ := http.NewRequest(http.MethodDelete, "/api/auth/sessions/not-a-uuid", nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, currentSessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockAuthService{
			revokeSessionFunc: func(u, s uuid.UUID) error {
				return errors.New("db delete error")
			},
		}
		router, _ := setupAuthTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/auth/sessions/"+targetSessionUuid.String(), nil)
		req.Header.Set("Authorization", generateValidToken(t, userUuid, currentSessionUuid))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Direct handler call with missing userUuid returns 401", func(t *testing.T) {
		_, ctn := setupAuthTestRouter(&mockAuthService{})
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodDelete, "/api/auth/sessions/"+targetSessionUuid.String(), nil)
		ctn.revokeSession(c)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
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
