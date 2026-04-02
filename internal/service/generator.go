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
	return fmt.Sprintf(`You are an expert music record grader and professional copywriter for music collectors. 
Your task is to generate a high-quality, positive, and accurate sale description for a vinyl record.

### Record Details:
- **Artist**: %s
- **Title**: %s
- **Media Condition**: %s
- **Sleeve Condition**: %s
- **User Notes**: %s

### Grading Standards (Discogs/Goldmine Basis):
- **NM (Near Mint)**: Shiny, no visible defects, nearly perfect. Plays with no surface noise. Covers are free of creases, ring wear, and seam splits.
- **VG+ (Very Good Plus)**: Slight signs of wear, light scuffs/scratches that don't affect listening. Minor signs of handling. Covers have minor wear, small seam splits (less than 1 inch).
- **VG (Very Good)**: Obvious flaws, light surface noise especially in soft passages but not overpowering. Covers have ring wear, creases, or more obvious seam splits.
- **G+ (Good Plus)**: Significant surface noise and groove wear, but plays through without skipping. Covers have heavy wear, splitting, or writing.

### Instructions:
1. Generate a **single long paragraph** describing the item.
2. Maintain a **professional and positive tone** to encourage buyers.
3. Be **strictly accurate** to the provided grading. If the condition is VG, do not describe it as looking new.
4. Integrate the user's specific notes naturally into the description.
5. Focus on why a collector would want this specific copy.
6. Do not use placeholders; output the final description text directly.`,
		req.GetArtist(),
		req.GetRecordTitle(),
		req.GetMediaCondition().String(),
		req.GetSleeveCondition().String(),
		req.GetUserNotes(),
	)
}
