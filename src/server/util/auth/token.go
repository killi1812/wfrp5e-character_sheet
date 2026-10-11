package auth

import (
	"errors"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Claims struct {
	jwt.RegisteredClaims
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	TokenUuid uuid.UUID `json:"uuid"`
	UserUuid  uuid.UUID `json:"userUuid,omitempty"`
}

const (
	_ACCESS_TOKEN_DURATION  = 24 * time.Hour
	_REFRESH_TOKEN_DURATION = 7 * 24 * time.Hour
)

func ParseToken(authHeader string) (*jwt.Token, *Claims, error) {
	if len(authHeader) <= len("Bearer ") || authHeader[:len("Bearer ")] != "Bearer " {
		zap.S().Debugf("token: %s", authHeader)
		return nil, nil, cerror.ErrInvalidTokenFormat
	}
	tokenString := authHeader[len("Bearer "):]
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(app.AccessKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if claims.UserUuid == uuid.Nil && claims.ID != "" {
		if parsed, err := uuid.Parse(claims.ID); err == nil {
			claims.UserUuid = parsed
		}
	}

	return token, &claims, nil
}

// ParseTokenAllowExpired parses a Bearer JWT token, verifying its signature while allowing expired tokens.
func ParseTokenAllowExpired(authHeader string) (*jwt.Token, *Claims, error) {
	if len(authHeader) <= len("Bearer ") || authHeader[:len("Bearer ")] != "Bearer " {
		zap.S().Debugf("token: %s", authHeader)
		return nil, nil, cerror.ErrInvalidTokenFormat
	}
	tokenString := authHeader[len("Bearer "):]
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(app.AccessKey), nil
	})
	if err != nil {
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorSignatureInvalid == 0 &&
				validationErr.Errors&jwt.ValidationErrorMalformed == 0 &&
				validationErr.Errors&jwt.ValidationErrorUnverifiable == 0 {
				if claims.UserUuid == uuid.Nil && claims.ID != "" {
					if parsed, err := uuid.Parse(claims.ID); err == nil {
						claims.UserUuid = parsed
					}
				}
				return token, &claims, nil
			}
		}
		return nil, nil, err
	}

	if claims.UserUuid == uuid.Nil && claims.ID != "" {
		if parsed, err := uuid.Parse(claims.ID); err == nil {
			claims.UserUuid = parsed
		}
	}

	return token, &claims, nil
}

// GenerateTokens returns a jwt access token and refresh token or an error
func GenerateTokens(email, username, role string, userUuid uuid.UUID) (string, string, error) {
	if userUuid == uuid.Nil {
		return "", "", cerror.ErrUserIsNil
	}
	uuidPair := uuid.New()
	accessTokenClaims := &Claims{
		Email:     email,
		Username:  username,
		Role:      role,
		TokenUuid: uuidPair,
		UserUuid:  userUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_ACCESS_TOKEN_DURATION)),
			ID:        userUuid.String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(app.AccessKey))
	if err != nil {
		zap.S().Errorf("Failed to generate access token err = %w", err)
		return "", "", err
	}

	refreshTokenClaims := &Claims{
		Username:  username,
		Email:     email,
		Role:      role,
		TokenUuid: uuidPair,
		UserUuid:  userUuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(_REFRESH_TOKEN_DURATION)),
			ID:        userUuid.String(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(app.RefreshKey))
	if err != nil {
		zap.S().Errorf("Failed to generate refresh token err = %w", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
