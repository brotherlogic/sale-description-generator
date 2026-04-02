package server

import (
	"context"
	"log"
	"net"
	"testing"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type mockGenerator struct{}

func (m *mockGenerator) Generate(ctx context.Context, req *pb.GenerateDescriptionRequest) (string, error) {
	return "Mock description for testing", nil
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSaleDescriptionServiceServer(s, &Server{
		Generator: &mockGenerator{},
	})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGenerateDescription(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSaleDescriptionServiceClient(conn)

	req := &pb.GenerateDescriptionRequest{
		RecordTitle:     "Wish You Were Here",
		Artist:          "Pink Floyd",
		MediaCondition:  pb.Grading_GRADING_VERY_GOOD_PLUS,
		SleeveCondition: pb.Grading_GRADING_NEAR_MINT,
		UserNotes:       "Slight edge wear on sleeve",
	}

	resp, err := client.GenerateDescription(ctx, req)
	if err != nil {
		t.Fatalf("GenerateDescription failed: %v", err)
	}

	if resp.GetDescription() == "" {
		t.Error("expected description to be non-empty")
	}

	if len(resp.GetDescription()) < 10 {
		t.Errorf("expected description to be reasonably long, got: %s", resp.GetDescription())
	}
}
