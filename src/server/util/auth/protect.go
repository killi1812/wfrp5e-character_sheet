package auth

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SessionValidatorFunc represents a function that verifies if a session is currently active.
type SessionValidatorFunc func(sessionUuid uuid.UUID) bool

var _sessionValidator SessionValidatorFunc

// SetSessionValidator configures a real-time session validation callback.
func SetSessionValidator(validator SessionValidatorFunc) {
	_sessionValidator = validator
}

// Protect protects routes allowing access only to given roles
// if roles are empty it only checks for the validity of tokens
func Protect(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		token, claims, err := ParseToken(authHeader)
		if err != nil {
			zap.S().Infof("Auth failed with err = %+v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if len(roles) != 0 && !slices.Contains(roles, claims.Role) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if _sessionValidator != nil && claims.SessionUuid != uuid.Nil {
			if !_sessionValidator(claims.SessionUuid) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session revoked or expired"})
				return
			}
		}

		c.Set("claims", claims)
		c.Set("userUuid", claims.UserUuid)
		c.Set("role", claims.Role)
		c.Set("sessionUuid", claims.SessionUuid)

		c.Next()
	}
}
