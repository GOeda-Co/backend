package grpc

import (
	"context"
	"errors"
	"fmt"

	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	"github.com/google/uuid"
	statClient "github.com/tomatoCoderq/card/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/card/internal/controller"
	"github.com/tomatoCoderq/card/internal/lib/convert"
	"github.com/tomatoCoderq/card/internal/lib/security"
	schemes "github.com/tomatoCoderq/card/pkg/scheme"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetAuthUser(ctx context.Context) (*security.AuthUser, error) {
	val := ctx.Value(security.UserContextKey)
	if val == nil {
		return nil, errors.New("no user in context")
	}

	authUser, ok := val.(security.AuthUser)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}

	return &authUser, nil
}

type ServerAPI struct {
	cardv1.UnimplementedCardServiceServer
	service    controller.Card
	statClient *statClient.Client
}

func Register(gRPCServer *grpc.Server, card controller.Card, statClient *statClient.Client) {
	cardv1.RegisterCardServiceServer(
		gRPCServer,
		&ServerAPI{
			service:    card,
			statClient: statClient,
		})
}

func (s *ServerAPI) AddCard(ctx context.Context, in *cardv1.AddCardRequest) (*cardv1.AddCardResponse, error) {
	if in.Card.Word == "" {
		return nil, status.Error(codes.InvalidArgument, "Word is required")
	}
	if in.Card.Translation == "" {
		return nil, status.Error(codes.InvalidArgument, "Translation is required")
	}

	card, err := convert.ProtoToModel(in.Card)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed during converting proto card model to inner card model")
	}

	fullCard, err := s.service.AddCard(card)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed during adding card")
	}

	return &cardv1.AddCardResponse{Card: convert.ModelToProto(fullCard)}, nil
}

func (s *ServerAPI) ReadAllCards(ctx context.Context, in *cardv1.ReadAllCardsRequest) (*cardv1.ReadAllCardsResponse, error) {
	_, err := uuid.Parse(in.UserId)

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to auth user: %v", err))
	}
	// ctx = withToken(ctx, ctx.Value("token").(string))

	// if err != nil {
	// 	return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	// }

	cards, err := s.service.ReadAllCards(authUser.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to read cards")
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.ModelToProto(&card))
	}

	return &cardv1.ReadAllCardsResponse{Cards: protoCards}, nil
}

func (s *ServerAPI) ReadAllCardsByUser(ctx context.Context, in *cardv1.ReadAllCardsByUserRequest) (*cardv1.ReadAllCardsByUserResponse, error) {
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	cards, err := s.service.ReadAllCardsByUser(userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to read cards")
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.ModelToProto(&card))
	}

	return &cardv1.ReadAllCardsByUserResponse{Cards: protoCards}, nil
}

func (s *ServerAPI) UpdateCard(ctx context.Context, in *cardv1.UpdateCardRequest) (*cardv1.UpdateCardResponse, error) {
	cardId, err := uuid.Parse(in.CardId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid card ID")
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	cardUpdate := convert.ProtoToUpdateCardScheme(in)
	if cardUpdate == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid update payload")
	}

	updatedCard, err := s.service.UpdateCard(cardId, cardUpdate, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update card")
	}

	return &cardv1.UpdateCardResponse{Card: convert.ModelToProto(updatedCard)}, nil
}

func (s *ServerAPI) DeleteCard(ctx context.Context, in *cardv1.DeleteCardRequest) (*cardv1.DeleteCardResponse, error) {
	cardId, err := uuid.Parse(in.CardId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid card ID")
	}
	userId, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	err = s.service.DeleteCard(cardId, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete card: %v", err))
	}

	return &cardv1.DeleteCardResponse{}, nil
}

func (s *ServerAPI) AddAnswers(ctx context.Context, in *cardv1.AddAnswersRequest) (*cardv1.AddAnswersResponse, error) {
	var answers []schemes.AnswerScheme

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to auth user: %v", err))
	}

	for _, answer := range in.Answers {
		if answer.CardId == "" {
			return nil, status.Error(codes.InvalidArgument, "Card ID is required in answers")
		}
		if answer.Grade < 0 || answer.Grade > 5 {
			return nil, status.Error(codes.InvalidArgument, "Answer is required in answers")
		}
		answerConverted, err := convert.ProtoToAnswerSchemes(answer)
		if err != nil {
			return nil, status.Error(codes.Internal, "Failed during converting answer")
		}
		fmt.Printf("CONVER: %v", answerConverted)
		answers = append(answers, *answerConverted)
	}

	fmt.Println("ANS", answers)

	err = s.service.AddAnswers(ctx, authUser.ID, answers)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to add answers: %v", err))
	}

	return &cardv1.AddAnswersResponse{Message: "added answers successfully"}, nil
}
