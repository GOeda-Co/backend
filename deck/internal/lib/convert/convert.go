package convert

import (
	"github.com/google/uuid"
	"github.com/tomatoCoderq/deck/pkg/model"
	deckv1 "github.com/GOeda-Co/proto-contract/gen/go/deck"
	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoToModelDeck converts a protobuf AddDeckRequest to a model.Deck
func ProtoToModelDeck(in *deckv1.AddDeckRequest, userId string) *model.Deck {
	uid, _ := uuid.Parse(userId)

	return &model.Deck{
		Name:        in.Name,
		Description: in.Description,
		CreatedBy:   uid,
	}
}

// ModelToProtoDeck converts a model.Deck to a protobuf Deck
func ModelToProtoDeck(deck *model.Deck) *deckv1.Deck {
	if deck == nil {
		return nil
	}

	return &deckv1.Deck{
		DeckId:      deck.DeckId.String(),
		CreatedBy:   deck.CreatedBy.String(),
		CreatedAt:   timestamppb.New(deck.CreatedAt),
		Name:        deck.Name,
		CardsQuantity: uint32(deck.CardsQuantity),
		Description: deck.Description,
	}
}


func ProtoToModel(card *cardv1.Card) (*model.Card, error) {
	cardId, err := uuid.Parse(card.CardId)
	if err != nil {
		return nil, err
	}
	createdBy, err := uuid.Parse(card.CreatedBy)
	if err != nil {
		return nil, err
	}
	deckId, err := uuid.Parse(card.DeckId)
	if err != nil {
		return nil, err
	}

	return &model.Card{
		CardId:           cardId,
		CreatedBy:        createdBy,
		CreatedAt:        card.CreatedAt.AsTime(),
		Word:             card.Word,
		Translation:      card.Translation,
		Easiness:         card.Easiness,
		UpdatedAt:        card.UpdatedAt.AsTime(),
		Interval:         int(card.Interval),
		ExpiresAt:        card.ExpiresAt.AsTime(),
		RepetitionNumber: int(card.RepetitionNumber),
		DeckID:           deckId,
		Tags:             card.Tags,
	}, nil
}


func ModelToProto(card *model.Card) *cardv1.Card {
	return &cardv1.Card{
		CardId:           card.CardId.String(),
		CreatedBy:        card.CreatedBy.String(),
		CreatedAt:        timestamppb.New(card.CreatedAt),
		Word:             card.Word,
		Translation:      card.Translation,
		Easiness:         card.Easiness,
		UpdatedAt:        timestamppb.New(card.UpdatedAt),
		Interval:         int32(card.Interval),
		ExpiresAt:        timestamppb.New(card.ExpiresAt),
		RepetitionNumber: int32(card.RepetitionNumber),
		DeckId:           card.DeckID.String(),
		Tags:             card.Tags,
	}
}

func ModelToProtoCard(card *model.Card) *deckv1.Card {
	return &deckv1.Card{
		CardId:           card.CardId.String(),
		CreatedBy:        card.CreatedBy.String(),
		CreatedAt:        timestamppb.New(card.CreatedAt),
		Word:             card.Word,
		Translation:      card.Translation,
		Easiness:         card.Easiness,
		UpdatedAt:        timestamppb.New(card.UpdatedAt),
		Interval:         int32(card.Interval),
		ExpiresAt:        timestamppb.New(card.ExpiresAt),
		RepetitionNumber: int32(card.RepetitionNumber),
		DeckId:           card.DeckID.String(),
		Tags:             card.Tags,
	}
}

// func ProtoToUpdateCardScheme(card *cardv1.UpdateCardRequest) *schemes.UpdateCardScheme{
// 	return &schemes.UpdateCardScheme{
// 		Word:             card.Word,
// 		Translation:      card.Translation,
// 		Easiness:         card.Easiness,
// 		Interval:         int(card.Interval),
// 		ExpiresAt:        card.ExpiresAt.AsTime(),
// 		RepetitionNumber: int(card.RepetitionNumber),
// 		Tags:             card.Tags,
// 	}
// }

// func ProtoToAnswerSchemes(answer *cardv1.Answer) (*schemes.AnswerScheme, error) {
// 	cardId, err := uuid.Parse(answer.CardId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &schemes.AnswerScheme{
// 		CardId: cardId,
// 		Grade:  int(answer.Grade),
// 	}, nil
// }

// func AnswerToProtoScheme(answer *cardv1.Answer) (*schemes.AnswerScheme, error) {
// 	cardId, err := uuid.Parse(answer.CardId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &schemes.AnswerScheme{
// 		CardId: cardId,
// 		Grade:  int(answer.Grade),
// 	}, nil
// }
