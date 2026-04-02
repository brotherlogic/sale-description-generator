package main

import (
	"log"
	"net"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"github.com/brotherlogic/sale-description-generator/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSaleDescriptionServiceServer(s, &server.Server{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
