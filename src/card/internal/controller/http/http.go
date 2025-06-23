package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"repeatro/src/card/internal/service"
	"repeatro/src/card/pkg/model"
	"repeatro/src/card/pkg/scheme"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"repeatro/internal/tools"
)

type CardController struct {
	CardService services.CardServiceInterface
}

func CreateNewController(cardService *services.CardService) *CardController {
	return &CardController{CardService: cardService}
}

func (cc CardController) AddCard(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
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
	userId, err := tools.GetUserIdFromHeader(ctx)
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
	
	userId, err := tools.GetUserIdFromHeader(ctx)
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

	userId, err := tools.GetUserIdFromHeader(ctx)
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

	userId, err := tools.GetUserIdFromHeader(ctx)
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
