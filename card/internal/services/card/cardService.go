package services

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/tomatoCoderq/card/internal/lib/sm2"

	"github.com/GOeda-Co/proto-contract/model/card"
	schemes "github.com/GOeda-Co/proto-contract/scheme/card"
)

type CardRepository interface {
	AddCard(card *model.Card) error
	ReadAllOwnCardsToLearn(userId uuid.UUID) ([]model.Card, error)
	ReadAllOwnCards(userId uuid.UUID) ([]model.Card, error)
	SearchAllPublicCards() ([]model.Card, error)
	SearchUserPublicCards(userId uuid.UUID) ([]model.Card, error)
	ReadCard(cardId uuid.UUID) (*model.Card, error)
	PureUpdate(card *model.Card) error
	UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error)
	DeleteCard(cardId uuid.UUID) error
}

type StatsClient interface {
	AddRecord(ctx context.Context, deckId, cardId string, grade int) (string, error)
}

type Card struct {
	log            *slog.Logger
	cardRepository CardRepository
	statClient     StatsClient
}

func New(
	log *slog.Logger,
	cardRepo CardRepository,
	statClient StatsClient,
) *Card {
	return &Card{
		log:            log,
		cardRepository: cardRepo,
		statClient:     statClient,
	}
}

func (cs Card) AddCard(card *model.Card) (*model.Card, error) {
	err := cs.cardRepository.AddCard(card)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (cm Card) ReadAllOwnCardsToLearn(userId uuid.UUID) ([]model.Card, error) {
	cards, err := cm.cardRepository.ReadAllOwnCardsToLearn(userId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cm Card) ReadAllOwnCards(userId uuid.UUID) ([]model.Card, error) {
	cards, err := cm.cardRepository.ReadAllOwnCards(userId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cm Card) UpdateCard(cardId uuid.UUID, cardUpdate *schemes.UpdateCardScheme, userId uuid.UUID) (*model.Card, error) {
	cardFound, err := cm.cardRepository.ReadCard(cardId)
	if err != nil {
		return nil, err
	}

	if cardFound.CreatedBy != userId {
		return nil, fmt.Errorf("cannot update other's user card")
	}

	cardUpdated, err := cm.cardRepository.UpdateCard(cardFound, cardUpdate)
	if err != nil {
		return nil, err
	}

	return cardUpdated, nil
}

func (cm Card) DeleteCard(cardId uuid.UUID, userId uuid.UUID) error {
	cardFound, err := cm.cardRepository.ReadCard(cardId)
	if err != nil {
		return err
	}

	if cardFound == nil || reflect.DeepEqual(cardFound, &model.Card{}) {
		return fmt.Errorf("card not found")
	}

	if cardFound.CreatedBy != userId {
		return fmt.Errorf("cannot delete other's user card")
	}

	err = cm.cardRepository.DeleteCard(cardId)
	if err != nil {

		return nil
	}

	return nil
}

func (cm Card) AddAnswers(ctx context.Context, userId uuid.UUID, answers []schemes.AnswerScheme) error {

	for _, answer := range answers {
		if answer.Grade < 0 || answer.Grade > 5 {
			return fmt.Errorf("invalid grade")
		}

		card, err := cm.cardRepository.ReadCard(answer.CardId)
		if err != nil {
			return err
		}

		// NOTE: If expire_time not reached yet the card will be just skipped
		if time.Now().Compare(card.ExpiresAt) == -1 {
			cm.log.Info("Card not expired yet, skipping", "cardId", card.CardId, "expiresAt", card.ExpiresAt)
			continue
		}

		cardOwnerId := card.CreatedBy
		if userId != cardOwnerId {
			return fmt.Errorf("invalid card owner. got %v. want %v", cardOwnerId, card.CreatedBy)
		}

		// recalculate values
		reviewResult := sm2.SM2(time.Now(),
			card.Interval,
			card.Easiness,
			card.RepetitionNumber,
			answer.Grade)

		// write back to db
		card.UpdatedAt = time.Now()
		card.ExpiresAt = reviewResult.NextReviewTime
		card.Easiness = reviewResult.Easiness
		card.Interval = int(reviewResult.Interval)
		card.RepetitionNumber = reviewResult.Repetitions

		if err = cm.cardRepository.PureUpdate(card); err != nil {
			return err
		}

		md, _ := metadata.FromIncomingContext(ctx)
		cm.log.Info("Authorization Metadata", "authorization", md["authorization"])

		reviewId, err := cm.statClient.AddRecord(ctx, card.DeckID.String(), card.CardId.String(), answer.Grade)
		if err != nil {
			cm.log.Error("Failed to add stat record", "error", err, "reviewId", reviewId)
			return err
		}
	}
	return nil
}

func (cm Card) SearchAllPublicCards() ([]model.Card, error) {
	cards, err := cm.cardRepository.SearchAllPublicCards()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cm Card) SearchUserPublicCards(userId string) ([]model.Card, error) {
	userIdParsed, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	if userIdParsed == uuid.Nil {
		return nil, fmt.Errorf("user ID cannot be empty")
	}

	cards, err := cm.cardRepository.SearchUserPublicCards(userIdParsed)
	if err != nil {
		return nil, err
	}
	return cards, nil
}
