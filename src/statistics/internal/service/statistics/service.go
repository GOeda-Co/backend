package statistics

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"repeatro/src/statistics/pkg/model"
)

func getMean(grades []int) int {
	var sum int
	for _, grade := range grades {
		sum += grade
	}

	return sum / len(grades)
}

type statsRepository interface {
	/* mean grade some period / cards learned by date / Add results / Delete results*/
	Add(result *model.Stat) error
	Delete(resultId uuid.UUID) error
	// same as next
	GetAllGradesForPeriod(dtStart time.Time, dtEnd time.Time, userId uuid.UUID) ([]int, error)
	// Here i basically get all card for specific user over a period
	GetLearnedCardsForPeriod(dtStart time.Time, dtEnd time.Time, userId uuid.UUID) ([]uuid.UUID, error)
}

type Service struct {
statsRepo statsRepository
}

func NewService(statsRepo statsRepository) *Service {
	return &Service{
		statsRepo: statsRepo,
	}
}

func (rs *Service) GetMeanGradeOfPeriod(dtStart time.Time, dtEnd time.Time, userId uuid.UUID) (int, error) {
	grades, err := rs.statsRepo.GetAllGradesForPeriod(dtStart, dtEnd, userId)
	if err != nil {
		return 0, err
	}

	fmt.Println("graaades", grades)

	if len(grades) == 0 {
		return 0, fmt.Errorf("grades over this period are not found")
	}

	return getMean(grades), nil
}

//TODO: Encapsilate this logic	
/*
	So basically i need to get all ids as from microversive
	Get them intro repeatro
	And then in service layer using some ids to cards controller from cards
	translate into card objects

*/


// func (rs *Service) GetLearnedCardsForPeriod(dtStart time.Time, dtEnd time.Time, userId uuid.UUID) ([]*model.Stat, error) {
// 	cardIds, err := rs.statsRepo.GetLearnedCardsForPeriod(dtStart, dtEnd, userId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if len(cardIds) == 0 {
// 		return nil, fmt.Errorf("learned cards over this period are not found")
// 	}
// 	// cards := make([]*model.Stat.Card, 1)

// 	// for _, cardId := range cardIds {
// 	// 	card, err := rs.cardRepository.ReadCard(cardId)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// 	cards = append(cards, card)
// 	// }

// 	return cards, nil
// }
