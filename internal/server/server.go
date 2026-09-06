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
	Generator      DescriptionGenerator
	LocalGenerator DescriptionGenerator
}

// GenerateDescription handles the gRPC request to generate a sale description
func (s *Server) GenerateDescription(ctx context.Context, req *pb.GenerateDescriptionRequest) (*pb.GenerateDescriptionResponse, error) {
	var gen DescriptionGenerator
	if req.GetUseLocalModel() {
		if s.LocalGenerator == nil {
			return nil, fmt.Errorf("local generator not initialized")
		}
		gen = s.LocalGenerator
	} else {
		if s.Generator == nil {
			return nil, fmt.Errorf("generator not initialized")
		}
		gen = s.Generator
	}

	description, err := gen.Generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate description: %v", err)
	}

	return &pb.GenerateDescriptionResponse{
		Description: description,
	}, nil
}
