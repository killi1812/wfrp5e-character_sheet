package auth

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/v2/bson"
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

	dbName := "wfrp5e_test_auth_" + uuid.New().String()[:8]
	db := client.Database(dbName)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_ = db.Drop(cleanupCtx)
	})
	return db
}

func TestAuthService(t *testing.T) {
	app.AccessKey = "test-auth-access-key-12345"
	app.RefreshKey = "test-auth-refresh-key-12345"

	db := getTestDb(t)
	logger := zap.NewNop().Sugar()
	svc := &AuthService{db: db, logger: logger}

	// Insert test user
	userUuid := uuid.New()
	hash, err := authUtil.HashPassword("validPassword123")
	assert.NoError(t, err)

	testUser := user.User{
		Uuid:         userUuid,
		Username:     "tester",
		Email:        "tester@example.com",
		PasswordHash: hash,
		Role:         user.ROLE_USER,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = db.Collection("users").InsertOne(ctx, testUser)
	assert.NoError(t, err)

	var activeSessionUuid uuid.UUID
	var activeRefreshToken string

	t.Run("Login with valid username creates new session", func(t *testing.T) {
		accessToken, refreshToken, sessionUuid, err := svc.Login("tester", "validPassword123", nil, "Mozilla/5.0", "127.0.0.1")
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.NotEqual(t, uuid.Nil, sessionUuid)

		activeSessionUuid = sessionUuid
		activeRefreshToken = refreshToken

		// Verify session was inserted in MongoDB
		var session user.Session
		findErr := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": sessionUuid}).Decode(&session)
		assert.NoError(t, findErr)
		assert.Equal(t, userUuid, session.UserUuid)
		assert.Equal(t, authUtil.HashToken(refreshToken), session.RefreshTokenHash)
		assert.Equal(t, "Mozilla/5.0", session.UserAgent)
		assert.Equal(t, "127.0.0.1", session.ClientIP)
		assert.WithinDuration(t, time.Now().Add(30*24*time.Hour), session.ExpiresAt, 5*time.Second)
	})

	t.Run("Login with existing requestedSession reuses session UUID", func(t *testing.T) {
		accessToken, refreshToken, sessionUuid, err := svc.Login("tester", "validPassword123", &activeSessionUuid, "UpdatedAgent/1.0", "192.168.1.100")
		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		assert.Equal(t, activeSessionUuid, sessionUuid)

		activeRefreshToken = refreshToken

		// Verify session was updated in MongoDB
		var session user.Session
		findErr := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": activeSessionUuid}).Decode(&session)
		assert.NoError(t, findErr)
		assert.Equal(t, "UpdatedAgent/1.0", session.UserAgent)
		assert.Equal(t, "192.168.1.100", session.ClientIP)
		assert.Equal(t, authUtil.HashToken(refreshToken), session.RefreshTokenHash)
	})

	t.Run("Login with non-matching requestedSession creates new session UUID", func(t *testing.T) {
		randomSessionID := uuid.New()
		_, _, sessionUuid, err := svc.Login("tester", "validPassword123", &randomSessionID, "Agent/2.0", "10.0.0.1")
		assert.NoError(t, err)
		assert.NotEqual(t, randomSessionID, sessionUuid)
	})

	t.Run("Login with wrong password returns ErrInvalidCredentials", func(t *testing.T) {
		_, _, _, err := svc.Login("tester", "wrongPassword", nil, "", "")
		assert.ErrorIs(t, err, cerror.ErrInvalidCredentials)
	})

	t.Run("Login with unknown username returns ErrInvalidCredentials", func(t *testing.T) {
		_, _, _, err := svc.Login("nonexistent", "validPassword123", nil, "", "")
		assert.ErrorIs(t, err, cerror.ErrInvalidCredentials)
	})

	t.Run("Login with user having nil uuid triggers GenerateTokens error", func(t *testing.T) {
		nilUser := user.User{
			Uuid:         uuid.Nil,
			Username:     "niluser",
			Email:        "nil@example.com",
			PasswordHash: hash,
			Role:         user.ROLE_USER,
		}
		_, _ = db.Collection("users").InsertOne(ctx, nilUser)
		_, _, _, err := svc.Login("niluser", "validPassword123", nil, "", "")
		assert.Error(t, err)
	})

	t.Run("RefreshTokens rotates token successfully", func(t *testing.T) {
		newAccessToken, newRefreshToken, sessionUuid, err := svc.RefreshTokens(activeRefreshToken, "RotatedAgent/1.0", "127.0.0.1")
		assert.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, activeRefreshToken, newRefreshToken)
		assert.Equal(t, activeSessionUuid, sessionUuid)

		// Verify session hash updated in MongoDB
		var session user.Session
		findErr := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": activeSessionUuid}).Decode(&session)
		assert.NoError(t, findErr)
		assert.Equal(t, authUtil.HashToken(newRefreshToken), session.RefreshTokenHash)

		// Replay attack test: calling RefreshTokens again with the OLD token MUST fail and terminate session
		_, _, _, replayErr := svc.RefreshTokens(activeRefreshToken, "AttackerAgent/1.0", "6.6.6.6")
		assert.Error(t, replayErr)
		assert.Contains(t, replayErr.Error(), "replay attack detected")

		// Session should have been deleted by replay protection
		findErrAfterReplay := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": activeSessionUuid}).Decode(&session)
		assert.ErrorIs(t, findErrAfterReplay, mongo.ErrNoDocuments)
	})

	t.Run("RefreshTokens with malformed token", func(t *testing.T) {
		_, _, _, err := svc.RefreshTokens("not-a-valid-jwt", "", "")
		assert.ErrorIs(t, err, cerror.ErrInvalidTokenFormat)
	})

	t.Run("RefreshTokens with missing sessionUuid or userUuid", func(t *testing.T) {
		claims := &authUtil.Claims{
			UserUuid: uuid.Nil,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, _ := token.SignedString([]byte(app.RefreshKey))

		_, _, _, err := svc.RefreshTokens(signed, "", "")
		assert.ErrorIs(t, err, cerror.ErrInvalidTokenFormat)
	})

	t.Run("RefreshTokens with session not in MongoDB", func(t *testing.T) {
		orphanSessionUuid := uuid.New()
		_, refreshToken, err := authUtil.GenerateTokens("orphan@test.com", "orphan", "user", userUuid, orphanSessionUuid)
		assert.NoError(t, err)

		_, _, _, err = svc.RefreshTokens(refreshToken, "", "")
		assert.ErrorIs(t, err, cerror.ErrInvalidTokenFormat)
	})

	t.Run("RefreshTokens with user not in users collection", func(t *testing.T) {
		ghostUserUuid := uuid.New()
		ghostSessionUuid := uuid.New()
		_, refreshToken, err := authUtil.GenerateTokens("ghost@test.com", "ghost", "user", ghostUserUuid, ghostSessionUuid)
		assert.NoError(t, err)

		// Insert session for ghost user
		_, err = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      ghostSessionUuid,
			UserUuid:         ghostUserUuid,
			RefreshTokenHash: authUtil.HashToken(refreshToken),
			ExpiresAt:        time.Now().Add(time.Hour),
		})
		assert.NoError(t, err)

		_, _, _, err = svc.RefreshTokens(refreshToken, "", "")
		assert.Error(t, err)
	})

	t.Run("ListSessions returns sessions with IsCurrent", func(t *testing.T) {
		session1Uuid := uuid.New()
		session2Uuid := uuid.New()

		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      session1Uuid,
			UserUuid:         userUuid,
			RefreshTokenHash: "hash1",
			UserAgent:        "Device1",
			ClientIP:         "1.1.1.1",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			ExpiresAt:        time.Now().Add(time.Hour),
		})
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      session2Uuid,
			UserUuid:         userUuid,
			RefreshTokenHash: "hash2",
			UserAgent:        "Device2",
			ClientIP:         "2.2.2.2",
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			ExpiresAt:        time.Now().Add(time.Hour),
		})

		sessions, err := svc.ListSessions(userUuid, session1Uuid)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(sessions), 2)

		var found1, found2 bool
		for _, s := range sessions {
			if s.SessionUuid == session1Uuid {
				found1 = true
				assert.True(t, s.IsCurrent)
			}
			if s.SessionUuid == session2Uuid {
				found2 = true
				assert.False(t, s.IsCurrent)
			}
		}
		assert.True(t, found1)
		assert.True(t, found2)
	})

	t.Run("RevokeSession removes specified session", func(t *testing.T) {
		targetSession := uuid.New()
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      targetSession,
			UserUuid:         userUuid,
			RefreshTokenHash: "targetHash",
			ExpiresAt:        time.Now().Add(time.Hour),
		})

		err := svc.RevokeSession(userUuid, targetSession)
		assert.NoError(t, err)

		var session user.Session
		findErr := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": targetSession}).Decode(&session)
		assert.ErrorIs(t, findErr, mongo.ErrNoDocuments)
	})

	t.Run("ValidateSession returns true for active session and false for non-existent", func(t *testing.T) {
		validSession := uuid.New()
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      validSession,
			UserUuid:         userUuid,
			RefreshTokenHash: "valHash",
			ExpiresAt:        time.Now().Add(time.Hour),
		})

		assert.True(t, svc.ValidateSession(validSession))
		assert.False(t, svc.ValidateSession(uuid.New()))
	})

	t.Run("LogoutCurrent removes specific session", func(t *testing.T) {
		currentSession := uuid.New()
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      currentSession,
			UserUuid:         userUuid,
			RefreshTokenHash: "currHash",
			ExpiresAt:        time.Now().Add(time.Hour),
		})

		err := svc.LogoutCurrent(userUuid, currentSession)
		assert.NoError(t, err)

		var session user.Session
		findErr := db.Collection("sessions").FindOne(ctx, bson.M{"session_uuid": currentSession}).Decode(&session)
		assert.ErrorIs(t, findErr, mongo.ErrNoDocuments)
	})

	t.Run("LogoutAll removes all sessions for user", func(t *testing.T) {
		s1 := uuid.New()
		s2 := uuid.New()
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      s1,
			UserUuid:         userUuid,
			RefreshTokenHash: "s1Hash",
			ExpiresAt:        time.Now().Add(time.Hour),
		})
		_, _ = db.Collection("sessions").InsertOne(ctx, user.Session{
			SessionUuid:      s2,
			UserUuid:         userUuid,
			RefreshTokenHash: "s2Hash",
			ExpiresAt:        time.Now().Add(time.Hour),
		})

		err := svc.LogoutAll(userUuid)
		assert.NoError(t, err)

		count, err := db.Collection("sessions").CountDocuments(ctx, bson.M{"user_uuid": userUuid})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

func TestNewAuthService(t *testing.T) {
	db := getTestDb(t)
	app.Test()
	app.Provide(func() *mongo.Database { return db })
	app.Provide(zap.S)

	svc := NewAuthService()
	assert.NotNil(t, svc)
}
