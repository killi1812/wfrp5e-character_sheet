package auth_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/user"
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/cerror"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	app.AccessKey = "test-jwt-key"
	app.RefreshKey = "test-refresh-key"

	now := time.Now()
	accessTokenDuration := 5 * time.Minute
	validClaims := auth.Claims{
		Username: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenDuration)),
			ID:        "test-uuid",
		},
	}
	validToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	validTokenString, _ := validToken.SignedString([]byte(app.AccessKey))
	validAuthHeader := "Bearer " + validTokenString

	invalidSignatureToken := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	invalidSignatureTokenString, _ := invalidSignatureToken.SignedString([]byte("wrong-key"))
	invalidSignatureAuthHeader := "Bearer " + invalidSignatureTokenString

	expiredClaims := auth.Claims{
		Username: "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
			ID:        "expired-uuid",
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(app.AccessKey))
	expiredAuthHeader := "Bearer " + expiredTokenString

	tests := []struct {
		name       string
		authHeader string
		wantToken  bool
		wantClaims *auth.Claims
		wantErr    error
	}{
		{
			name:       "Valid token",
			authHeader: validAuthHeader,
			wantToken:  true,
			wantClaims: &validClaims,
			wantErr:    nil,
		},
		{
			name:       "Invalid token format - missing Bearer",
			authHeader: validTokenString,
			wantToken:  false,
			wantClaims: nil,
			wantErr:    cerror.ErrInvalidTokenFormat,
		},
		{
			name:       "Invalid token format - too short",
			authHeader: "Bearer",
			wantToken:  false,
			wantClaims: nil,
			wantErr:    cerror.ErrInvalidTokenFormat,
		},
		{
			name:       "Invalid signature",
			authHeader: invalidSignatureAuthHeader,
			wantToken:  false,
			wantClaims: nil,
			wantErr:    jwt.ErrSignatureInvalid,
		},
		{
			name:       "Expired token",
			authHeader: expiredAuthHeader,
			wantToken:  false,
			wantClaims: nil,
			wantErr:    jwt.ErrTokenExpired,
		},
		{
			name:       "Malformed token",
			authHeader: "Bearer malformed.token.string",
			wantToken:  false,
			wantClaims: nil,
			wantErr:    jwt.ErrTokenMalformed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, gotClaims, gotErr := auth.ParseToken(tt.authHeader)

			if gotErr != nil {
				if tt.wantErr == nil || !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("ParseToken() error = %v, wantErr %v", gotErr, tt.wantErr)
				}
				return
			}
			if tt.wantErr != nil {
				t.Errorf("ParseToken() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if (gotToken != nil) != tt.wantToken {
				t.Errorf("ParseToken() gotToken = %v, wantToken %v", gotToken != nil, tt.wantToken)
			}

			if tt.wantClaims != nil {
				if gotClaims == nil {
					t.Errorf("ParseToken() gotClaims = nil, want non-nil")
					return
				}
				if !reflect.DeepEqual(gotClaims, tt.wantClaims) {
					t.Errorf("ParseToken() gotClaims = %v, wantClaims %v", gotClaims, tt.wantClaims)
				}
			} else if gotClaims != nil {
				t.Errorf("ParseToken() gotClaims = %v, want nil", gotClaims)
			}
		})
	}
}

func TestParseTokenAllowExpired(t *testing.T) {
	app.AccessKey = "test-jwt-key"

	now := time.Now()
	expiredClaims := auth.Claims{
		Username: "expired@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour)),
			ID:        "expired-uuid",
		},
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(app.AccessKey))
	expiredAuthHeader := "Bearer " + expiredTokenString

	token, claims, err := auth.ParseTokenAllowExpired(expiredAuthHeader)
	if err != nil {
		t.Fatalf("ParseTokenAllowExpired() returned error for expired token: %v", err)
	}
	if claims == nil || claims.Username != "expired@example.com" {
		t.Fatalf("ParseTokenAllowExpired() claims mismatch: %+v", claims)
	}
	if token == nil {
		t.Fatalf("ParseTokenAllowExpired() returned nil token")
	}

	// Invalid signature should still error
	invalidSigTokenString, _ := expiredToken.SignedString([]byte("wrong-key"))
	_, _, err = auth.ParseTokenAllowExpired("Bearer " + invalidSigTokenString)
	if err == nil {
		t.Fatalf("ParseTokenAllowExpired() expected error for invalid signature, got nil")
	}

	// Missing / invalid format
	_, _, err = auth.ParseTokenAllowExpired("invalid")
	assert.Error(t, err)

	_, _, err = auth.ParseTokenAllowExpired("Bearer short")
	assert.Error(t, err)

	// Malformed
	_, _, err = auth.ParseTokenAllowExpired("Bearer not.a.valid.jwt")
	assert.Error(t, err)

	// Valid non-expired token
	validClaims := auth.Claims{
		Username: "valid@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			ID:        "valid-uuid",
		},
	}
	validTok := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims)
	validTokStr, _ := validTok.SignedString([]byte(app.AccessKey))
	_, vClaims, err := auth.ParseTokenAllowExpired("Bearer " + validTokStr)
	assert.NoError(t, err)
	assert.Equal(t, "valid@example.com", vClaims.Username)
}

func TestParseToken_UserUuid(t *testing.T) {
	app.AccessKey = "test-jwt-key"
	expectedUserUuid := uuid.New()

	// Token generated via GenerateTokens has UserUuid populated
	accToken, _, err := auth.GenerateTokens("test@example.com", "testuser", "user", expectedUserUuid)
	assert.NoError(t, err)

	_, claims, err := auth.ParseToken("Bearer " + accToken)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserUuid, claims.UserUuid)

	// Legacy token with only RegisteredClaims.ID populated gets parsed into UserUuid
	legacyClaims := auth.Claims{
		Username: "legacy@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			ID:        expectedUserUuid.String(),
		},
	}
	legacyToken := jwt.NewWithClaims(jwt.SigningMethodHS256, legacyClaims)
	legacyTokenString, _ := legacyToken.SignedString([]byte(app.AccessKey))

	_, parsedLegacyClaims, err := auth.ParseToken("Bearer " + legacyTokenString)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserUuid, parsedLegacyClaims.UserUuid)
}

func TestGenerateTokens(t *testing.T) {
	app.AccessKey = "test-jwt-key"
	app.RefreshKey = "test-refresh-key"

	usr := &user.User{
		Username: "test@example.com",
		Uuid:     uuid.New(),
	}

	accessTokenDuration := 15 * time.Minute
	refreshTokenDuration := 30 * 24 * time.Hour

	tests := []struct {
		name                     string
		user                     *user.User
		wantAccessTokenNonEmpty  bool
		wantRefreshTokenNonEmpty bool
		wantErr                  bool
	}{
		{
			name:                     "Valid user",
			user:                     usr,
			wantAccessTokenNonEmpty:  true,
			wantRefreshTokenNonEmpty: true,
			wantErr:                  false,
		},
		{
			name:                     "Nil user",
			user:                     nil,
			wantAccessTokenNonEmpty:  false,
			wantRefreshTokenNonEmpty: false,
			wantErr:                  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var email, username, role string
			var userUuid uuid.UUID
			if tt.user != nil {
				email = tt.user.Email
				username = tt.user.Username
				role = string(tt.user.Role)
				userUuid = tt.user.Uuid
			}
			customSessionUuid := uuid.New()
			gotAccessToken, gotRefreshToken, err := auth.GenerateTokens(email, username, role, userUuid, customSessionUuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokens() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantAccessTokenNonEmpty && gotAccessToken == "" {
				t.Errorf("GenerateTokens() gotAccessToken = %q, want non-empty", gotAccessToken)
			}
			if tt.wantRefreshTokenNonEmpty && gotRefreshToken == "" {
				t.Errorf("GenerateTokens() gotRefreshToken = %q, want non-empty", gotRefreshToken)
			}

			if gotAccessToken != "" {
				token, _, err := new(jwt.Parser).ParseUnverified(gotAccessToken, &auth.Claims{})
				if err != nil {
					t.Errorf("GenerateTokens() generated invalid access token: %+v", err)
					return
				}
				if claims, ok := token.Claims.(*auth.Claims); ok {
					if claims.Username != tt.user.Username || claims.ID != tt.user.Uuid.String() || claims.UserUuid != tt.user.Uuid || claims.SessionUuid != customSessionUuid {
						t.Errorf("GenerateTokens() access token claims are incorrect: got = %+v, want = %+v", claims, tt.user)
					}
					if !claims.ExpiresAt.After(time.Now().Add(accessTokenDuration-time.Minute)) || !claims.ExpiresAt.Before(time.Now().Add(accessTokenDuration+time.Minute)) {
						t.Errorf("GenerateTokens() access token expiry is not within expected range: %+v", claims.ExpiresAt)
					}
				} else {
					t.Errorf("GenerateTokens() failed to parse access token claims")
				}
			}

			if gotRefreshToken != "" {
				token, _, err := new(jwt.Parser).ParseUnverified(gotRefreshToken, &auth.Claims{})
				if err != nil {
					t.Errorf("GenerateTokens() generated invalid refresh token: %+v", err)
					return
				}
				if claims, ok := token.Claims.(*auth.Claims); ok {
					if claims.Username != tt.user.Username || claims.ID != tt.user.Uuid.String() || claims.UserUuid != tt.user.Uuid || claims.SessionUuid != customSessionUuid {
						t.Errorf("GenerateTokens() refresh token claims are incorrect: got = %+v, want = %+v", claims, tt.user)
					}
					if !claims.ExpiresAt.After(time.Now().Add(refreshTokenDuration-time.Minute)) || !claims.ExpiresAt.Before(time.Now().Add(refreshTokenDuration+time.Minute)) {
						t.Errorf("GenerateTokens() refresh token expiry is not within expected range: %+v", claims.ExpiresAt)
					}
				} else {
					t.Errorf("GenerateTokens() failed to parse refresh token claims")
				}
			}
		})
	}
}

func TestHashToken(t *testing.T) {
	t1 := "sample-token-123"
	h1 := auth.HashToken(t1)
	h2 := auth.HashToken(t1)
	assert.NotEmpty(t, h1)
	assert.Equal(t, h1, h2)
	assert.NotEqual(t, h1, auth.HashToken("different-token"))
}

func TestParseRefreshToken(t *testing.T) {
	app.AccessKey = "test-access-key"
	app.RefreshKey = "test-refresh-key"

	userUuid := uuid.New()
	sessionUuid := uuid.New()
	_, refreshToken, err := auth.GenerateTokens("u@test.com", "u", "user", userUuid, sessionUuid)
	assert.NoError(t, err)

	t.Run("Valid raw refresh token", func(t *testing.T) {
		token, claims, err := auth.ParseRefreshToken(refreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotNil(t, claims)
		assert.Equal(t, userUuid, claims.UserUuid)
		assert.Equal(t, sessionUuid, claims.SessionUuid)
	})

	t.Run("Valid Bearer prefixed refresh token", func(t *testing.T) {
		token, claims, err := auth.ParseRefreshToken("Bearer " + refreshToken)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotNil(t, claims)
		assert.Equal(t, userUuid, claims.UserUuid)
	})

	t.Run("Empty refresh token", func(t *testing.T) {
		_, _, err := auth.ParseRefreshToken("")
		assert.ErrorIs(t, err, cerror.ErrInvalidTokenFormat)
	})

	t.Run("Malformed refresh token", func(t *testing.T) {
		_, _, err := auth.ParseRefreshToken("malformed-token")
		assert.Error(t, err)
	})

	t.Run("Token signed with wrong key", func(t *testing.T) {
		// access token is signed with AccessKey, not RefreshKey
		accessToken, _, _ := auth.GenerateTokens("u@test.com", "u", "user", userUuid, sessionUuid)
		_, _, err := auth.ParseRefreshToken(accessToken)
		assert.Error(t, err)
	})
}

