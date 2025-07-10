package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	cardv1 "github.com/GOeda-Co/proto-contract/gen/go/card"
	"github.com/google/uuid"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/tomatoCoderq/repeatro/internal/lib/convert"
	model "github.com/tomatoCoderq/repeatro/pkg/models"
	"github.com/tomatoCoderq/repeatro/pkg/schemes"
	"google.golang.org/grpc"
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

func (c *Client) AddCard(ctx context.Context, card *model.Card) (model.Card, error) {
	const op = "grpc.AddCard"

	resp, err := c.api.AddCard(ctx, &cardv1.AddCardRequest{
		Card: convert.ModelToProto(card),
	})
	if err != nil {
		return model.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	cardModel, err := convert.ProtoToModel(resp.Card)
	if err != nil {
		return model.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	return *cardModel, nil
}

func (c *Client) ReadAllCards(ctx context.Context, uid uuid.UUID) ([]model.Card, error) {
	const op = "grpc.ReadAllCards"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.ReadAllCards(ctx, &cardv1.ReadAllCardsRequest{UserId: uid.String()})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	cards := make([]model.Card, 0, len(resp.Cards))
	for _, protoCard := range resp.Cards {
		card, err := convert.ProtoToModel(protoCard)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		cards = append(cards, *card)
	}
	return cards, nil
}

func (c *Client) UpdateCard(ctx context.Context, cid uuid.UUID, card *schemes.UpdateCardScheme) (model.Card, error) {
	const op = "grpc.UpdateCard"

	ctx = withToken(ctx, ctx.Value("token").(string))
	
	resp, err := c.api.UpdateCard(ctx, &cardv1.UpdateCardRequest{
		CardId: cid.String(),
		Word:  card.Word,
		Translation: card.Translation,
		Easiness: card.Easiness,
		UpdatedAt: timestamppb.New(card.UpdatedAt),
		Interval: int32(card.Interval),
		ExpiresAt: timestamppb.New(card.ExpiresAt),
		RepetitionNumber: int32(card.RepetitionNumber),
		Tags: card.Tags,
	})
	if err != nil {
		return model.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	cardModel, err := convert.ProtoToModel(resp.Card)
	if err != nil {
		return model.Card{}, fmt.Errorf("%s: %w", op, err)
	}
	return *cardModel, nil
}

func (c *Client) DeleteCard(ctx context.Context, cid uuid.UUID, uid uuid.UUID) (bool, error) {
	const op = "grpc.DeleteCard"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.DeleteCard(ctx, &cardv1.DeleteCardRequest{
		CardId: cid.String(),
		UserId: uid.String(),
	})
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return resp.Success, nil
}

func (c *Client) AddAnswers(ctx context.Context, uid uuid.UUID, answers []*schemes.AnswerScheme) (string, error) {
	const op = "grpc.AddAnswers"

	ctx = withToken(ctx, ctx.Value("token").(string))

	convertedAnswers, err := convert.AnswersToProtoSchemes(answers)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	resp, err := c.api.AddAnswers(ctx, &cardv1.AddAnswersRequest{
		Answers: convertedAnswers,
	})
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return resp.Message, nil
}
