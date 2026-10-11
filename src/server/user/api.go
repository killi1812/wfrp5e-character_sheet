package user

import (
	"net/http"
	"strings"

	"github.com/killi1812/wfrp5e-character_sheet/app"
	"github.com/killi1812/wfrp5e-character_sheet/util/auth"
	"github.com/killi1812/wfrp5e-character_sheet/util/ginutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserCtn struct {
	UserCrud IUserCrudService
	logger   *zap.SugaredLogger
}

func NewUserCtn() app.Controller {
	var controller *UserCtn

	// Call dependency injection
	app.Invoke(func(UserService IUserCrudService, logger *zap.SugaredLogger) {
		// create controller
		controller = &UserCtn{
			UserCrud: UserService,
			logger:   logger,
		}
	})

	return controller
}

func (u *UserCtn) RegisterEndpoints(api *gin.RouterGroup) {
	// create a group with the name of the router
	group := api.Group("/user")

	// Public unauthenticated endpoints
	group.POST("/register", u.register)

	// Protected endpoint
	group.GET("/my-data", auth.Protect(), u.getLoggedInUser)

	// Admin-only endpoints
	adminGroup := group.Group("")
	adminGroup.Use(auth.Protect(string(ROLE_ADMIN)))
	adminGroup.POST("/", u.create)
	adminGroup.GET("/:uuid", u.get)
	adminGroup.PUT("/:uuid", u.update)
	adminGroup.DELETE("/:uuid", u.delete)
	adminGroup.GET("/all-users", u.getAllUsersForSuperAdmin)
	adminGroup.GET("/search", u.searchUsersByName)
}

// UserExample godoc
//
//	@Summary		get user with uuid
//	@Description	get a user with uuid
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	UserDto
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string	true	"user uuid"
//	@Router			/user/{uuid} [get]
func (u *UserCtn) get(c *gin.Context) {
	userUuid, ok := ginutil.ParseUUIDParam(c, "uuid")
	if !ok {
		u.logger.Errorf("error parsing uuid value = %s", c.Param("uuid"))
		return
	}

	user, err := u.UserCrud.Read(userUuid)
	if ginutil.HandleServiceError(c, err, ErrRecordNotFound) {
		return
	}

	uDto := UserDto{}
	c.JSON(http.StatusOK, uDto.FromModel(user))
}

// Register godoc
//
//	@Summary	Register a new user
//	@Tags		user
//	@Produce	json
//	@Success	201	{object}	UserDto
//	@Failure	400
//	@Failure	409
//	@Failure	500
//	@Param		model	body	RegisterUserDto	true	"Data for registering user"
//	@Router		/user/register [post]
func (u *UserCtn) register(c *gin.Context) {
	var inputDto RegisterUserDto
	if err := c.BindJSON(&inputDto); err != nil {
		u.logger.Errorf("Failed to bind register error = %+v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration data. Username must be at least 2 characters and password at least 6 characters."})
		return
	}

	newUser := inputDto.ToModel()

	user, err := u.UserCrud.Create(newUser, inputDto.Password)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		u.logger.Errorf("Failed to register user error = %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	responseDto := UserDto{}.FromModel(user)
	c.JSON(http.StatusCreated, responseDto)
}

// UserExample godoc
//
//	@Summary	Create new user
//	@Tags		user
//	@Produce	json
//	@Success	201	{object}	UserDto
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Param		model	body	NewUserDto	true	"Data for new user"
//	@Router		/user [post]
func (u *UserCtn) create(c *gin.Context) {
	var inputDto NewUserDto
	if err := c.BindJSON(&inputDto); err != nil {
		u.logger.Errorf("Failed to bind error = %+v", err)
		return
	}

	newUser, err := inputDto.ToModel()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Create(newUser, inputDto.Password)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Log the response for debugging
	responseDto := inputDto.FromModel(user)
	tokenValue := ""

	u.logger.Infof("Response DTO: %+v", responseDto)
	u.logger.Infof("Response PoliceToken: %q", tokenValue)

	c.JSON(http.StatusCreated, responseDto)
}

// UserExample godoc
//
//	@Summary	Update user with new dat
//	@Tags		user
//	@Produce	json
//	@Success	200	{object}	UserDto
//	@Failure	400
//	@Failure	404
//	@Failure	500
//	@Param		uuid	path	string	true	"uuid of user to be updated"
//	@Param		model	body	UserDto	true	"Data for updating user"
//	@Router		/user/{uuid} [put]
func (u *UserCtn) update(c *gin.Context) {
	userUuid, ok := ginutil.ParseUUIDParam(c, "uuid")
	if !ok {
		u.logger.Errorf("Error parsing UUID = %s", c.Param("uuid"))
		return
	}

	var inputDto UserDto
	if err := c.BindJSON(&inputDto); err != nil {
		u.logger.Errorf("Failed to bind error = %+v", err)
		return
	}

	newUser, err := inputDto.ToModel()
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := u.UserCrud.Update(userUuid, newUser)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, inputDto.FromModel(user))
}

// UserExample  godoc
//
//	@Summary		delete user with uuid
//	@Description	delete a user with uuid
//	@Tags			user
//	@Produce		json
//	@Success		204
//	@Failure		400
//	@Failure		404
//	@Failure		500
//	@Param			uuid	path	string	true	"user uuid"
//	@Router			/user/{uuid} [delete]
func (u *UserCtn) delete(c *gin.Context) {
	userUuid, ok := ginutil.ParseUUIDParam(c, "uuid")
	if !ok {
		u.logger.Errorf("error parsing uuid value = %s", c.Param("uuid"))
		return
	}

	err := u.UserCrud.Delete(userUuid)
	if ginutil.HandleServiceError(c, err, ErrRecordNotFound) {
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetLoggedInUser godoc
//
//	@Summary		Get logged-in user data
//	@Description	Fetches the currently logged-in user's data based on the JWT token
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	UserDto
//	@Failure		400
//	@Failure		401
//	@Failure		404
//	@Failure		500
//	@Router			/user/my-data [get]
func (u *UserCtn) getLoggedInUser(c *gin.Context) {
	_, claims, err := auth.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		u.logger.Errorf("Failed to parse token: %v", err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	userUuid := claims.UserUuid
	if userUuid == uuid.Nil {
		var parseErr error
		userUuid, parseErr = uuid.Parse(claims.ID)
		if parseErr != nil {
			u.logger.Errorf("Error parsing UUID = %s", parseErr)
			c.AbortWithError(http.StatusBadRequest, parseErr)
			return
		}
	}

	user, err := u.UserCrud.Read(userUuid)
	if ginutil.HandleServiceError(c, err, ErrRecordNotFound) {
		return
	}

	userDto := UserDto{}
	c.JSON(http.StatusOK, userDto.FromModel(user))
}

// GetAllUsers godoc
//
//	@Summary		Get all users for admin
//	@Description	Fetches all users for admin
//	@Tags			user
//	@Produce		json
//	@Success		200	{array}	UserDto
//	@Failure		401
//	@Failure		403
//	@Failure		500
//	@Router			/user/all-users [get]
func (u *UserCtn) getAllUsersForSuperAdmin(c *gin.Context) {
	_, claims, err := auth.ParseToken(c.Request.Header.Get("Authorization"))
	if err != nil {
		u.logger.Errorf("Failed to parse token: %v", err)
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	if claims.Role != string(ROLE_ADMIN) {
		u.logger.Warnf("Unauthorized access attempt by user with role: %s", claims.Role)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	users, err := u.UserCrud.GetAllUsers()
	if err != nil {
		u.logger.Errorf("Failed to fetch users: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userDtos := make([]UserDto, 0, len(users))
	for _, user := range users {
		itemDto := UserDto{}
		userDtos = append(userDtos, itemDto.FromModel(&user))
	}

	c.JSON(http.StatusOK, userDtos)
}

// SearchUsersByName godoc
//
//	@Summary		Search users by name
//	@Description	Performs a fuzzy search for users by first name, last name, or full name with similarity matching
//	@Tags			user
//	@Produce		json
//	@Param			query	query	string	true	"Search query"
//	@Success		200		{array}	UserDto
//	@Failure		400
//	@Failure		500
//	@Router			/user/search [get]
func (u *UserCtn) searchUsersByName(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		u.logger.Warn("Search query is empty")
		c.JSON(http.StatusBadRequest, "Search query is required")
		return
	}

	u.logger.Infof("Searching users with query: %s", query)

	users, err := u.UserCrud.SearchUsersByName(query)
	if err != nil {
		u.logger.Errorf("Failed to search users: %v", err)
		c.JSON(http.StatusInternalServerError, "Failed to search users")
		return
	}

	userDtos := make([]UserDto, 0, len(users))
	for _, user := range users {
		itemDto := UserDto{}
		userDtos = append(userDtos, itemDto.FromModel(&user))
	}

	c.JSON(http.StatusOK, userDtos)
}
