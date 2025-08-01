package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	_ "github.com/swaggo/swag/example/celler/httputil"

	_ "github.com/GOeda-Co/proto-contract/model/user"
	// "github.com/tomatoCoderq/repeatro/pkg/schemes"
	schemes "github.com/GOeda-Co/proto-contract/scheme/sso"
)

// Register godoc
//
//	@Summary		Registers new user to the system
//	@Description	Register by email, name, and password, getting user_id
//	@Tags			sso
//	@Accept			json
//	@Produce		json
//	@Param			request	body		schemes.RegisterScheme	true	"Registration data"
//	@Success		200		{object}	model.RegisterResponse	"User registered successfully"
//	@Failure		400		{object}	model.ErrorResponse		"Bad Request - Invalid request body"
//	@Failure		500		{object}	model.ErrorResponse		"Internal Server Error - Failed to register user"
//	@Router			/register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var registerScheme schemes.RegisterScheme

	if err := ctx.ShouldBindBodyWithJSON(&registerScheme); err != nil {
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
		return
	}

	uid, err := c.ssoClient.Register(ctx.Request.Context(), registerScheme.Email, registerScheme.Password, registerScheme.Name)
	if err != nil {
		c.log.Debug("Failed to register user", "error", err)
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Failed to register user: %v", err)})
		return
	}
	ctx.JSON(200, gin.H{
		"user_id": uid,
		"message": "User registered successfully",
	})
}

// Login godoc
//
//	@Summary		Logs in a user
//	@Description	Logs in a user and returns a JWT token
//	@Tags			sso
//	@Accept			json
//	@Produce		json
//	@Param			request	body		schemes.LoginScheme	true	"Login credentials"
//	@Success		200		{object}	model.LoginResponse		"User logged in successfully"
//	@Failure		400		{object}	model.ErrorResponse		"Bad Request - Invalid request body"
//	@Failure		500		{object}	model.ErrorResponse		"Internal Server Error - Failed to login user"
//	@Router			/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var loginScheme schemes.LoginScheme

	if err := ctx.ShouldBindBodyWithJSON(&loginScheme); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := c.ssoClient.Login(ctx.Request.Context(), loginScheme.Email, loginScheme.Password, loginScheme.AppId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Failed to login user: %v", err)})
		return
	}
	ctx.JSON(200, gin.H{
		"token":   token,
		"message": "User logged in successfully",
	})
}

// IsAdmin godoc
//
//	@Summary		Checks if user is admin
//	@Description	Verifies if the user has admin privileges
//	@Tags			sso
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		string	true	"ID of user"
//	@Success		200		{object}	model.AdminCheckResponse	"Admin status check result"
//	@Failure		400		{object}	model.ErrorResponse			"Bad Request - user_id is required or invalid format"
//	@Failure		500		{object}	model.ErrorResponse			"Internal Server Error - Failed to check admin status"
//	@Router			/is-admin [get]
func (c *Controller) IsAdmin(ctx *gin.Context) {
	userId := ctx.Query("user_id")

	if userId == "" {
		ctx.JSON(400, gin.H{"error": "user_id is required"})
		return
	}

	userIdParse, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid user_id format"})
		return
	}

	isAdmin, err := c.ssoClient.IsAdmin(ctx.Request.Context(), userIdParse)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Failed to check admin status: %v", err)})
		return
	}

	ctx.JSON(200, gin.H{
		"is_admin": isAdmin,
	})
}
