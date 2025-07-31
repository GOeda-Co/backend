package stats

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	statv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
	"github.com/tomatoCoderq/card/internal/lib/security"

	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func withToken(ctx context.Context, token string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %s", token))
}

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
	const op = "grpc.GetAverageGrade"

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
	const op = "grpc.GetCardsReviewedCount"

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

func (c *Client) AddRecord(ctx context.Context, deckId, cid string, grade int) (string, error) {
	const op = "grpc.AddRecord"

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing metadata in context")
	}

	// Get the authorization token
	authValues := md["authorization"]
	if len(authValues) == 0 {
		return "", fmt.Errorf("authorization token not found in metadata")
	}
	token := authValues[0]

	// Create new outgoing context with the token
	outCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", token))

	fmt.Println("INFODEK", deckId, cid, grade)

	resp, err := c.api.AddRecording(outCtx, &statv1.AddRecordingRequest{
		DeckId:    deckId,
		CardId:    cid,
		CreatedAt: timestamppb.New(time.Now()),
		Grade:     int32(grade),
	})

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.ReviewId, nil
}
