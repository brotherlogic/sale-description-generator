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
	prompt := g.constructPrompt(req)

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

	return description, nil
}

func (g *Generator) constructPrompt(req *pb.GenerateDescriptionRequest) string {
	return fmt.Sprintf(`"Act as an expert record grader and Discogs seller. I am going to provide you with the Artist, Album Title, and Condition Notes for a vinyl record.

Please write a concise, professional listing description following these rules:

No Marketing Fluff: Do not explain why the album is 'classic' or 'influential.' The buyer already knows.

Structure: Break the description into three clear sections: Media, Sleeve, and Overall.

Technical Specifics: Use standard collector terminology (e.g., 'hairlines,' 'shelf wear,' 'seam splits,' 'surface noise').

Tone: Professional, honest, and minimalist.

Length: Keep the entire description under 80 words.

Here are the details:

Artist/Album: [Insert Name Here]

Media Grade: [e.g., VG+]

Sleeve Grade: [e.g., VG]

Specific Defects/Notes: [e.g., minor corner ding, plays through perfectly, includes original insert]"`,
		req.GetArtist(),
		req.GetRecordTitle(),
		req.GetMediaCondition().String(),
		req.GetSleeveCondition().String(),
		req.GetUserNotes(),
	)
}
