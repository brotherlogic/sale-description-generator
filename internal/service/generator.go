package service

import (
	"context"
	"fmt"
	"os"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Generator handles the business logic for description generation
type Generator struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

// NewGenerator creates a new description generator
func NewGenerator(ctx context.Context) (*Generator, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %v", err)
	}

	// Use gemini-2.5-flash as the standard for this project
	model := client.GenerativeModel("gemini-2.5-flash")

	return &Generator{
		client: client,
		model:  model,
	}, nil
}

// Close closes the gemini client
func (g *Generator) Close() {
	if g.client != nil {
		g.client.Close()
	}
}

// Generate produces a sale description based on the request
func (g *Generator) Generate(ctx context.Context, req *pb.GenerateDescriptionRequest) (string, error) {
	prompt := ConstructPrompt(req)

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %v", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("no response from gemini")
	}

	// Extract text from response
	var description string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			description += string(text)
		}
	}

	return appendShippingInfo(description), nil
}
