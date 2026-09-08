package charactersheet

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

	dbName := "wfrp5e_test_cs_" + uuid.New().String()[:8]
	db := client.Database(dbName)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})
	return db
}

func TestCharacterSheetService_CRUD(t *testing.T) {
	db := getTestDb(t)
	logger := zap.NewNop().Sugar()
	svc := &CharacterSheetService{db: db, logger: logger}

	user1 := uuid.New()
	user2 := uuid.New()

	var sheet1Uuid uuid.UUID

	t.Run("Create sets uuid and timestamps", func(t *testing.T) {
		sheet := &CharacterSheet{
			UserUuid: user1,
			Name:     "Felix Jaeger",
			Species:  "Human",
			Class:    "Ranger",
		}
		created, err := svc.Create(sheet)
		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.NotEqual(t, uuid.Nil, created.Uuid)
		assert.False(t, created.CreatedAt.IsZero())
		assert.False(t, created.UpdatedAt.IsZero())
		assert.False(t, created.ID.IsZero())
		sheet1Uuid = created.Uuid
	})

	t.Run("Read existing sheet", func(t *testing.T) {
		found, err := svc.Read(sheet1Uuid)
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Felix Jaeger", found.Name)
		assert.Equal(t, user1, found.UserUuid)
	})

	t.Run("Read non-existent sheet returns ErrSheetNotFound", func(t *testing.T) {
		_, err := svc.Read(uuid.New())
		assert.ErrorIs(t, err, ErrSheetNotFound)
	})

	t.Run("Create second sheet and ReadAll", func(t *testing.T) {
		sheet2 := &CharacterSheet{
			UserUuid: user2,
			Name:     "Max Schreiber",
			Species:  "Human",
		}
		_, err := svc.Create(sheet2)
		assert.NoError(t, err)

		all, err := svc.ReadAll()
		assert.NoError(t, err)
		assert.Len(t, all, 2)
	})

	t.Run("ReadByUser filters correctly", func(t *testing.T) {
		user1Sheets, err := svc.ReadByUser(user1)
		assert.NoError(t, err)
		assert.Len(t, user1Sheets, 1)
		assert.Equal(t, "Felix Jaeger", user1Sheets[0].Name)

		user2Sheets, err := svc.ReadByUser(user2)
		assert.NoError(t, err)
		assert.Len(t, user2Sheets, 1)
		assert.Equal(t, "Max Schreiber", user2Sheets[0].Name)

		emptySheets, err := svc.ReadByUser(uuid.New())
		assert.NoError(t, err)
		assert.Empty(t, emptySheets)
	})

	t.Run("Update modifies sheet and preserves ID/CreatedAt", func(t *testing.T) {
		orig, err := svc.Read(sheet1Uuid)
		assert.NoError(t, err)

		updatePayload := &CharacterSheet{
			Name:    "Felix the Poet",
			Species: "Human",
			Class:   "Scholar",
		}
		updated, err := svc.Update(sheet1Uuid, updatePayload)
		assert.NoError(t, err)
		assert.Equal(t, "Felix the Poet", updated.Name)
		assert.Equal(t, orig.ID, updated.ID)
		assert.Equal(t, orig.CreatedAt.Unix(), updated.CreatedAt.Unix())
	})

	t.Run("Update non-existent sheet returns ErrSheetNotFound", func(t *testing.T) {
		_, err := svc.Update(uuid.New(), &CharacterSheet{Name: "None"})
		assert.ErrorIs(t, err, ErrSheetNotFound)
	})

	t.Run("Delete removes sheet", func(t *testing.T) {
		err := svc.Delete(sheet1Uuid)
		assert.NoError(t, err)

		_, err = svc.Read(sheet1Uuid)
		assert.ErrorIs(t, err, ErrSheetNotFound)
	})

	t.Run("Delete non-existent sheet returns ErrSheetNotFound", func(t *testing.T) {
		err := svc.Delete(sheet1Uuid)
		assert.ErrorIs(t, err, ErrSheetNotFound)
	})
}

func TestNewCharacterSheetCrudService(t *testing.T) {
	db := getTestDb(t)
	app.Test()
	app.Provide(func() *mongo.Database { return db })
	app.Provide(zap.S)

	svc := NewCharacterSheetCrudService()
	assert.NotNil(t, svc)
}
