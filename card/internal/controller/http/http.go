package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tomatoCoderq/card/internal/services/card"
	"github.com/tomatoCoderq/card/pkg/model"
	"github.com/tomatoCoderq/card/pkg/scheme"
	"github.com/tomatoCoderq/card/internal/controller"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

func GetUserIdFromContext(ctx *gin.Context) (uuid.UUID, error) {
	userClaims, exists := ctx.Get("userClaims")
	if !exists {
		return uuid.UUID{}, fmt.Errorf("user claims do not exist")
	}

	claimsMap, ok := userClaims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot convert claims to map")
	}

	fmt.Println("VIUST")

	userIdString, ok := claimsMap["uid"].(string)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("cannot get user_id from claims map")
	}

	userId, err := uuid.Parse(userIdString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("cannot parse uuid")
	}
	
	return userId, nil
}

type CardController struct {
	CardService controller.Card
}

func CreateNewController(cardService *services.Card) *CardController {
	return &CardController{CardService: cardService}
}

func (cc CardController) AddCard(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var card model.Card
	if err = json.Unmarshal(body, &card); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	card.CreatedBy = userId

	response, err := cc.CardService.AddCard(&card)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (cc CardController) ReadAllCardsToLearn(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	fmt.Println("USS", userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, err := cc.CardService.ReadAllCards(userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (cc CardController) UpdateCard(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	
	userId, err := GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	
	var cardUpdate schemes.UpdateCardScheme
	if err = ctx.ShouldBindJSON(&cardUpdate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	fmt.Println("COMEHERE")


	card, err := cc.CardService.UpdateCard(cardId, &cardUpdate, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("CO", card)
	ctx.JSON(http.StatusOK, card)
}

func (cc CardController) DeleteCard(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = cc.CardService.DeleteCard(cardId, userId)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

func (cc CardController) AddAnswers(ctx *gin.Context) {
	var answers []schemes.AnswerScheme

	if err := ctx.ShouldBindJSON(&answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userId, err := GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if err = cc.CardService.AddAnswers(userId, answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "added answers succesfully "})
}
