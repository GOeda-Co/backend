package controller

import (
	"context"

	"github.com/google/uuid"
	// "github.com/tomatoCoderq/card/pkg/model"
	"github.com/GOeda-Co/proto-contract/model/card"
	schemes "github.com/GOeda-Co/proto-contract/scheme/card"
)

type Card interface {
	AddCard(card *model.Card) (*model.Card, error)
	ReadAllOwnCardsToLearn(userId uuid.UUID) ([]model.Card, error)
	ReadAllOwnCards(userId uuid.UUID) ([]model.Card, error)
	SearchAllPublicCards() ([]model.Card, error)
	SearchUserPublicCards(useId string) ([]model.Card, error)
	UpdateCard(id uuid.UUID, card *schemes.UpdateCardScheme, userId uuid.UUID) (*model.Card, error)
	DeleteCard(id uuid.UUID, userId uuid.UUID) error
	AddAnswers(ctx context.Context, userId uuid.UUID, answers []schemes.AnswerScheme) error
}
