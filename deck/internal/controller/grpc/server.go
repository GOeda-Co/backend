package grpc

import (
	"context"
	"errors"
	"fmt"

	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	deckv1 "github.com/GOeda-Co/proto-contract/gen/go/deck"
	"github.com/google/uuid"
	"github.com/tomatoCoderq/deck/internal/controller"
	"github.com/tomatoCoderq/deck/internal/lib/convert"
	"github.com/tomatoCoderq/deck/internal/lib/security"
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

type DeckServerAPI struct {
	deckv1.UnimplementedDeckServiceServer
	service controller.Deck
}

func Register(gRPCServer *grpc.Server, deck controller.Deck) {
	deckv1.RegisterDeckServiceServer(gRPCServer, &DeckServerAPI{service: deck})
}

func (s *DeckServerAPI) AddDeck(ctx context.Context, in *deckv1.AddDeckRequest) (*deckv1.DeckResponse, error) {
	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}

	authUser, err := GetAuthUser(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}

	deck := convert.ProtoToModelDeck(in, authUser.ID.String())
	createdDeck, err := s.service.AddDeck(deck)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to add deck")
	}

	return &deckv1.DeckResponse{Deck: convert.ModelToProtoDeck(createdDeck)}, nil
}

func (s *DeckServerAPI) ReadAllDecks(ctx context.Context, in *emptypb.Empty) (*deckv1.DeckListResponse, error) {
	authUser, err := GetAuthUser(ctx)
	fmt.Println("Auth User", authUser)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}

	decks, err := s.service.ReadAllDecksOfUser(authUser.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to fetch decks")
	}

	var protoDecks []*deckv1.Deck
	for _, deck := range decks {
		protoDecks = append(protoDecks, convert.ModelToProtoDeck(&deck))
	}

	return &deckv1.DeckListResponse{Decks: protoDecks}, nil
}

func (s *DeckServerAPI) ReadDeck(ctx context.Context, in *deckv1.ReadDeckRequest) (*deckv1.DeckResponse, error) {
	deckId, err := uuid.Parse(in.DeckId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid deck ID")
	}

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}

	deck, err := s.service.ReadDeck(deckId, authUser.ID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &deckv1.DeckResponse{Deck: convert.ModelToProtoDeck(deck)}, nil
}

func (s *DeckServerAPI) DeleteDeck(ctx context.Context, in *deckv1.ReadDeckRequest) (*emptypb.Empty, error) {
	deckId, err := uuid.Parse(in.DeckId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid deck ID")
	}

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}

	err = s.service.DeleteDeck(deckId, authUser.ID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *DeckServerAPI) AddCardToDeck(ctx context.Context, in *deckv1.AddCardToDeckRequest) (*emptypb.Empty, error) {
	cardId, err := uuid.Parse(in.CardId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid card ID")
	}
	deckId, err := uuid.Parse(in.DeckId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid deck ID")
	}

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}

	err = s.service.AddCardToDeck(cardId, deckId, authUser.ID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *DeckServerAPI) ReadCardsFromDeck(ctx context.Context, in *deckv1.ReadDeckRequest) (*deckv1.CardListResponse, error) {
	deckId, err := uuid.Parse(in.DeckId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid deck ID")
	}

	authUser, err := GetAuthUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "User not authenticated")
	}
	cards, err := s.service.ReadAllCardsFromDeck(deckId, authUser.ID)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	var protoCards []*cardv1.Card
	for _, card := range cards {
		protoCards = append(protoCards, convert.ModelToProto(&card))
	}

	return &deckv1.CardListResponse{Cards: protoCards}, nil
}
