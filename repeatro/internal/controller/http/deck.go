package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	_ "github.com/swaggo/swag/example/celler/httputil"

	// model "github.com/tomatoCoderq/repeatro/pkg/models"
	model "github.com/GOeda-Co/proto-contract/model/deck"
)

// AddDeck godoc
//
//	@Summary		Add a deck
//	@Description	Create a new deck
//	@Tags			decks
//	@Accept			json
//	@Produce		json
//	@Param			deck	body		model.Deck	true	"Deck to add"
//	@Success		200		{object}	model.Deck
//	@Failure		400		{object}	model.ErrorResponse	"Bad Request - Invalid request body"
//	@Failure		500		{object}	model.ErrorResponse	"Internal Server Error - Failed to get user ID or add deck"
//	@Router			/deck [post]
func (cc *Controller) AddDeck(ctx *gin.Context) {
	var deck model.Deck
	if err := ctx.ShouldBindJSON(&deck); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	deck.CreatedBy = userId

	response, err := cc.deckClient.AddDeck(ctx, &deck)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add deck"})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ReadAllDecks godoc
//
//	@Summary		Get all decks
//	@Description	Retrieves all decks in the system
//	@Tags			decks
//	@Produce		json
//	@Success		200	{array}		model.Deck
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to retrieve decks"
//	@Router			/decks [get]
func (cc *Controller) ReadAllDecks(ctx *gin.Context) {
	response, err := cc.deckClient.ReadAllDecks(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// ReadDeck godoc
//
//	@Summary		Get deck by ID
//	@Description	Retrieves a deck by its ID
//	@Tags			decks
//	@Param			id	path		string	true	"Deck ID"
//	@Success		200	{object}	model.Deck
//	@Failure		400	{object}	model.ErrorResponse	"Bad Request - Invalid deck ID format"
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to read deck"
//	@Router			/deck/{id} [get]
func (cc *Controller) ReadDeck(ctx *gin.Context) {
	deckId := ctx.Param("id")

	dId, err := uuid.Parse(deckId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	response, err := cc.deckClient.ReadDeck(ctx, dId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read deck"})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// SearchPublicDecks godoc
//
//	@Summary		Search public decks
//	@Description	Retrieves public decks. If user_id query parameter is provided, returns that user's public decks. Otherwise returns all public decks in the system.
//	@Tags			decks
//	@Produce		json
//	@Param			user_id	query		string	false	"User ID to filter by specific user's public decks"
//	@Success		200		{array}		model.Deck
//	@Failure		400		{object}	model.ErrorResponse	"Bad Request - Invalid user ID format"
//	@Failure		500		{object}	model.ErrorResponse	"Internal Server Error - Failed to get user ID or search public decks"
//	@Router			/decks/search [get]
func (cc *Controller) SearchPublicDecks(ctx *gin.Context) {
	userIdParam := ctx.Query("user_id")

	if userIdParam != "" {
		_, err := uuid.Parse(userIdParam)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			return
		}

		response, err := cc.deckClient.SearchUserPublicDecks(ctx, userIdParam)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search user public decks"})
			return
		}
		ctx.JSON(http.StatusOK, response)
	} else {
		// Search all public decks
		response, err := cc.deckClient.SearchAllPublicDecks(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search public decks"})
			return
		}
		ctx.JSON(http.StatusOK, response)
	}
}

// DeleteDeck godoc
//
//	@Summary		Delete a deck
//	@Description	Delete a deck by ID
//	@Tags			decks
//	@Param			id	path	string	true	"Deck ID"
//	@Success		200
//	@Failure		400	{object}	model.ErrorResponse	"Bad Request - Invalid deck ID format"
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to delete deck"
//	@Router			/deck/{id} [delete]
func (cc *Controller) DeleteDeck(ctx *gin.Context) {
	deckId := ctx.Param("id")

	dId, err := uuid.Parse(deckId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deck ID"})
		return
	}

	err = cc.deckClient.DeleteDeck(ctx, dId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deck"})
		return
	}
	ctx.Status(http.StatusOK)
}

// AddCardToDeck godoc
//
//	@Summary		Add card to deck
//	@Description	Add a card to a specific deck. Is_public will be updated to deck's is_public
//	@Tags			decks
//	@Param			card_id	path	string	true	"Card ID"
//	@Param			deck_id	path	string	true	"Deck ID"
//	@Success		200
//	@Failure		400	{object}	model.ErrorResponse	"Bad Request - Invalid card ID or deck ID format"
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to add card to deck"
//	@Router			/deck/{deck_id}/card/{card_id} [post]
func (cc *Controller) AddCardToDeck(ctx *gin.Context) {
	cid, err := uuid.Parse(ctx.Param("card_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid card ID"})
		return
	}

	did, err := uuid.Parse(ctx.Param("deck_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	err = cc.deckClient.AddCardToDeck(ctx, did, cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to deck"})
		return
	}
	ctx.Status(http.StatusOK)
}

// ReadCardsFromDeck godoc
//
//	@Summary		Get cards from deck
//	@Description	Retrieve all cards from a specific deck
//	@Tags			decks
//	@Param			id	path		string	true	"Deck ID"
//	@Success		200	{array}		model.Card
//	@Failure		400	{object}	model.ErrorResponse	"Bad Request - Invalid deck ID format"
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to get cards from deck"
//	@Router			/deck/{id}/cards [get]
func (cc *Controller) ReadCardsFromDeck(ctx *gin.Context) {
	did, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID format"})
		return
	}

	response, err := cc.deckClient.ReadCardsFromDeck(ctx, did)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cards from deck"})
		return
	}
	ctx.JSON(http.StatusOK, response)
}
