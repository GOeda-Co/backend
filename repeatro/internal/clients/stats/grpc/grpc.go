package stats

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	statv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}

type Client struct {
	api statv1.StatServiceClient
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

	grpcClient := statv1.NewStatServiceClient(cc)

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

func (c *Client) GetAverageGrade(ctx context.Context, userId, deckId string, timeRange statv1.TimeRange) (float64, error) {
	const op = "grpc.AddDeck"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.GetAverageGrade(ctx, &statv1.GetAverageGradeRequest{
		UserId:    userId,
		DeckId:    deckId,
		TimeRange: timeRange,
	})

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.AverageGrade, nil
}

func (c *Client) GetCardsReviewedCount(ctx context.Context, uid, deckId string, timeRange statv1.TimeRange) (int32, error) {
	const op = "grpc.ReadAllDecks"

	ctx = withToken(ctx, ctx.Value("token").(string))

	resp, err := c.api.GetCardsReviewedCount(ctx, &statv1.GetCardsReviewedCountRequest{
		UserId:    uid,
		DeckId:    deckId,
		TimeRange: timeRange,
	})

	fmt.Println("ASASASA")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.ReviewedCount, nil
}
