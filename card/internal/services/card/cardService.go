package services

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	statClient "github.com/tomatoCoderq/card/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/card/internal/lib/sm2"
	"github.com/tomatoCoderq/card/pkg/model"
	"github.com/tomatoCoderq/card/pkg/scheme"
)

type CardRepository interface {
	AddCard(card *model.Card) error
	ReadAllCards(userId uuid.UUID) ([]model.Card, error)
	ReadAllCardsByUser(userId uuid.UUID) ([]model.Card, error)
	ReadCard(cardId uuid.UUID) (*model.Card, error)
	PureUpdate(card *model.Card) error
	UpdateCard(card *model.Card, cardUpdate *schemes.UpdateCardScheme) (*model.Card, error)
	DeleteCard(cardId uuid.UUID) error
}

type Card struct {
	log            *slog.Logger
	cardRepository CardRepository
	statClient *statClient.Client
}

func New(
	log *slog.Logger,
	cardRepo CardRepository,
	statClient *statClient.Client,
) *Card {
	return &Card{
		log:            log,
		cardRepository: cardRepo,
		statClient: statClient,
	}
}

func (cs Card) AddCard(card *model.Card) (*model.Card, error) {
	err := cs.cardRepository.AddCard(card)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (cm Card) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	cards, err := cm.cardRepository.ReadAllCards(userId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cm Card) ReadAllCardsByUser(userId uuid.UUID) ([]model.Card, error) {
	cards, err := cm.cardRepository.ReadAllCardsByUser(userId)
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

		fmt.Println("ANSID", answer.CardId)

		card, err := cm.cardRepository.ReadCard(answer.CardId)
		if err != nil {
			return err
		}

		fmt.Println("CARD", card)
		// NOTE: If expire_time not reached yet the card will be just skipped
		if time.Now().Compare(card.ExpiresAt) == -1 {
			fmt.Println("SKIPPED", time.Now(), card.ExpiresAt)
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
			fmt.Println("WTF")
			return err
		}

		fmt.Println("NEWCARD", card)

		// result := model.Review {
		// 	UserID: card.CreatedBy,
		// 	CardID: card.CardId,
		// 	Grade: int32(answer.Grade),
		// }

		// fmt.Println("Res", result.CreatedAt)
		fmt.Println("CAME TO RES")
		md, _ := metadata.FromIncomingContext(ctx)
		fmt.Println(md["authorization"])

		fmt.Println("INFODEK", card.DeckID.String())

		reviewId, err := cm.statClient.AddRecord(ctx, card.DeckID.String(), card.CardId.String(), answer.Grade)
		if err != nil {
			fmt.Println("SO HERE", reviewId)
			return err
		}
	}
	return nil
}
