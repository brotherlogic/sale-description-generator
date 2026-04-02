package server

import (
	"context"
	"fmt"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

// Server implements the SaleDescriptionService gRPC server
type Server struct {
	pb.UnimplementedSaleDescriptionServiceServer
}

// GenerateDescription handles the gRPC request to generate a sale description
func (s *Server) GenerateDescription(ctx context.Context, req *pb.GenerateDescriptionRequest) (*pb.GenerateDescriptionResponse, error) {
	// Placeholder for now
	description := fmt.Sprintf("This is a %s by %s in %s media condition and %s sleeve condition. Notes: %s",
		req.GetRecordTitle(), req.GetArtist(), req.GetMediaCondition(), req.GetSleeveCondition(), req.GetUserNotes())

	return &pb.GenerateDescriptionResponse{
		Description: description,
	}, nil
}
