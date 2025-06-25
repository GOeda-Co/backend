package services

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"repeatro/src/card/internal/repository/postgresql"
	"repeatro/src/card/pkg/model"
	"repeatro/src/card/pkg/scheme"
)

//TODO: Add results to repository

type ReviewResult struct {
	NextReviewTime time.Time
	Interval       int // in minutes
	Easiness       float64
	Repetitions    int
}

func SM2(
	now time.Time,
	previousInterval int, // in minutes
	previousEasiness float64,
	repetitions int,
	grade int, // 0–5 scale
) ReviewResult {
	var interval int
	easiness := previousEasiness
	if grade < 3 {
		repetitions = 0
		interval = 1 // reset to 5 minutes
	} else {
		switch repetitions {
		case 0:
			interval = 5
		case 1:
			interval = 30
		default:
			minutes := float64(previousInterval) * easiness
			interval = int(math.Round(minutes))
		}
		repetitions++
	}

	easiness += 0.1 - float64(5-grade)*(0.08+float64(5-grade)*0.02)
	if easiness < 1.3 {
		easiness = 1.3
	}

	fmt.Println(time.Duration(interval))

	nextReviewTime := now.Add(time.Duration(interval) * time.Minute)

	return ReviewResult{
		NextReviewTime: nextReviewTime,
		Interval:       interval,
		Easiness:       easiness,
		Repetitions:    repetitions,
	}
}

type CardService struct {
	cardRepository   postgresql.CardRepositoryInterface
}

func CreateNewService(cardRepository *postgresql.Repository) *CardService {
	return &CardService{
		cardRepository:   cardRepository,
	}
}

type CardServiceInterface interface {
	AddCard(card *model.Card) (*model.Card, error)
	ReadAllCards(userId uuid.UUID) ([]model.Card, error)
	UpdateCard(id uuid.UUID, card *schemes.UpdateCardScheme, userId uuid.UUID) (*model.Card, error)
	DeleteCard(id uuid.UUID, userId uuid.UUID) error
	AddAnswers(userId uuid.UUID, answers []schemes.AnswerScheme) error
}

func (cs CardService) AddCard(card *model.Card) (*model.Card, error) {
	err := cs.cardRepository.AddCard(card)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (cm CardService) ReadAllCards(userId uuid.UUID) ([]model.Card, error) {
	cards, err := cm.cardRepository.ReadAllCards(userId)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (cm CardService) UpdateCard(cardId uuid.UUID, cardUpdate *schemes.UpdateCardScheme, userId uuid.UUID) (*model.Card, error) {
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

func (cm CardService) DeleteCard(cardId uuid.UUID, userId uuid.UUID) error {
	cardFound, err := cm.cardRepository.ReadCard(cardId)
	if err != nil {
		return err
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

func (cm CardService) AddAnswers(userId uuid.UUID, answers []schemes.AnswerScheme) error {
	for _, answer := range answers {
		if answer.Grade < 0 || answer.Grade > 5 {
			return fmt.Errorf("invalid grade")
		}

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
			return fmt.Errorf("invalid card owner")
		}
		
		// recalculate values
		reviewResult := SM2(time.Now(),
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
		
		// result := model.Result {
			// 	UserId: card.CreatedBy,
			// 	CardId: card.CardId,
			// 	Grade: answer.Grade,
			// }
			
			// fmt.Println("Res", result.CreatedAt)
			
			// if err = cm.resultRepository.AddResult(&result); err != nil {
				// 	fmt.Println("so here")
				// 	return err
				// }
	}
	return nil
}

