package grpc

import (
	"context"
	"fmt"

	"github.com/tomatoCoderq/stats/internal/controller"

	statsv1 "github.com/GOeda-Co/proto-contract/gen/go/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	statsv1.UnimplementedStatServiceServer
	service controller.Service
}

func Register(gRPCServer *grpc.Server, card controller.Service) {
	statsv1.RegisterStatServiceServer(gRPCServer, &ServerAPI{service: card})
}

func (s *ServerAPI) GetAverageGrade(ctx context.Context, in *statsv1.GetAverageGradeRequest) (*statsv1.GetAverageGradeResponse, error) {
	// if in.DeckId == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "DeckId is required")
	// }

	if in.TimeRange == statsv1.TimeRange_TIME_RANGE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "TimeRange is not specified")
	}

	
	avGrage, err := s.service.GetAverageGrade(in.UserId, in.DeckId, in.TimeRange)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Error happened: %v", err))
	}

	return &statsv1.GetAverageGradeResponse{
		AverageGrade: avGrage,
	}, nil
}

// GetCardsReviewedCount returns how many cards the user reviewed in a given time range and optional deck
func (s *ServerAPI) GetCardsReviewedCount(ctx context.Context, in *statsv1.GetCardsReviewedCountRequest) (*statsv1.GetCardsReviewedCountResponse, error) {
	// if in.DeckId == "" {
	// 	return nil, status.Error(codes.InvalidArgument, "DeckId is required")
	// }

	if in.TimeRange == statsv1.TimeRange_TIME_RANGE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "TimeRange is not specified")
	}

	revCards, err := s.service.GetCardsReviewedCount(in.UserId, in.DeckId, in.TimeRange)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Error happened: %v", err))
	}

	return &statsv1.GetCardsReviewedCountResponse{
		ReviewedCount: revCards,
	}, nil
}

// GetCardsLearnedCount returns how many cards the user learned in a given time range and optional deck
// func (s *ServerAPI) GetCardsLearnedCount(ctx context.Context, req *statsv1.GetCardsLearnedCountRequest) (*statsv1.GetCardsLearnedCountResponse, error) {
// 	// TODO: implement logic to count "learned" cards (e.g., status or easiness threshold)
// 	fmt.Printf("Counting learned cards for user: %s, deck: %s, range: %v\n", req.UserId, req.DeckId, req.TimeRange)

// 	return &statsv1.GetCardsLearnedCountResponse{
// 		LearnedCount: 17,
// 	}, nil
// }