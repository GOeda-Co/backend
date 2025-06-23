package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"repeatro/internal/tools"
	"repeatro/src/deck/internal/service"
	"repeatro/src/deck/pkg/model"

	"github.com/google/uuid"
)


type DeckController struct {
	DeckService *services.Service
}

func CreateNewController(deckService *services.Service) *DeckController {
	return &DeckController{DeckService: deckService}
}

func (dc DeckController) AddDeck(ctx *gin.Context) {
	var deck models.Deck

	if err := ctx.ShouldBindJSON(&deck); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deck.CreatedBy = userId
	createdDeck, err := dc.DeckService.AddCard(&deck, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdDeck)
}

func (dc DeckController) ReadAllDecks(ctx *gin.Context) {
	userId, err := tools.GetUserIdFromHeader(ctx)
	fmt.Println("SS", userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	decks, err := dc.DeckService.ReadAllDecksOfUser(userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, decks)
}

func (dc DeckController) ReadDeck(ctx *gin.Context) {
	deckId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deck, err := dc.DeckService.ReadDeck(deckId, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deck)
}

func (dc DeckController) ReadCardsFromDeck(ctx *gin.Context) {
	deckId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cards, err := dc.DeckService.ReadAllCardsFromDeck(deckId, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cards)
}

// Delete
func (dc DeckController) DeleteDeck(ctx *gin.Context) {
	deckId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := dc.DeckService.DeleteDeck(deckId, userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (dc DeckController) AddCardToDeck(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("card_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	deckId, err := uuid.Parse(ctx.Param("deck_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	userId, err := tools.GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := dc.DeckService.AddCardToDeck(cardId, deckId, userId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "card added to deck"})
}
