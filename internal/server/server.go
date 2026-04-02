package server

import (
	"context"
	"fmt"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

// DescriptionGenerator defines the interface for generating sale descriptions
type DescriptionGenerator interface {
	Generate(ctx context.Context, req *pb.GenerateDescriptionRequest) (string, error)
}

// Server implements the SaleDescriptionService gRPC server
type Server struct {
	pb.UnimplementedSaleDescriptionServiceServer
	Generator DescriptionGenerator
}

// GenerateDescription handles the gRPC request to generate a sale description
func (s *Server) GenerateDescription(ctx context.Context, req *pb.GenerateDescriptionRequest) (*pb.GenerateDescriptionResponse, error) {
	if s.Generator == nil {
		return nil, fmt.Errorf("generator not initialized")
	}

	description, err := s.Generator.Generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate description: %v", err)
	}

	return &pb.GenerateDescriptionResponse{
		Description: description,
	}, nil
}
