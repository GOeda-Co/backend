package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/GOeda-Co/proto-contract/convert"
	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	schemes "github.com/GOeda-Co/proto-contract/scheme/card"
	"github.com/google/uuid"
	statClient "github.com/tomatoCoderq/card/internal/clients/stats/grpc"
	"github.com/tomatoCoderq/card/internal/controller"
	"github.com/tomatoCoderq/card/internal/lib/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

	card, err := convert.FromProtoToModelCard(in.Card)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed during converting proto card model to inner card model")
	}

	fullCard, err := s.service.AddCard(card)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed during adding card")
	}

	return &cardv1.AddCardResponse{Card: convert.FromModelToProtoCard(fullCard)}, nil
}

func (s *ServerAPI) ReadAllCards(ctx context.Context, in *cardv1.ReadAllCardsRequest) (*cardv1.ReadAllCardsResponse, error) {

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to auth user: %v", err))
	}

	cards, err := s.service.ReadAllCards(authUser.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to read cards")
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.FromModelToProtoCard(&card))
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
		protoCards = append(protoCards, convert.FromModelToProtoCard(&card))
	}

	return &cardv1.ReadAllCardsByUserResponse{Cards: protoCards}, nil
}

func (s *ServerAPI) SearchAllPublicCards(ctx context.Context, in *emptypb.Empty) (*cardv1.SearchAllPublicCardsResponse, error) {
	cards, err := s.service.SearchAllPublicCards()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to search public cards")
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.FromModelToProtoCard(&card))
	}

	return &cardv1.SearchAllPublicCardsResponse{Cards: protoCards}, nil
}

func (s *ServerAPI) SearchUserPublicCards(ctx context.Context, in *cardv1.SearchUserPublicCardsRequest) (*cardv1.SearchUserPublicCardsResponse, error) {
	cards, err := s.service.SearchUserPublicCards(in.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to search user public cards")
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.FromModelToProtoCard(&card))
	}

	return &cardv1.SearchUserPublicCardsResponse{Cards: protoCards}, nil
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

	cardUpdate := convert.FromProtoToUpdateSchemeCard(in)
	if cardUpdate == nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid update payload")
	}

	updatedCard, err := s.service.UpdateCard(cardId, cardUpdate, userId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to update card")
	}

	return &cardv1.UpdateCardResponse{Card: convert.FromModelToProtoCard(updatedCard)}, nil
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
		answerConverted, err := convert.FromProtoToAnswerSchemeCard(answer)
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
