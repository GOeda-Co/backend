package controller

import (
	"github.com/google/uuid"
	"github.com/tomatoCoderq/card/pkg/model"
	schemes "github.com/tomatoCoderq/card/pkg/scheme"
)

type Card interface {
	AddCard(card *model.Card) (*model.Card, error)
	ReadAllCards(userId uuid.UUID) ([]model.Card, error)
	ReadAllCardsByUser(userId uuid.UUID) ([]model.Card, error)
	UpdateCard(id uuid.UUID, card *schemes.UpdateCardScheme, userId uuid.UUID) (*model.Card, error)
	DeleteCard(id uuid.UUID, userId uuid.UUID) error
	AddAnswers(userId uuid.UUID, answers []schemes.AnswerScheme) error
}

