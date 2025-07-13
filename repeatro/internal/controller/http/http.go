package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/swaggo/swag/example/celler/httputil"
	"github.com/google/uuid"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	deckClient "github.com/tomatoCoderq/repeatro/internal/clients/deck/grpc"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	model "github.com/tomatoCoderq/repeatro/pkg/models"
	"github.com/tomatoCoderq/repeatro/pkg/schemes"
)

type Controller struct {
	ssoClient  *ssoClient.Client
	cardClient *cardClient.Client
	deckClient *deckClient.Client
}

func New(ssoClient *ssoClient.Client, cardClient *cardClient.Client, deckClient *deckClient.Client) *Controller {
	return &Controller{
		ssoClient:  ssoClient,
		cardClient: cardClient,
		deckClient: deckClient,
	}
}

// Register godoc
//	@Summary		Registers new user to the system
//	@Description	Register by email, name, and password, getting user_id
//	@Tags			sso
//	@Accept			json
//	@Produce		json
//	@Param			name		body		string	true	"Name of user"
//	@Param			email		body		string	true	"Email of user"
//	@Param			password	body		string	true	"Password of user (> 6 letters)"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/register [post]
func (c *Controller) Register(ctx *gin.Context) {
	var registerScheme schemes.RegisterScheme

	if err := ctx.ShouldBindBodyWithJSON(&registerScheme); err != nil {
		fmt.Println("HERE")
		ctx.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request body: %v", err)})
	}

	fmt.Println(registerScheme)

	uid, err := c.ssoClient.Register(ctx.Request.Context(), registerScheme.Email, registerScheme.Password, registerScheme.Name)
	if err != nil {
		fmt.Println("Err", err)
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Failed to register user: %v", err)})
		return
	}
	ctx.JSON(200, gin.H{
		"user_id": uid,
		"message": "User registered successfully",
	})
}

// Login godoc
//	@Summary		Logs in a user
//	@Description	Logs in a user and returns a JWT token
//	@Tags			sso
//	@Accept			json
//	@Produce		json
//	@Param			email		body		string	true	"Email of user"
//	@Param			password	body		string	true	"Password of user"
//	@Param			app_id		body		int		true	"Application ID"
//	@Success		200			{object}	map[string]interface{}
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/login [post]
func (c *Controller) Login(ctx *gin.Context) {
	var loginScheme schemes.LoginScheme

	if err := ctx.ShouldBindBodyWithJSON(&loginScheme); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
	}

	token, err := c.ssoClient.Login(ctx.Request.Context(), loginScheme.Email, loginScheme.Password, loginScheme.AppId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("Failed to login user: %v", err)})
		return
	}
	ctx.JSON(200, gin.H{
		"token":   token,
		"message": "User registered successfully",
	})
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

// AddCard godoc
//	@Summary		Add a card
//	@Description	Add a new card for the authenticated user
//	@Tags			cards
//	@Accept			json
//	@Produce		json
//	@Param			card	body		model.Card	true	"Card to add"
//	@Success		200		{object}	model.Card
//	@Failure		500		{object}	map[string]string
//	@Router			/card [post]
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
//	@Summary		Get all cards to learn
//	@Description	Retrieves all cards assigned to the user for learning
//	@Tags			cards
//	@Produce		json
//	@Success		200	{array}		model.Card
//	@Failure		500	{object}	map[string]string
//	@Router			/cards [get]
func (cc *Controller) ReadAllCardsToLearn(ctx *gin.Context) {
	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		fmt.Println("USS", userId)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, err := cc.cardClient.ReadAllCards(ctx, userId)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// UpdateCard godoc
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
	fmt.Println("COMEHERE")

	card, err := cc.cardClient.UpdateCard(ctx, uid, cardId, &cardUpdate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("CO", card)
	ctx.JSON(http.StatusOK, card)
}

// DeleteCard godoc
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := GetUserIdFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err = cc.cardClient.AddAnswers(ctx, userId, answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "added answers succesfully "})
}

// AddDeck godoc
//	@Summary		Add a deck
//	@Description	Create a new deck
//	@Tags			decks
//	@Accept			json
//	@Produce		json
//	@Param			deck	body		model.Deck	true	"Deck to add"
//	@Success		200		{object}	model.Deck
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
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
//	@Summary		Get all decks
//	@Description	Retrieves all decks in the system
//	@Tags			decks
//	@Produce		json
//	@Success		200	{array}		model.Deck
//	@Failure		500	{object}	map[string]string
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
//	@Summary		Get deck by ID
//	@Description	Retrieves a deck by its ID
//	@Tags			decks
//	@Param			id	path		string	true	"Deck ID"
//	@Success		200	{object}	model.Deck
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
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

// DeleteDeck godoc
//	@Summary		Delete a deck
//	@Description	Delete a deck by ID
//	@Tags			decks
//	@Param			id	path	string	true	"Deck ID"
//	@Success		200
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
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
//	@Summary		Add card to deck
//	@Description	Add a card to a specific deck
//	@Tags			decks
//	@Param			card_id	path	string	true	"Card ID"
//	@Param			deck_id	path	string	true	"Deck ID"
//	@Success		200
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
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

	// uid, err := GetUserIdFromContext(ctx)
	// if err != nil {
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	err = cc.deckClient.AddCardToDeck(ctx, did, cid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add card to deck"})
		return
	}
	ctx.Status(http.StatusOK)
}

// ReadCardsFromDeck godoc
//	@Summary		Get cards from deck
//	@Description	Retrieve all cards from a specific deck
//	@Tags			decks
//	@Param			id	path		string	true	"Deck ID"
//	@Success		200	{array}		model.Card
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/deck/{id}/cards [get]
func (cc *Controller) ReadCardsFromDeck(ctx *gin.Context) {
	did, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid deck ID"})
		return
	}

	response, err := cc.deckClient.ReadCardsFromDeck(ctx, did)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cards from deck"})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// add card
// add deck
// add many cards in one deck
// delete deck (and all cards inside)
// delete card (and from the deck it's assigned to)
// update card info
// update deck info

/*in the future:
admin can check all possible deecks
deck can be open/closed
pictures taken from other repo
*/
