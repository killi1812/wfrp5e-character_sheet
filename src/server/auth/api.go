package auth

import (
	"net/http"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	authUtil "github.com/killi1812/wfrp5e-character_sheet/util/auth"

	"github.com/gin-gonic/gin"
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

	group.POST("/login", c.login)
	group.POST("/refresh", authUtil.Protect(), c.refreshToken)
	group.POST("/logout", authUtil.Protect(), c.logout)
}

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticates a user and returns access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginDto	body		LoginDto	true	"Login credentials"
//	@Success		200			{object}	TokenDto
//	@Router			/auth/login [post]
func (l *AuthCtn) login(c *gin.Context) {
	var loginDto LoginDto

	if err := c.BindJSON(&loginDto); err != nil {
		l.logger.Errorf("Invalid login request err = %+v", err)
		return
	}

	accessToken, err := l.auth.Login(loginDto.Username, loginDto.Password)
	if err != nil {
		l.logger.Errorf("Login failed err = %+v", err)
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, TokenDto{
		AccessToken: accessToken,
	})
}

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Generates a new access token using a valid refresh token
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	TokenDto
//	@Router			/auth/refresh [post]
func (l *AuthCtn) refreshToken(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authorization")
	token, err := l.auth.RefreshTokens(tokenStr)
	if err != nil {
		l.logger.Errorf("Refresh failed err = %w", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, TokenDto{
		AccessToken: token,
	})
}

// Logout godoc
//
//	@Summary		User Logout
//	@Description	Logs out current user session
//	@Tags			auth
//	@Produce		json
//	@Success		200
//	@Router			/auth/logout [post]
func (l *AuthCtn) logout(c *gin.Context) {
	_, claims, err := authUtil.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		l.logger.Errorf("Logout failed err = %w", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = l.auth.Logout(claims.ID)
	if err != nil {
		l.logger.Errorf("Logout failed err = %w", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
