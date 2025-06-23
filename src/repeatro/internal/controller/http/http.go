package http

import (
	"encoding/json"
	"errors"
	// "fmt"
	"log"
	"repeatro/src/repeatro/internal/gateway"
	"repeatro/src/repeatro/internal/service"
	// schemes "repeatro/src/user/pkg/scheme"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"/
)

type Controller struct {
	service *service.Service
}

// New creates a new movie HTTP handler.
func New(service *service.Service) *Controller {
	return &Controller{service}
}

//TODO: Rewrite with gin

// GetMovieDetails handles GET /movie requests.
func (c *Controller) GetAllCards(ctx *gin.Context) {
	details, err := c.service.GetCards(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) AddCard(ctx *gin.Context) {
	details, err := c.service.AddCard(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) UpdateCard(ctx *gin.Context) {
	details, err := c.service.UpdateCard(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) DeleteCard(ctx *gin.Context) {
	details, err := c.service.DeleteCard(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) AddAnswers(ctx *gin.Context) {
	details, err := c.service.AddAnswers(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(details); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}


func (c *Controller) Login(ctx *gin.Context) {
	token, err := c.service.Login(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(token); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) Register(ctx *gin.Context) {
	token, err := c.service.Register(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	if err := json.NewEncoder(ctx.Writer).Encode(token); err != nil {
		log.Printf("Encode error: %v\n", err)
	}
}

func (c *Controller) AddDeck(ctx *gin.Context) {
	deck, err := c.service.AddDeck(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("AddDeck error: %v\n", err)
		ctx.AbortWithStatusJSON(200, gin.H{"message": err.Error()})
		return
	}
	_ = json.NewEncoder(ctx.Writer).Encode(deck)
}

func (c *Controller) ReadAllDecks(ctx *gin.Context) {
	decks, err := c.service.ReadAllDecksOfUser(ctx.Request.Context(), ctx)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("ReadAllDecks error: %v\n", err)
		ctx.AbortWithStatusJSON(500, gin.H{"message": err.Error()})
		return
	}
	_ = json.NewEncoder(ctx.Writer).Encode(decks)
}

func (c *Controller) ReadDeck(ctx *gin.Context) {
	deckId := ctx.Param("id")
	if deckId == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "missing deck ID"})
		return
	}

	deck, err := c.service.ReadDeck(ctx.Request.Context(), ctx, deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("ReadDeck error: %v\n", err)
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	}
	_ = json.NewEncoder(ctx.Writer).Encode(deck)
}

func (c *Controller) ReadCardsFromDeck(ctx *gin.Context) {
	deckId := ctx.Param("id")
	if deckId == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "missing deck ID"})
		return
	}

	cards, err := c.service.ReadAllCardsFromDeck(ctx.Request.Context(), ctx, deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("ReadCardsFromDeck error: %v\n", err)
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	}
	_ = json.NewEncoder(ctx.Writer).Encode(cards)
}

func (c *Controller) DeleteDeck(ctx *gin.Context) {
	deckId := ctx.Param("id")
	if deckId == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "missing deck ID"})
		return
	}

	err := c.service.DeleteDeck(ctx.Request.Context(), ctx, deckId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("DeleteDeck error: %v\n", err)
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "deck deleted"})
}

func (c *Controller) AddCardToDeck(ctx *gin.Context) {
	deckId := ctx.Param("deck_id")
	cardId := ctx.Param("card_id")

	if deckId == "" || cardId == "" {
		ctx.AbortWithStatusJSON(400, gin.H{"message": "missing deck_id or card_id"})
		return
	}

	err := c.service.AddCardToDeck(ctx.Request.Context(), ctx, deckId, cardId)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	} else if err != nil {
		log.Printf("AddCardToDeck error: %v\n", err)
		ctx.AbortWithStatusJSON(403, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"message": "card added to deck"})
}