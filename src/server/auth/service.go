package auth

import (
	"context"
	"errors"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
)

type IAuthService interface {
	Login(username, password string, requestedSession *uuid.UUID, userAgent, clientIP string) (accessToken, refreshToken string, sessionUuid uuid.UUID, err error)
	RefreshTokens(rawRefreshToken string, userAgent, clientIP string) (newAccessToken, newRefreshToken string, sessionUuid uuid.UUID, err error)
	LogoutCurrent(userUuid, sessionUuid uuid.UUID) error
	LogoutAll(userUuid uuid.UUID) error
	ListSessions(userUuid, currentSessionUuid uuid.UUID) ([]SessionDto, error)
	RevokeSession(userUuid, sessionUuid uuid.UUID) error
	ValidateSession(sessionUuid uuid.UUID) bool
}

type AuthService struct {
	db     *mongo.Database
	logger *zap.SugaredLogger
}

func NewAuthService() IAuthService {
	var service IAuthService

	app.Invoke(func(db *mongo.Database, logger *zap.SugaredLogger) {
		s := &AuthService{
			db:     db,
			logger: logger,
		}
		service = s
		authUtil.SetSessionValidator(s.ValidateSession)
	})

	return service
}

func (s *AuthService) usersColl() *mongo.Collection {
	return s.db.Collection("users")
}

func (s *AuthService) sessionsColl() *mongo.Collection {
	return s.db.Collection("sessions")
}

func (s *AuthService) Login(username, password string, requestedSession *uuid.UUID, userAgent, clientIP string) (string, string, uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usr user.User
	err := s.usersColl().FindOne(ctx, bson.M{"username": username}).Decode(&usr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Debugf("User not found username = %s", username)
			return "", "", uuid.Nil, cerror.ErrInvalidCredentials
		}
		s.logger.Errorf("Failed to query user: %v", err)
		return "", "", uuid.Nil, err
	}

	if !authUtil.VerifyPassword(usr.PasswordHash, password) {
		s.logger.Debugf("Invalid password for user username: %s, uuid: %s", usr.Username, usr.Uuid)
		return "", "", uuid.Nil, cerror.ErrInvalidCredentials
	}

	sessionUuid := uuid.New()
	if requestedSession != nil && *requestedSession != uuid.Nil {
		var existing user.Session
		findErr := s.sessionsColl().FindOne(ctx, bson.M{
			"session_uuid": *requestedSession,
			"user_uuid":    usr.Uuid,
		}).Decode(&existing)
		if findErr == nil {
			sessionUuid = *requestedSession
		}
	}

	accessToken, refreshToken, err := authUtil.GenerateTokens(usr.Email, usr.Username, string(usr.Role), usr.Uuid, sessionUuid)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens: %v", err)
		return "", "", uuid.Nil, err
	}

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	refreshTokenHash := authUtil.HashToken(refreshToken)

	opts := options.UpdateOne().SetUpsert(true)
	_, err = s.sessionsColl().UpdateOne(ctx, bson.M{"session_uuid": sessionUuid}, bson.M{
		"$set": bson.M{
			"user_uuid":          usr.Uuid,
			"refresh_token_hash": refreshTokenHash,
			"user_agent":         userAgent,
			"client_ip":          clientIP,
			"updated_at":         now,
			"expires_at":         expiresAt,
		},
		"$setOnInsert": bson.M{
			"session_uuid": sessionUuid,
			"created_at":   now,
		},
	}, opts)
	if err != nil {
		s.logger.Errorf("Failed to persist session: %v", err)
		return "", "", uuid.Nil, err
	}

	return accessToken, refreshToken, sessionUuid, nil
}

func (s *AuthService) RefreshTokens(rawRefreshToken string, userAgent, clientIP string) (string, string, uuid.UUID, error) {
	_, claims, err := authUtil.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		s.logger.Errorf("Failed to parse refresh token: %v", err)
		return "", "", uuid.Nil, cerror.ErrInvalidTokenFormat
	}

	sessionUuid := claims.SessionUuid
	userUuid := claims.UserUuid
	if sessionUuid == uuid.Nil || userUuid == uuid.Nil {
		return "", "", uuid.Nil, cerror.ErrInvalidTokenFormat
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session user.Session
	err = s.sessionsColl().FindOne(ctx, bson.M{"session_uuid": sessionUuid}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Warnf("Session not found during refresh: %s", sessionUuid)
			return "", "", uuid.Nil, cerror.ErrInvalidTokenFormat
		}
		s.logger.Errorf("Failed to query session during refresh: %v", err)
		return "", "", uuid.Nil, err
	}

	incomingHash := authUtil.HashToken(rawRefreshToken)
	if session.RefreshTokenHash != incomingHash {
		// Replay attack detected! Invalidate compromised session immediately
		s.logger.Warnf("Replay attack detected for session: %s. Terminating session.", sessionUuid)
		_, _ = s.sessionsColl().DeleteOne(ctx, bson.M{"session_uuid": sessionUuid})
		return "", "", uuid.Nil, errors.New("replay attack detected: token already rotated")
	}

	var usr user.User
	err = s.usersColl().FindOne(ctx, bson.M{"uuid": userUuid}).Decode(&usr)
	if err != nil {
		s.logger.Errorf("Failed to find user during refresh: %v", err)
		return "", "", uuid.Nil, err
	}

	newAccessToken, newRefreshToken, err := authUtil.GenerateTokens(usr.Email, usr.Username, string(usr.Role), usr.Uuid, sessionUuid)
	if err != nil {
		s.logger.Errorf("Failed to generate rotated tokens: %v", err)
		return "", "", uuid.Nil, err
	}

	now := time.Now()
	expiresAt := now.Add(30 * 24 * time.Hour)
	newHash := authUtil.HashToken(newRefreshToken)

	_, err = s.sessionsColl().UpdateOne(ctx, bson.M{"session_uuid": sessionUuid}, bson.M{
		"$set": bson.M{
			"refresh_token_hash": newHash,
			"user_agent":         userAgent,
			"client_ip":          clientIP,
			"updated_at":         now,
			"expires_at":         expiresAt,
		},
	})
	if err != nil {
		s.logger.Errorf("Failed to update rotated session: %v", err)
		return "", "", uuid.Nil, err
	}

	return newAccessToken, newRefreshToken, sessionUuid, nil
}

func (s *AuthService) LogoutCurrent(userUuid, sessionUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.sessionsColl().DeleteOne(ctx, bson.M{
		"user_uuid":    userUuid,
		"session_uuid": sessionUuid,
	})
	if err != nil {
		s.logger.Errorf("Failed to delete session: %v", err)
		return err
	}
	return nil
}

func (s *AuthService) LogoutAll(userUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.sessionsColl().DeleteMany(ctx, bson.M{"user_uuid": userUuid})
	if err != nil {
		s.logger.Errorf("Failed to delete all sessions for user %s: %v", userUuid, err)
		return err
	}
	return nil
}

func (s *AuthService) ListSessions(userUuid, currentSessionUuid uuid.UUID) ([]SessionDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.sessionsColl().Find(ctx, bson.M{
		"user_uuid":  userUuid,
		"expires_at": bson.M{"$gt": time.Now()},
	})
	if err != nil {
		s.logger.Errorf("Failed to list sessions: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var sessions []user.Session
	if err = cursor.All(ctx, &sessions); err != nil {
		s.logger.Errorf("Failed to decode sessions: %v", err)
		return nil, err
	}

	dtos := make([]SessionDto, 0, len(sessions))
	for _, sess := range sessions {
		dtos = append(dtos, SessionDto{
			SessionUuid: sess.SessionUuid,
			UserAgent:   sess.UserAgent,
			ClientIP:    sess.ClientIP,
			CreatedAt:   sess.CreatedAt,
			UpdatedAt:   sess.UpdatedAt,
			IsCurrent:   sess.SessionUuid == currentSessionUuid,
		})
	}

	return dtos, nil
}

func (s *AuthService) RevokeSession(userUuid, sessionUuid uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.sessionsColl().DeleteOne(ctx, bson.M{
		"user_uuid":    userUuid,
		"session_uuid": sessionUuid,
	})
	if err != nil {
		s.logger.Errorf("Failed to revoke session: %v", err)
		return err
	}
	return nil
}

func (s *AuthService) ValidateSession(sessionUuid uuid.UUID) bool {
	if s == nil || s.db == nil || sessionUuid == uuid.Nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	count, err := s.sessionsColl().CountDocuments(ctx, bson.M{
		"session_uuid": sessionUuid,
		"expires_at":   bson.M{"$gt": time.Now()},
	})
	if err != nil || count == 0 {
		return false
	}
	return true
}
