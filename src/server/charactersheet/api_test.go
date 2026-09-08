package charactersheet

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

// Mock ICharacterSheetService
type mockSheetService struct {
	createFunc     func(sheet *CharacterSheet) (*CharacterSheet, error)
	readFunc       func(uuid uuid.UUID) (*CharacterSheet, error)
	readAllFunc    func() ([]CharacterSheet, error)
	readByUserFunc func(userUuid uuid.UUID) ([]CharacterSheet, error)
	updateFunc     func(uuid uuid.UUID, sheet *CharacterSheet) (*CharacterSheet, error)
	deleteFunc     func(uuid uuid.UUID) error
}

func (m *mockSheetService) Create(sheet *CharacterSheet) (*CharacterSheet, error) {
	if m.createFunc != nil {
		return m.createFunc(sheet)
	}
	return sheet, nil
}

func (m *mockSheetService) Read(sheetUuid uuid.UUID) (*CharacterSheet, error) {
	if m.readFunc != nil {
		return m.readFunc(sheetUuid)
	}
	return nil, ErrSheetNotFound
}

func (m *mockSheetService) ReadAll() ([]CharacterSheet, error) {
	if m.readAllFunc != nil {
		return m.readAllFunc()
	}
	return []CharacterSheet{}, nil
}

func (m *mockSheetService) ReadByUser(userUuid uuid.UUID) ([]CharacterSheet, error) {
	if m.readByUserFunc != nil {
		return m.readByUserFunc(userUuid)
	}
	return []CharacterSheet{}, nil
}

func (m *mockSheetService) Update(sheetUuid uuid.UUID, sheet *CharacterSheet) (*CharacterSheet, error) {
	if m.updateFunc != nil {
		return m.updateFunc(sheetUuid, sheet)
	}
	return sheet, nil
}

func (m *mockSheetService) Delete(sheetUuid uuid.UUID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(sheetUuid)
	}
	return nil
}

func setupTestRouter(service ICharacterSheetService) (*gin.Engine, *CharacterSheetCtn) {
	gin.SetMode(gin.TestMode)
	app.AccessKey = "test-secret-access-key-12345"
	app.RefreshKey = "test-secret-refresh-key-12345"

	logger := zap.NewNop().Sugar()
	ctn := &CharacterSheetCtn{
		service: service,
		logger:  logger,
	}

	r := gin.New()
	api := r.Group("/api")
	ctn.RegisterEndpoints(api)
	return r, ctn
}

func generateAuthHeader(t *testing.T, userUuid uuid.UUID, role string) string {
	token, _, err := auth.GenerateTokens("test@example.com", "testuser", role, userUuid)
	assert.NoError(t, err)
	return "Bearer " + token
}

func TestCharacterSheetCtn_GetAll(t *testing.T) {
	adminUuid := uuid.New()
	userUuid := uuid.New()

	t.Run("Admin role reads all sheets", func(t *testing.T) {
		mock := &mockSheetService{
			readAllFunc: func() ([]CharacterSheet, error) {
				return []CharacterSheet{{Name: "Sheet1"}, {Name: "Sheet2"}}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []CharacterSheet
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp, 2)
	})

	t.Run("Regular user reads own sheets", func(t *testing.T) {
		mock := &mockSheetService{
			readByUserFunc: func(u uuid.UUID) ([]CharacterSheet, error) {
				assert.Equal(t, userUuid, u)
				return []CharacterSheet{{Name: "UserSheet"}}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []CharacterSheet
		assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Len(t, resp, 1)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readByUserFunc: func(u uuid.UUID) ([]CharacterSheet, error) {
				return nil, errors.New("db error")
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, userUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCharacterSheetCtn_GetOne(t *testing.T) {
	sheetUuid := uuid.New()
	ownerUuid := uuid.New()
	otherUserUuid := uuid.New()
	adminUuid := uuid.New()

	t.Run("Owner retrieves own sheet", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid, Name: "Hero"}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Admin retrieves any sheet", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid, Name: "Hero"}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Other user cannot retrieve sheet (403)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid, Name: "Hero"}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, otherUserUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Sheet not found returns 404", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, ErrSheetNotFound
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, errors.New("db failure")
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupTestRouter(&mockSheetService{})

		req, _ := http.NewRequest(http.MethodGet, "/api/character-sheets/invalid-uuid", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCharacterSheetCtn_Create(t *testing.T) {
	userUuid := uuid.New()

	t.Run("Valid creation sets UserUuid and returns 201", func(t *testing.T) {
		mock := &mockSheetService{
			createFunc: func(s *CharacterSheet) (*CharacterSheet, error) {
				assert.Equal(t, userUuid, s.UserUuid)
				s.Uuid = uuid.New()
				return s, nil
			},
		}
		router, _ := setupTestRouter(mock)

		payload := CharacterSheet{Name: "New Character", Species: "Human"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/character-sheets/", bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, userUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Invalid JSON returns 400", func(t *testing.T) {
		router, _ := setupTestRouter(&mockSheetService{})

		req, _ := http.NewRequest(http.MethodPost, "/api/character-sheets/", bytes.NewReader([]byte("not json")))
		req.Header.Set("Authorization", generateAuthHeader(t, userUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service error returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			createFunc: func(s *CharacterSheet) (*CharacterSheet, error) {
				return nil, errors.New("db error")
			},
		}
		router, _ := setupTestRouter(mock)

		payload := CharacterSheet{Name: "Fail Sheet"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPost, "/api/character-sheets/", bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, userUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCharacterSheetCtn_Update(t *testing.T) {
	sheetUuid := uuid.New()
	ownerUuid := uuid.New()
	otherUuid := uuid.New()

	t.Run("Owner updates sheet successfully (200)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid, Name: "Old"}, nil
			},
			updateFunc: func(id uuid.UUID, s *CharacterSheet) (*CharacterSheet, error) {
				s.Uuid = id
				return s, nil
			},
		}
		router, _ := setupTestRouter(mock)

		payload := CharacterSheet{Name: "Updated"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Other user cannot update sheet (403)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid, Name: "Old"}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		payload := CharacterSheet{Name: "Hijack"}
		body, _ := json.Marshal(payload)

		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, otherUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Update non-existent sheet returns 404", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, ErrSheetNotFound
			},
		}
		router, _ := setupTestRouter(mock)

		body, _ := json.Marshal(CharacterSheet{Name: "X"})
		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Read error returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, errors.New("read error")
			},
		}
		router, _ := setupTestRouter(mock)

		body, _ := json.Marshal(CharacterSheet{Name: "X"})
		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid update body returns 400", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader([]byte("bad")))
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Update service failure returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
			updateFunc: func(id uuid.UUID, s *CharacterSheet) (*CharacterSheet, error) {
				return nil, errors.New("update failed")
			},
		}
		router, _ := setupTestRouter(mock)

		body, _ := json.Marshal(CharacterSheet{Name: "X"})
		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/"+sheetUuid.String(), bytes.NewReader(body))
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupTestRouter(&mockSheetService{})

		req, _ := http.NewRequest(http.MethodPut, "/api/character-sheets/invalid", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCharacterSheetCtn_Delete(t *testing.T) {
	sheetUuid := uuid.New()
	ownerUuid := uuid.New()
	otherUuid := uuid.New()
	adminUuid := uuid.New()

	t.Run("Owner deletes sheet (200)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
			deleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Admin deletes sheet (200)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
			deleteFunc: func(id uuid.UUID) error {
				return nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, adminUuid, "admin"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Other user cannot delete (403)", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, otherUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Delete non-existent sheet returns 404", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, ErrSheetNotFound
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Read failure returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return nil, errors.New("read failure")
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Delete service error returns 500", func(t *testing.T) {
		mock := &mockSheetService{
			readFunc: func(id uuid.UUID) (*CharacterSheet, error) {
				return &CharacterSheet{Uuid: id, UserUuid: ownerUuid}, nil
			},
			deleteFunc: func(id uuid.UUID) error {
				return errors.New("delete failure")
			},
		}
		router, _ := setupTestRouter(mock)

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/"+sheetUuid.String(), nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Invalid UUID parameter returns 400", func(t *testing.T) {
		router, _ := setupTestRouter(&mockSheetService{})

		req, _ := http.NewRequest(http.MethodDelete, "/api/character-sheets/not-a-uuid", nil)
		req.Header.Set("Authorization", generateAuthHeader(t, ownerUuid, "user"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestNewCharacterSheetCtn(t *testing.T) {
	app.Test()
	app.Provide(func() ICharacterSheetService {
		return &mockSheetService{}
	})
	app.Provide(zap.S)

	ctn := NewCharacterSheetCtn()
	assert.NotNil(t, ctn)
}
