package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/GOeda-Co/proto-contract/convert"
	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	// model "github.com/tomatoCoderq/repeatro/pkg/models"
	modelCard "github.com/GOeda-Co/proto-contract/model/card"
	// "github.com/tomatoCoderq/repeatro/pkg/schemes"
	schemes "github.com/GOeda-Co/proto-contract/scheme/card"

	"google.golang.org/grpc"
	// "google.golang.org/protobuf/types/known/emptypb"
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
	api cardv1.CardServiceClient
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

	grpcClient := cardv1.NewCardServiceClient(cc)

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

func (c *Client) AddCard(ctx context.Context, card *modelCard.Card) (modelCard.Card, error) {
	const op = "grpc.AddCard"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.AddCard(ctx, &cardv1.AddCardRequest{
		Card: convert.FromModelToProtoCard(card),
	})
	if err != nil {
		return modelCard.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	cardModel, err := convert.FromProtoToModelCard(resp.Card)
	if err != nil {
		return modelCard.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	return *cardModel, nil
}

func (c *Client) ReadAllCardsToLearn(ctx context.Context, uid uuid.UUID) ([]modelCard.Card, error) {
	const op = "grpc.ReadAllCards"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadAllOwnCardsToLearn(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]modelCard.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.FromProtoToModelCard(protoCard)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
}

func (c *Client) ReadAllCards(ctx context.Context, uid uuid.UUID) ([]modelCard.Card, error) {
	const op = "grpc.ReadAllCards"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadAllOwnCards(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]modelCard.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.FromProtoToModelCard(protoCard)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
}

func (c *Client) SearchAllPublicCards(ctx context.Context) ([]modelCard.Card, error) {
	const op = "grpc.SearchAllPublicCards"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.SearchAllPublicCards(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]modelCard.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.FromProtoToModelCard(protoCard)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
}

func (c *Client) SearchUserPublicCards(ctx context.Context, uid uuid.UUID) ([]modelCard.Card, error) {
	const op = "grpc.SearchUserPublicCards"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.SearchUserPublicCards(ctx, &cardv1.SearchUserPublicCardsRequest{UserId: uid.String()})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]modelCard.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.FromProtoToModelCard(protoCard)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
	// return nil, nil
}

func (c *Client) UpdateCard(ctx context.Context, uid uuid.UUID, cid uuid.UUID, card *schemes.UpdateCardScheme) (modelCard.Card, error) {
	const op = "grpc.UpdateCard"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.UpdateCard(ctx, &cardv1.UpdateCardRequest{
		CardId:           cid.String(),
		UserId:           uid.String(),
		Word:             card.Word,
		Translation:      card.Translation,
		Easiness:         card.Easiness,
		UpdatedAt:        timestamppb.New(card.UpdatedAt),
		Interval:         int32(card.Interval),
		ExpiresAt:        timestamppb.New(card.ExpiresAt),
		RepetitionNumber: int32(card.RepetitionNumber),
		Tags:             card.Tags,
	})
	if err != nil {
		return modelCard.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	cardModel, err := convert.FromProtoToModelCard(resp.Card)
	if err != nil {
		return modelCard.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	return *cardModel, nil
}

func (c *Client) DeleteCard(ctx context.Context, cid uuid.UUID, uid uuid.UUID) (bool, error) {
	const op = "grpc.DeleteCard"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.DeleteCard(ctx, &cardv1.DeleteCardRequest{
		CardId: cid.String(),
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.Success, nil
}

func (c *Client) AddAnswers(ctx context.Context, uid uuid.UUID, answers []*schemes.AnswerScheme) (string, error) {
	const op = "grpc.AddAnswers"

	ctx = withToken(ctx, ctx.Value("token").(string))

	convertedAnswers, err := convert.FromAnswerSchemesToProtosCard(answers)
	if err != nil {
		fmt.Printf("%s: %v", op, err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	resp, err := c.api.AddAnswers(ctx, &cardv1.AddAnswersRequest{
		Answers: convertedAnswers,
	})
	if err != nil {
		fmt.Printf("%s: %v", op, err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.Message, nil
}
