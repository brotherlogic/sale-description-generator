package main

import (
	"context"
	"log"
	"net"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"github.com/brotherlogic/sale-description-generator/internal/server"
	"github.com/brotherlogic/sale-description-generator/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx := context.Background()

	// Initialize Gemini Generator
	var gen *service.Generator
	var err error
	gen, err = service.NewGenerator(ctx)
	if err != nil {
		log.Printf("Warning: failed to initialize Gemini generator: %v", err)
	} else {
		defer gen.Close()
	}

	// Initialize Local Ollama Generator
	localGen := service.NewLocalGenerator("", "")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterSaleDescriptionServiceServer(s, &server.Server{
		Generator:      gen,
		LocalGenerator: localGen,
	})

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
