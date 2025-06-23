package http

import (
	// "fmt"
	// "fmt"
	"net/http"

	"repeatro/src/user/internal/service"
	"repeatro/src/user/pkg/scheme"

	"repeatro/internal/security"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *services.Service
	security    *security.Security
}

func CreateNewUserController(userService *services.Service, security *security.Security) *UserController {
	return &UserController{
		UserService: userService,
		security:    security,
	}
}

func (uc *UserController) Register(ctx *gin.Context) {
	var userRegister schemes.AuthUser

	if err := ctx.ShouldBindJSON(&userRegister); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	token, err := uc.UserService.Register(userRegister)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(200, token)
}

func (uc *UserController) Login(ctx *gin.Context) {
	// get email and password
	var userLogin schemes.AuthUser

	if err := ctx.ShouldBindJSON(&userLogin); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	token, err := uc.UserService.Login(userLogin)
	if err != nil {
		ctx.AbortWithError(500, err)
		return
	}

	ctx.JSON(200, token)
}
