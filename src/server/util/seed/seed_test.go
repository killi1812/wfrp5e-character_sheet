package seed

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/user"
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

	dbName := "wfrp5e_test_seed_" + uuid.New().String()[:8]
	db := client.Database(dbName)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})
	return db
}

func TestSeed(t *testing.T) {
	db := getTestDb(t)
	app.Test()
	app.Provide(func() *mongo.Database { return db })
	app.Provide(zap.S)
	app.Provide(user.NewUserCrudService)

	origPass := os.Getenv(_PASSWORD_ENV)
	defer os.Setenv(_PASSWORD_ENV, origPass)

	t.Run("Empty password env returns error", func(t *testing.T) {
		os.Setenv(_PASSWORD_ENV, "")
		err := createSuperAdmin()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "env variable is empty")
	})

	t.Run("Short password returns error", func(t *testing.T) {
		os.Setenv(_PASSWORD_ENV, "short")
		err := createSuperAdmin()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 8 characters")
	})

	t.Run("Valid password creates superadmin and Insert succeeds", func(t *testing.T) {
		os.Setenv(_PASSWORD_ENV, "superpassword123")
		err := createSuperAdmin()
		assert.NoError(t, err)

		// Running again detects existing superadmin and skips
		err = createSuperAdmin()
		assert.NoError(t, err)

		// Insert() wrapper runs without panic
		assert.NotPanics(t, func() {
			Insert()
		})
	})

	t.Run("DB error on GetUserByEmail returns error", func(t *testing.T) {
		closedDb := getTestDb(t)
		_ = closedDb.Client().Disconnect(context.Background())
		app.Test()
		app.Provide(func() *mongo.Database { return closedDb })
		app.Provide(zap.S)

		err := createSuperAdmin()
		assert.Error(t, err)
	})

	t.Run("Insert handles error gracefully when createSuperAdmin fails", func(t *testing.T) {
		os.Setenv(_PASSWORD_ENV, "")
		assert.NotPanics(t, func() {
			Insert()
		})
	})
}
