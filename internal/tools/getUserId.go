package tools

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUserIdFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return uuid.UUID{}, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot convert claims to map")
	}

	userIdString, ok := claimsMap["user_id"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetUserIdFromHeader(ctx *gin.Context) (uuid.UUID, error) {
	userClaims := ctx.GetHeader("userClaims")

	var claimsMap jwt.MapClaims
	err := json.Unmarshal([]byte(userClaims), &claimsMap)
	if err != nil {
		return uuid.UUID{}, err
	}

	userIdString, ok := claimsMap["user_id"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}

func GetClaimsMapFromContext(ctx *gin.Context) (map[string]interface{}, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return nil, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("cannot convert claims to map")
	}

	return claimsMap, nil
}

func GetUserIdFromClaims(userClaims any) (uuid.UUID, error) {
	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot convert claims to map")
	}

	userIdString, ok := claimsMap["user_id"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, err
	}
	return userId, nil
}
