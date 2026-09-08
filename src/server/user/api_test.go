package user

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
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockUserCrudService struct {
	createFunc            func(user *User, password string) (*User, error)
	readFunc              func(uuid uuid.UUID) (*User, error)
	readAllFunc           func() ([]User, error)
	updateFunc            func(uuid uuid.UUID, user *User) (*User, error)
	deleteFunc            func(uuid uuid.UUID) error
	getAllUsersFunc       func() ([]User, error)
	searchUsersByNameFunc func(query string) ([]User, error)
	getUserByEmailFunc    func(email string) (*User, error)
}

func (m *mockUserCrudService) Create(user *User, password string) (*User, error) {
	if m.createFunc != nil {
		return m.createFunc(user, password)
	}
	return user, nil
}

func (m *mockUserCrudService) Read(id uuid.UUID) (*User, error) {
	if m.readFunc != nil {
		return m.readFunc(id)
	}
	return nil, ErrRecordNotFound
}

func (m *mockUserCrudService) ReadAll() ([]User, error) {
	if m.readAllFunc != nil {
		return m.readAllFunc()
	}
	return []User{}, nil
}

func (m *mockUserCrudService) Update(id uuid.UUID, user *User) (*User, error) {
	if m.updateFunc != nil {
		return m.updateFunc(id, user)
	}
	return user, nil
}

func (m *mockUserCrudService) Delete(id uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(id)
	}
	return nil
}

func (m *mockUserCrudService) GetAllUsers() ([]User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc()
	}
	return []User{}, nil
}

func (m *mockUserCrudService) SearchUsersByName(query string) ([]User, error) {
	if m.searchUsersByNameFunc != nil {
		return m.searchUsersByNameFunc(query)
	}
	return []User{}, nil
}

func (m *mockUserCrudService) GetUserByEmail(email string) (*User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(email)
	}
	return nil, ErrRecordNotFound
}

func setupUserTestRouter(service IUserCrudService) (*gin.Engine, *UserCtn) {
	gin.SetMode(gin.TestMode)
	app.AccessKey = "test-user-secret-access-key-12345"
	app.RefreshKey = "test-user-secret-refresh-key-12345"

	logger := zap.NewNop().Sugar()
	ctn := &UserCtn{
		UserCrud: service,
		logger:   logger,
	}

	r := gin.New()
	api := r.Group("/api")
	ctn.RegisterEndpoints(api)
	return r, ctn
}

func generateUserAuthHeader(t *testing.T, userUuid uuid.UUID, role string) string {
	token, _, err := auth.GenerateTokens("user@example.com", "testuser", role, userUuid)
	assert.NoError(t, err)
	return "Bearer " + token
}

func TestUserCtn_Get(t *testing.T) {
	targetUuid := uuid.New()
	adminUuid := uuid.New()

	t.Run("Admin gets existing user (200)", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return &User{Uuid: id, Username: "target", Email: "target@test.com", Role: ROLE_USER}, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var uDto UserDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &uDto))
		assert.Equal(t, "target", uDto.Username)
	})

	t.Run("User not found returns 404", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return nil, ErrRecordNotFound
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return nil, errors.New("db error")
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		req, _ := http.NewRequest(http.MethodGet, "/api/user/bad-uuid", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserCtn_Create(t *testing.T) {
	adminUuid := uuid.New()

	t.Run("Valid creation returns 201", func(t *testing.T) {
		mock := &mockUserCrudService{
			createFunc: func(u *User, pass string) (*User, error) {
				assert.Equal(t, "newuser", u.Username)
				assert.Equal(t, "secretpass", pass)
				u.Uuid = uuid.New()
				return u, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		payload := NewUserDto{
			Username: "newuser",
			Email:    "new@test.com",
			Password: "secretpass",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid role returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		payload := map[string]any{
			"username": "baduser",
			"email":    "bad@test.com",
			"password": "secretpass",
			"role":     "superhacker",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON returns without crash", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})
		req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader([]byte("not-json")))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			createFunc: func(u *User, pass string) (*User, error) {
				return nil, errors.New("duplicate key")
			},
		}
		router, _ := setupUserTestRouter(mock)

		payload := NewUserDto{
			Username: "user",
			Email:    "user@test.com",
			Password: "password",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/user/", bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserCtn_Update(t *testing.T) {
	targetUuid := uuid.New()
	adminUuid := uuid.New()

	t.Run("Valid update returns 200", func(t *testing.T) {
		mock := &mockUserCrudService{
			updateFunc: func(id uuid.UUID, u *User) (*User, error) {
				u.Uuid = id
				return u, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		payload := UserDto{
			Uuid:     targetUuid.String(),
			Username: "updatedname",
			Email:    "updated@test.com",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/user/"+targetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Invalid UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		req, _ := http.NewRequest(http.MethodPut, "/api/user/not-a-uuid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid body role returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		payload := map[string]any{
			"uuid":     targetUuid.String(),
			"username": "updatedname",
			"role":     "invalid_role",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/user/"+targetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Invalid JSON body returns without crash", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})
		req, _ := http.NewRequest(http.MethodPut, "/api/user/"+targetUuid.String(), bytes.NewReader([]byte("not-json")))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			updateFunc: func(id uuid.UUID, u *User) (*User, error) {
				return nil, errors.New("update error")
			},
		}
		router, _ := setupUserTestRouter(mock)

		payload := UserDto{
			Uuid:     targetUuid.String(),
			Username: "test",
			Role:     "user",
		}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/user/"+targetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserCtn_Delete(t *testing.T) {
	targetUuid := uuid.New()
	adminUuid := uuid.New()

	t.Run("Valid delete returns 204", func(t *testing.T) {
		mock := &mockUserCrudService{
			deleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Not found returns 404", func(t *testing.T) {
		mock := &mockUserCrudService{
			deleteFunc: func(id uuid.UUID) error {
				return ErrRecordNotFound
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			deleteFunc: func(id uuid.UUID) error {
				return errors.New("delete error")
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/user/"+targetUuid.String(), nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid UUID returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		req, _ := http.NewRequest(http.MethodDelete, "/api/user/invalid", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserCtn_GetLoggedInUser(t *testing.T) {
	userUuid := uuid.New()

	t.Run("Authenticated user gets their data (200)", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return &User{Uuid: id, Username: "loggedin", Role: ROLE_USER}, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/my-data", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var uDto UserDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &uDto))
		assert.Equal(t, "loggedin", uDto.Username)
	})

	t.Run("User not found returns 404", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return nil, ErrRecordNotFound
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/my-data", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Read failure returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			readFunc: func(id uuid.UUID) (*User, error) {
				return nil, errors.New("read failed")
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/my-data", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid token returns 401", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})
		req, _ := http.NewRequest(http.MethodGet, "/api/user/my-data", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestUserCtn_GetAllUsersForSuperAdmin(t *testing.T) {
	adminUuid := uuid.New()
	userUuid := uuid.New()

	t.Run("Invalid token returns 401", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})
		req, _ := http.NewRequest(http.MethodGet, "/api/user/all-users", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Admin retrieves all users (200)", func(t *testing.T) {
		mock := &mockUserCrudService{
			getAllUsersFunc: func() ([]User, error) {
				return []User{
					{Username: "u1", Role: ROLE_USER},
					{Username: "u2", Role: ROLE_ADMIN},
				}, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/all-users", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var dtos []UserDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &dtos))
		assert.Len(t, dtos, 2)
	})

	t.Run("Regular user gets 403 Forbidden", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		req, _ := http.NewRequest(http.MethodGet, "/api/user/all-users", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			getAllUsersFunc: func() ([]User, error) {
				return nil, errors.New("db failure")
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/all-users", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestUserCtn_SearchUsersByName(t *testing.T) {
	adminUuid := uuid.New()

	t.Run("Valid query returns matched users (200)", func(t *testing.T) {
		mock := &mockUserCrudService{
			searchUsersByNameFunc: func(q string) ([]User, error) {
				assert.Equal(t, "felix", q)
				return []User{{Username: "felix", Role: ROLE_USER}}, nil
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/search?query=felix", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var dtos []UserDto
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &dtos))
		assert.Len(t, dtos, 1)
	})

	t.Run("Empty query returns 400", func(t *testing.T) {
		router, _ := setupUserTestRouter(&mockUserCrudService{})

		req, _ := http.NewRequest(http.MethodGet, "/api/user/search?query=", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockUserCrudService{
			searchUsersByNameFunc: func(q string) ([]User, error) {
				return nil, errors.New("search error")
			},
		}
		router, _ := setupUserTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/user/search?query=fail", nil)
		req.Header.Set("Authorization", generateUserAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestNewUserCtn(t *testing.T) {
	app.Test()
	app.Provide(func() IUserCrudService {
		return &mockUserCrudService{}
	})
	app.Provide(zap.S)

	ctn := NewUserCtn()
	assert.NotNil(t, ctn)
}
