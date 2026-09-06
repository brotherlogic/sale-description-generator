package server

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type mockGenerator struct {
	prefix string
}

func (m *mockGenerator) Generate(ctx context.Context, req *pb.GenerateDescriptionRequest) (string, error) {
	if m.prefix == "" {
		return "Mock description for testing", nil
	}
	return fmt.Sprintf("%s: %s by %s", m.prefix, req.GetRecordTitle(), req.GetArtist()), nil
}

func setupTestServer(t *testing.T, srv *Server) pb.SaleDescriptionServiceClient {
	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSaleDescriptionServiceServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Logf("Server exited: %v", err)
		}
	}()
	t.Cleanup(func() {
		s.Stop()
		lis.Close()
	})

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})

	return pb.NewSaleDescriptionServiceClient(conn)
}

func TestGenerateDescription_DefaultGenerator(t *testing.T) {
	client := setupTestServer(t, &Server{
		Generator:      &mockGenerator{prefix: "Gemini"},
		LocalGenerator: &mockGenerator{prefix: "Local"},
	})

	req := &pb.GenerateDescriptionRequest{
		RecordTitle:     "Wish You Were Here",
		Artist:          "Pink Floyd",
		MediaCondition:  pb.Grading_GRADING_VERY_GOOD_PLUS,
		SleeveCondition: pb.Grading_GRADING_NEAR_MINT,
		UserNotes:       "Slight edge wear on sleeve",
		UseLocalModel:   false,
	}

	resp, err := client.GenerateDescription(context.Background(), req)
	if err != nil {
		t.Fatalf("GenerateDescription failed: %v", err)
	}

	if !strings.HasPrefix(resp.GetDescription(), "Gemini:") {
		t.Errorf("expected Gemini generator to be used, got: %s", resp.GetDescription())
	}
}

func TestGenerateDescription_LocalGenerator(t *testing.T) {
	client := setupTestServer(t, &Server{
		Generator:      &mockGenerator{prefix: "Gemini"},
		LocalGenerator: &mockGenerator{prefix: "Local"},
	})

	req := &pb.GenerateDescriptionRequest{
		RecordTitle:     "Wish You Were Here",
		Artist:          "Pink Floyd",
		MediaCondition:  pb.Grading_GRADING_VERY_GOOD_PLUS,
		SleeveCondition: pb.Grading_GRADING_NEAR_MINT,
		UserNotes:       "Slight edge wear on sleeve",
		UseLocalModel:   true,
	}

	resp, err := client.GenerateDescription(context.Background(), req)
	if err != nil {
		t.Fatalf("GenerateDescription failed: %v", err)
	}

	if !strings.HasPrefix(resp.GetDescription(), "Local:") {
		t.Errorf("expected Local generator to be used, got: %s", resp.GetDescription())
	}
}

func TestGenerateDescription_UninitializedGenerators(t *testing.T) {
	client := setupTestServer(t, &Server{})

	_, err := client.GenerateDescription(context.Background(), &pb.GenerateDescriptionRequest{
		UseLocalModel: false,
	})
	if err == nil || !strings.Contains(err.Error(), "generator not initialized") {
		t.Errorf("expected 'generator not initialized' error, got: %v", err)
	}

	_, err = client.GenerateDescription(context.Background(), &pb.GenerateDescriptionRequest{
		UseLocalModel: true,
	})
	if err == nil || !strings.Contains(err.Error(), "local generator not initialized") {
		t.Errorf("expected 'local generator not initialized' error, got: %v", err)
	}
}
