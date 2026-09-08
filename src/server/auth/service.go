package auth

import (
	"context"
	"errors"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type IAuthService interface {
	Login(identifier, password string) (string, error)
	RefreshTokens(accessToken string) (string, error)
	Logout(userUuid string) error
}

type AuthService struct {
	db     *mongo.Database
	logger *zap.SugaredLogger
}

func NewAuthService() IAuthService {
	var service IAuthService

	app.Invoke(func(db *mongo.Database, logger *zap.SugaredLogger) {
		service = &AuthService{
			db:     db,
			logger: logger,
		}
	})

	return service
}

func (s *AuthService) usersColl() *mongo.Collection {
	return s.db.Collection("users")
}

func (s *AuthService) sessionsColl() *mongo.Collection {
	return s.db.Collection("sessions")
}

func (s *AuthService) Login(identifier, password string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var usr user.User
	err := s.usersColl().FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": identifier},
			{"email": identifier},
		},
	}).Decode(&usr)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			s.logger.Debugf("User not found Identifier = %s", identifier)
			return "", cerror.ErrInvalidCredentials
		}
		s.logger.Errorf("Failed to query user, error = %+v", err)
		return "", err
	}

	if !authUtil.VerifyPassword(usr.PasswordHash, password) {
		s.logger.Debugf("Invalid password for user Email: %s, uuid: %s", usr.Email, usr.Uuid)
		return "", cerror.ErrInvalidCredentials
	}

	token, refresh, err := authUtil.GenerateTokens(usr.Email, usr.Username, string(usr.Role), usr.Uuid)
	if err != nil {
		s.logger.Errorf("Failed to generate token error = %+v", err)
		return "", err
	}

	// Remove previous session if exists
	s.Logout(usr.Uuid.String())

	session := user.Session{
		UserUuid:     usr.Uuid,
		RefreshToken: refresh,
	}
	_, err = s.sessionsColl().InsertOne(ctx, session)
	if err != nil {
		s.logger.Errorf("Failed to create a session, err = %w", err)
		return "", err
	}

	return token, nil
}

func (s *AuthService) RefreshTokens(accessToken string) (string, error) {
	token, claims, err := authUtil.ParseToken(accessToken)
	if err != nil {
		s.logger.Errorf("Error Parsing claims err = %+v", err)
		return "", err
	}
	if !token.Valid {
		s.logger.Errorf("Token is not valid")
		return "", cerror.ErrInvalidTokenFormat
	}

	parsedUserUuid, err := uuid.Parse(claims.ID)
	if err != nil {
		s.logger.Errorf("Error Parsing uuid err = %+v", err)
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var session user.Session
	err = s.sessionsColl().FindOne(ctx, bson.M{"user_uuid": parsedUserUuid}).Decode(&session)
	if err != nil {
		s.logger.Errorf("Failed to find session, err = %w", err)
		return "", err
	}

	var refreshClaims authUtil.Claims
	_, err = jwt.ParseWithClaims(session.RefreshToken, &refreshClaims, func(token *jwt.Token) (any, error) {
		return []byte(app.RefreshKey), nil
	})
	if err != nil {
		s.logger.Errorf("Error parsing refresh token claims, err = %w", err)
		return "", err
	}

	if claims.TokenUuid != refreshClaims.TokenUuid {
		s.logger.Errorf("Error token uuids don't match, claims=%s, refresh=%s", claims.TokenUuid, refreshClaims.TokenUuid)
		return "", cerror.ErrInvalidTokenFormat
	}

	var usr user.User
	err = s.usersColl().FindOne(ctx, bson.M{"uuid": parsedUserUuid}).Decode(&usr)
	if err != nil {
		return "", err
	}

	newAccessToken, refreshToken, err := authUtil.GenerateTokens(usr.Email, usr.Username, string(usr.Role), usr.Uuid)
	if err != nil {
		s.logger.Errorf("Failed to generate tokens, err = %w", err)
		return "", err
	}

	session.RefreshToken = refreshToken
	_, err = s.sessionsColl().ReplaceOne(ctx, bson.M{"user_uuid": parsedUserUuid}, session)
	if err != nil {
		return "", err
	}

	return newAccessToken, nil
}

func (s *AuthService) Logout(userUuidStr string) error {
	s.logger.Debugf("logging out user with uuid = %s", userUuidStr)
	parsedUuid, err := uuid.Parse(userUuidStr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.sessionsColl().DeleteMany(ctx, bson.M{"user_uuid": parsedUuid})
	if err != nil {
		s.logger.Errorf("Error deleting session: %+v", err)
		return err
	}

	return nil
}
