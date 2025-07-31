package deck

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	deckv1 "github.com/GOeda-Co/proto-contract/gen/go/deck"

	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/tomatoCoderq/repeatro/internal/lib/convert"
	modelDeck "github.com/GOeda-Co/proto-contract/model/deck"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"


	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}

type Client struct {
	api deckv1.DeckServiceClient
	log *slog.Logger
}

func New(
	ctx context.Context,
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadReceived, grpclog.PayloadSent),
	}

	cc, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	grpcClient := deckv1.NewDeckServiceClient(cc)

	return &Client{
		api: grpcClient,
		log: log,
	}, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

// type DeckServiceClient interface {
// 	AddDeck(ctx context.Context, in *AddDeckRequest, opts ...grpc.CallOption) (*DeckResponse, error)
// 	ReadAllDecks(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*DeckListResponse, error)
// 	ReadDeck(ctx context.Context, in *ReadDeckRequest, opts ...grpc.CallOption) (*DeckResponse, error)
// 	DeleteDeck(ctx context.Context, in *ReadDeckRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
// 	AddCardToDeck(ctx context.Context, in *AddCardToDeckRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
// 	ReadCardsFromDeck(ctx context.Context, in *ReadDeckRequest, opts ...grpc.CallOption) (*CardListResponse, error)
// }

func (c *Client) AddDeck(ctx context.Context, deck *modelDeck.Deck) (modelDeck.Deck, error) {
	const op = "grpc.AddDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.AddDeck(ctx, &deckv1.AddDeckRequest{
		Name:        deck.Name,
		Description: deck.Description,
	})
	if err != nil {
		return modelDeck.Deck{}, fmt.Errorf("%s: %w", op, err)
	}

	deckModel, err := convert.ProtoDeckToModel(resp.Deck)
	if err != nil {
		return modelDeck.Deck{}, fmt.Errorf("%s: %w", op, err)
	}
	return *deckModel, nil
}

func (c *Client) ReadAllDecks(ctx context.Context) ([]modelDeck.Deck, error) {
	const op = "grpc.ReadAllDecks"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadAllDecks(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Printf("%s: %s", op, err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	decks := make([]modelDeck.Deck, 0, len(resp.Decks))
	for _, protoDeck := range resp.Decks {
		deck, err := convert.ProtoDeckToModel(protoDeck)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		decks = append(decks, *deck)
	}
	return decks, nil
}

func (c *Client) ReadDeck(ctx context.Context, did uuid.UUID) (modelDeck.Deck, error) {
	const op = "grpc.ReadDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadDeck(ctx, &deckv1.ReadDeckRequest{
		DeckId: did.String(),
	})
	if err != nil {
		return modelDeck.Deck{}, fmt.Errorf("%s: %w", op, err)
	}
	deckModel, err := convert.ProtoDeckToModel(resp.Deck)
	if err != nil {
		return modelDeck.Deck{}, fmt.Errorf("%s: %w", op, err)
	}
	return *deckModel, nil
}

func (c *Client) DeleteDeck(ctx context.Context, did uuid.UUID) error {
	const op = "grpc.DeleteDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	_, err := c.api.DeleteDeck(ctx, &deckv1.ReadDeckRequest{
		DeckId: did.String(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *Client) AddCardToDeck(ctx context.Context, did, cid uuid.UUID) error {
	const op = "grpc.AddCardToDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	_, err := c.api.AddCardToDeck(ctx, &deckv1.AddCardToDeckRequest{
		DeckId: did.String(),
		CardId: cid.String(),
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (c *Client) ReadCardsFromDeck(ctx context.Context, did uuid.UUID) ([]modelCard.Card, error) {
	const op = "grpc.ReadCardsFromDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadCardsFromDeck(ctx, &deckv1.ReadDeckRequest{
		DeckId: did.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]modelCard.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.ProtoToModel(&cardv1.Card{CardId: protoCard.CardId,
			Word:             protoCard.Word,
			Translation:      protoCard.Translation,
			DeckId:           protoCard.DeckId,
			Easiness:         protoCard.Easiness,
			CreatedBy:        protoCard.CreatedBy,
			Interval:         protoCard.Interval,
			ExpiresAt:        timestamppb.New(protoCard.ExpiresAt.AsTime()),
			RepetitionNumber: protoCard.RepetitionNumber,
			Tags:             protoCard.Tags,
			CreatedAt:        timestamppb.New(protoCard.CreatedAt.AsTime()),
			UpdatedAt:        timestamppb.New(protoCard.UpdatedAt.AsTime())})
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
}
