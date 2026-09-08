package user

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

func getTestDb(t *testing.T) *mongo.Database {
	uri := os.Getenv("MONGO_CONN")
	if uri == "" {
		uri = "mongodb://127.0.0.1:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		t.Skipf("Skipping DB test: unable to connect to mongo: %v", err)
		return nil
	}
	if err := client.Ping(ctx, nil); err != nil {
		t.Skipf("Skipping DB test: mongo ping failed: %v", err)
		return nil
	}

	dbName := "wfrp5e_test_usr_" + uuid.New().String()[:8]
	db := client.Database(dbName)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})
	return db
}

func TestUserCrudService_CRUD(t *testing.T) {
	db := getTestDb(t)
	logger := zap.NewNop().Sugar()
	svc := &UserCrudService{db: db, logger: logger}

	var createdUser *User

	t.Run("Create user hashes password and sets fields", func(t *testing.T) {
		u := &User{
			Username: "sigmar",
			Email:    "sigmar@empire.com",
			Role:     ROLE_ADMIN,
		}
		res, err := svc.Create(u, "hammerofsigmar")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.NotEqual(t, uuid.Nil, res.Uuid)
		assert.NotEmpty(t, res.PasswordHash)
		assert.False(t, res.CreatedAt.IsZero())
		createdUser = res
	})

	t.Run("Read existing user", func(t *testing.T) {
		found, err := svc.Read(createdUser.Uuid)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "sigmar", found.Username)
		assert.Equal(t, "sigmar@empire.com", found.Email)
	})

	t.Run("Read non-existent user returns ErrRecordNotFound", func(t *testing.T) {
		_, err := svc.Read(uuid.New())
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("GetUserByEmail finds existing user", func(t *testing.T) {
		found, err := svc.GetUserByEmail("sigmar@empire.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, createdUser.Uuid, found.Uuid)
	})

	t.Run("GetUserByEmail non-existent returns ErrRecordNotFound", func(t *testing.T) {
		_, err := svc.GetUserByEmail("unknown@nowhere.com")
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Create second user and ReadAll / GetAllUsers", func(t *testing.T) {
		u2 := &User{
			Username: "karlfranz",
			Email:    "karl@empire.com",
			Role:     ROLE_USER,
		}
		_, err := svc.Create(u2, "emperorpass")
		assert.NoError(t, err)

		all1, err := svc.ReadAll()
		assert.NoError(t, err)
		assert.Len(t, all1, 2)

		all2, err := svc.GetAllUsers()
		assert.NoError(t, err)
		assert.Len(t, all2, 2)
	})

	t.Run("SearchUsersByName matches close query", func(t *testing.T) {
		matches, err := svc.SearchUsersByName("sigmar")
		assert.NoError(t, err)
		assert.NotEmpty(t, matches)
		assert.Equal(t, "sigmar", matches[0].Username)

		noMatches, err := svc.SearchUsersByName("completelyunrelatedstring12345")
		assert.NoError(t, err)
		assert.Empty(t, noMatches)
	})

	t.Run("Update modifies user", func(t *testing.T) {
		toUpdate := &User{
			Username: "sigmar_unifier",
			Email:    "unifier@empire.com",
			Role:     ROLE_ADMIN,
		}
		updated, err := svc.Update(createdUser.Uuid, toUpdate)
		assert.NoError(t, err)
		assert.Equal(t, "sigmar_unifier", updated.Username)
		assert.Equal(t, "unifier@empire.com", updated.Email)
	})

	t.Run("Update non-existent user returns ErrRecordNotFound", func(t *testing.T) {
		_, err := svc.Update(uuid.New(), &User{Username: "ghost"})
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Delete removes user", func(t *testing.T) {
		err := svc.Delete(createdUser.Uuid)
		assert.NoError(t, err)

		_, err = svc.Read(createdUser.Uuid)
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})

	t.Run("Delete non-existent user returns ErrRecordNotFound", func(t *testing.T) {
		err := svc.Delete(createdUser.Uuid)
		assert.ErrorIs(t, err, ErrRecordNotFound)
	})
}

func TestNewUserCrudService(t *testing.T) {
	db := getTestDb(t)
	app.Test()
	app.Provide(func() *mongo.Database { return db })
	app.Provide(zap.S)

	svc := NewUserCrudService()
	assert.NotNil(t, svc)
}
