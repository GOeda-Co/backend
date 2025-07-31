package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	_ "github.com/swaggo/swag/example/celler/httputil"

	model "github.com/tomatoCoderq/repeatro/pkg/models"
	"github.com/tomatoCoderq/repeatro/pkg/schemes"
)

// AddCard godoc
//
//	@Summary		Add a card
//	@Description	Add a new card for the authenticated user
//	@Tags			cards
//	@Accept			json
//	@Produce		json
//	@Param			card	body		model.Card	true	"Card to add"
//	@Success		200		{object}	model.Card
//	@Failure		500		{object}	model.ErrorResponse	"Internal Server Error - Failed to read request body, get user ID, or add card"
//	@Router			/cards [post]
func (cc *Controller) AddCard(ctx *gin.Context) {
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

	response, err := cc.cardClient.AddCard(ctx, &card)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ReadAllCardsToLearn godoc
//
//	@Summary		Get all cards to learn
//	@Description	Retrieves all cards assigned to the user for learning
//	@Tags			cards
//	@Produce		json
//	@Success		200	{array}		model.Card
//	@Failure		500	{object}	model.ErrorResponse	"Internal Server Error - Failed to get user ID or retrieve cards"
//	@Router			/cards/learn [get]
func (cc *Controller) ReadAllCardsToLearn(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, err := cc.cardClient.ReadAllCardsToLearn(ctx, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// ReadAllCards godoc
//
//	@Summary		Get all cards
//	@Description	Retrieves all cards for the authenticated user. Admins can access cards of other users by providing user_id query parameter
//	@Tags			cards
//	@Produce		json
//	@Param			user_id	query		string	false	"User ID (admin only)"
//	@Success		200		{array}		model.Card
//	@Failure		403		{object}	model.ErrorResponse	"Forbidden - User without admin rights cannot access other users' cards"
//	@Failure		400		{object}	model.ErrorResponse	"Bad Request - Invalid user_id format"
//	@Failure		500		{object}	model.ErrorResponse	"Internal Server Error - Failed to get user ID or retrieve cards"
//	@Router			/cards [get]
func (cc *Controller) ReadAllCards(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		cc.log.Debug("Error getting user ID from context", "error", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cc.log.Debug("Current User ID", "userId", userId)

	var uuidUserId uuid.UUID
	queryUserId := ctx.Query("user_id")
	cc.log.Debug("Query User ID", "queryUserId", queryUserId)

	if queryUserId != "" {
		if uuidUserId, err = uuid.Parse(queryUserId); err != nil {
			cc.log.Debug("Error parsing query user ID", "error", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		cc.log.Debug("Parsed Query User ID", "uuidUserId", uuidUserId)
	}

	isAdmin, err := GetIsAdminFromContext(ctx)
	if err != nil {
		cc.log.Debug("Error getting admin status (defaulting to false)", "error", err)
		isAdmin = false // Default to false if admin check fails
	}
	cc.log.Debug("Is Admin", "isAdmin", isAdmin)

	var targetUserId uuid.UUID

	if queryUserId != "" {
		cc.log.Debug("Query user specified, checking permissions...")
		if !isAdmin && uuidUserId != userId {
			cc.log.Debug("Access denied: non-admin user trying to access other user's cards")
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User without admin rights cannot access other users' cards"})
			return
		}
		targetUserId = uuidUserId
		cc.log.Debug("Using query user ID as target", "targetUserId", targetUserId)
	} else {
		targetUserId = userId
		cc.log.Debug("Using current user ID as target", "targetUserId", targetUserId)
	}

	cc.log.Debug("Final Target User ID", "targetUserId", targetUserId)

	response, err := cc.cardClient.ReadAllCards(ctx, targetUserId)
	if err != nil {
		cc.log.Debug("Error calling card client", "error", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	cc.log.Debug("Cards response", "length", len(response))
	if len(response) > 0 {
		cc.log.Debug("First card", "card", response[0])
	}

	ctx.JSON(http.StatusOK, response)
}

// UpdateCard godoc
//
//	@Summary		Update a card
//	@Description	Update a card's content by ID
//	@Tags			cards
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"Card ID"
//	@Param			card	body		schemes.UpdateCardScheme	true	"Updated card data"
//	@Success		200		{object}	model.Card
//	@Failure		400		{object}	map[string]string
//	@Router			/card/{id} [put]
func (cc *Controller) UpdateCard(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	uid, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var cardUpdate schemes.UpdateCardScheme
	if err = ctx.ShouldBindJSON(&cardUpdate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	card, err := cc.cardClient.UpdateCard(ctx, uid, cardId, &cardUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, card)
}

// DeleteCard godoc
//
//	@Summary		Delete a card
//	@Description	Delete a card by ID
//	@Tags			cards
//	@Param			id	path	string	true	"Card ID"
//	@Success		200
//	@Failure		400	{object}	map[string]string
//	@Router			/card/{id} [delete]
func (cc *Controller) DeleteCard(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = cc.cardClient.DeleteCard(ctx, cardId, userId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}

// AddAnswers godoc
//
//	@Summary		Submit answers
//	@Description	Submit answers to cards
//	@Tags			answers
//	@Accept			json
//	@Produce		json
//	@Param			answers	body		[]schemes.AnswerScheme	true	"List of answers"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/answers [post]
func (cc *Controller) AddAnswers(ctx *gin.Context) {
	var answers []*schemes.AnswerScheme

	if err := ctx.ShouldBindJSON(&answers); err != nil {
		fmt.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("ANSSASAS", answers[0])

	if _, err = cc.cardClient.AddAnswers(ctx, userId, answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "added answers succesfully "})
}
