package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func parseGrading(s string) (pb.Grading, error) {
	s = strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
	switch s {
	case "MINT", "M":
		return pb.Grading_GRADING_MINT, nil
	case "NEAR_MINT", "NM":
		return pb.Grading_GRADING_NEAR_MINT, nil
	case "VERY_GOOD_PLUS", "VG+", "VG_PLUS":
		return pb.Grading_GRADING_VERY_GOOD_PLUS, nil
	case "VERY_GOOD", "VG":
		return pb.Grading_GRADING_VERY_GOOD, nil
	case "GOOD_PLUS", "G+", "G_PLUS":
		return pb.Grading_GRADING_GOOD_PLUS, nil
	case "GOOD", "G":
		return pb.Grading_GRADING_GOOD, nil
	case "FAIR", "F":
		return pb.Grading_GRADING_FAIR, nil
	case "POOR", "P":
		return pb.Grading_GRADING_POOR, nil
	default:
		return pb.Grading_GRADING_UNSPECIFIED, fmt.Errorf("unknown grading: %s", s)
	}
}

func main() {
	serverAddr := flag.String("server", "localhost:50051", "The server address and port")
	title := flag.String("title", "", "Record title")
	artist := flag.String("artist", "", "Artist name")
	media := flag.String("media", "VG+", "Media condition")
	sleeve := flag.String("sleeve", "VG+", "Sleeve condition")
	notes := flag.String("notes", "", "Additional user notes")
	flag.Parse()

	if *title == "" || *artist == "" {
		log.Fatalf("title and artist are required")
	}

	mediaGrading, err := parseGrading(*media)
	if err != nil {
		log.Fatalf("invalid media condition: %v", err)
	}

	sleeveGrading, err := parseGrading(*sleeve)
	if err != nil {
		log.Fatalf("invalid sleeve condition: %v", err)
	}

	// Connect to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to %s: %v", *serverAddr, err)
	}
	defer conn.Close()

	client := pb.NewSaleDescriptionServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	resp, err := client.GenerateDescription(ctx, &pb.GenerateDescriptionRequest{
		RecordTitle:     *title,
		Artist:          *artist,
		MediaCondition:  mediaGrading,
		SleeveCondition: sleeveGrading,
		UserNotes:       *notes,
	})
	if err != nil {
		log.Fatalf("failed to generate description: %v", err)
	}

	fmt.Printf("Generated Description for %s by %s:\n\n%s\n", *title, *artist, resp.GetDescription())
}
