package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	ssoClient "github.com/tomatoCoderq/repeatro/internal/clients/sso/grpc"
	cardClient "github.com/tomatoCoderq/repeatro/internal/clients/card/grpc"
	model "github.com/tomatoCoderq/repeatro/pkg/models"
	"github.com/tomatoCoderq/repeatro/pkg/schemes"
)

type Controller struct {
	ssoClient *ssoClient.Client
	cardClient *cardClient.Client
}

func New(ssoClient *ssoClient.Client, cardClient *cardClient.Client) *Controller {
	return &Controller{
		ssoClient: ssoClient,
		cardClient: cardClient,
	}
}

func (c *Controller) Register(ctx *gin.Context) {
	var registerScheme schemes.RegisterScheme

	if err := ctx.ShouldBindBodyWithJSON(&registerScheme); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
	}

	uid, err := c.ssoClient.Register(ctx.Request.Context(), registerScheme.Email, registerScheme.Password)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to register user"})
		return
	}
	ctx.JSON(200, gin.H{
		"user_id": uid,
		"message": "User registered successfully",
	})
}

func (c *Controller) Login(ctx *gin.Context) {
	var loginScheme schemes.LoginScheme

	if err := ctx.ShouldBindBodyWithJSON(&loginScheme); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
	}

	token, err := c.ssoClient.Login(ctx.Request.Context(), loginScheme.Email, loginScheme.Password, loginScheme.AppId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to login user"})
		return
	}
	ctx.JSON(200, gin.H{
		"token": token,
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

func (cc *Controller) UpdateCard(ctx *gin.Context) {
	cardId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	
	// userId, err := GetUserIdFromHeader(ctx)
	// if err != nil {
	// 	ctx.AbortWithError(http.StatusInternalServerError, err)
	// 	return
	// }
	
	var cardUpdate schemes.UpdateCardScheme
	if err = ctx.ShouldBindJSON(&cardUpdate); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	fmt.Println("COMEHERE")


	card, err := cc.cardClient.UpdateCard(ctx, cardId, &cardUpdate)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	fmt.Println("CO", card)
	ctx.JSON(http.StatusOK, card)
}

func (cc *Controller) DeleteCard(ctx *gin.Context) {
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

	success, err := cc.cardClient.DeleteCard(ctx, cardId, userId)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !success {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Card"})
		return 
	}
	ctx.Status(http.StatusOK)
}

func (cc *Controller) AddAnswers(ctx *gin.Context) {
	var answers []*schemes.AnswerScheme

	if err := ctx.ShouldBindJSON(&answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	userId, err := GetUserIdFromHeader(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err = cc.cardClient.AddAnswers(ctx, userId, answers); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "added answers succesfully "})
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