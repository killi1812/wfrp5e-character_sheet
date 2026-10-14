package auth

import (
	"net/http"
	"time"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/ginutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthCtn struct {
	auth   IAuthService
	logger *zap.SugaredLogger
}

func NewAuthCtn() app.Controller {
	var controller *AuthCtn

	app.Invoke(func(loginService IAuthService, logger *zap.SugaredLogger) {
		controller = &AuthCtn{
			auth:   loginService,
			logger: logger,
		}
	})

	return controller
}

func (c *AuthCtn) RegisterEndpoints(api *gin.RouterGroup) {
	group := api.Group("/auth")

	// Public routes
	group.POST("/login", c.login)
	group.POST("/refresh", c.refreshToken)

	// Protected routes
	protected := group.Group("", authUtil.Protect())
	protected.POST("/logout", c.logout)
	protected.POST("/logout-all", c.logoutAll)
	protected.GET("/sessions", c.listSessions)
	protected.DELETE("/sessions/:uuid", c.revokeSession)
}

func setRefreshTokenCookie(c *gin.Context, refreshToken string) {
	secure := gin.Mode() == gin.ReleaseMode || c.Request.TLS != nil
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, int(30*24*time.Hour/time.Second), "/api/auth", "", secure, true)
}

func clearRefreshTokenCookie(c *gin.Context) {
	secure := gin.Mode() == gin.ReleaseMode || c.Request.TLS != nil
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/api/auth", "", secure, true)
}

func getContextAuth(c *gin.Context) (uuid.UUID, uuid.UUID) {
	var userUuid, sessionUuid uuid.UUID
	if val, exists := c.Get("userUuid"); exists {
		if u, ok := val.(uuid.UUID); ok {
			userUuid = u
		}
	}
	if val, exists := c.Get("sessionUuid"); exists {
		if s, ok := val.(uuid.UUID); ok {
			sessionUuid = s
		}
	}
	if userUuid == uuid.Nil || sessionUuid == uuid.Nil {
		if val, exists := c.Get("claims"); exists {
			if claims, ok := val.(*authUtil.Claims); ok {
				if userUuid == uuid.Nil {
					userUuid = claims.UserUuid
				}
				if sessionUuid == uuid.Nil {
					sessionUuid = claims.SessionUuid
				}
			}
		}
	}
	return userUuid, sessionUuid
}

func (l *AuthCtn) login(c *gin.Context) {
	var loginDto LoginDto
	if err := c.BindJSON(&loginDto); err != nil {
		l.logger.Errorf("Invalid login request: %v", err)
		ginutil.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	accessToken, refreshToken, sessionUuid, err := l.auth.Login(
		loginDto.Username,
		loginDto.Password,
		loginDto.SessionId,
		userAgent,
		clientIP,
	)
	if err != nil {
		l.logger.Debugf("Login failed: %v", err)
		ginutil.RespondWithError(c, http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	setRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusOK, TokenDto{
		AccessToken: accessToken,
		SessionId:   sessionUuid,
	})
}

func (l *AuthCtn) refreshToken(c *gin.Context) {
	var tokenStr string
	if cookieVal, err := c.Cookie("refresh_token"); err == nil && cookieVal != "" {
		tokenStr = cookieVal
	} else if authHeader := c.Request.Header.Get("Authorization"); authHeader != "" {
		tokenStr = authHeader
	}

	if tokenStr == "" {
		ginutil.RespondWithError(c, http.StatusUnauthorized, "Missing refresh token")
		return
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	newAccessToken, newRefreshToken, sessionUuid, err := l.auth.RefreshTokens(tokenStr, userAgent, clientIP)
	if err != nil {
		l.logger.Warnf("Refresh failed: %v", err)
		clearRefreshTokenCookie(c)
		ginutil.RespondWithError(c, http.StatusUnauthorized, err.Error())
		return
	}

	setRefreshTokenCookie(c, newRefreshToken)

	c.JSON(http.StatusOK, TokenDto{
		AccessToken: newAccessToken,
		SessionId:   sessionUuid,
	})
}

func (l *AuthCtn) logout(c *gin.Context) {
	userUuid, sessionUuid := getContextAuth(c)
	if userUuid != uuid.Nil && sessionUuid != uuid.Nil {
		_ = l.auth.LogoutCurrent(userUuid, sessionUuid)
	}

	clearRefreshTokenCookie(c)
	c.Status(http.StatusNoContent)
}

func (l *AuthCtn) logoutAll(c *gin.Context) {
	userUuid, _ := getContextAuth(c)
	if userUuid != uuid.Nil {
		_ = l.auth.LogoutAll(userUuid)
	}

	clearRefreshTokenCookie(c)
	c.Status(http.StatusNoContent)
}

func (l *AuthCtn) listSessions(c *gin.Context) {
	userUuid, currentSessionUuid := getContextAuth(c)
	if userUuid == uuid.Nil {
		ginutil.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	sessions, err := l.auth.ListSessions(userUuid, currentSessionUuid)
	if err != nil {
		l.logger.Errorf("Failed to list sessions: %v", err)
		ginutil.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve sessions")
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (l *AuthCtn) revokeSession(c *gin.Context) {
	userUuid, currentSessionUuid := getContextAuth(c)
	if userUuid == uuid.Nil {
		ginutil.RespondWithError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	targetSessionUuid, ok := ginutil.ParseUUIDParam(c, "uuid")
	if !ok {
		return
	}

	err := l.auth.RevokeSession(userUuid, targetSessionUuid)
	if err != nil {
		l.logger.Errorf("Failed to revoke session: %v", err)
		ginutil.RespondWithError(c, http.StatusInternalServerError, "Failed to revoke session")
		return
	}

	if targetSessionUuid == currentSessionUuid {
		clearRefreshTokenCookie(c)
	}

	c.Status(http.StatusNoContent)
}
