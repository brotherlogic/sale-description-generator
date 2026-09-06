package service

import (
	"fmt"
	"strings"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

const shippingInfo = "The record is packed outside the sleeve in a protective sleeve and a sturdy mailer."

// ConstructPrompt formats the prompt for the LLM based on the request details.
func ConstructPrompt(req *pb.GenerateDescriptionRequest) string {
	return fmt.Sprintf(`You are an expert music record grader. Generate a concise, positive sale description for a vinyl record.

### Instructions:
1.  **Maximum of three sentences**.
2.  Focus ONLY on the condition and user notes.
3.  Include the artist and title if they fit naturally, but prioritize the condition.
4.  Maintain a high-quality, professional tone.

### Examples of Good Descriptions:
- One Owner, Played a few times with care.
- Played maybe once or twice. Fantastic condition.
- Still has the original hype sticker. The side of the cellophane was carefully cut to play the album once.
- Never played, well kept copy sold in a carefully packaged bubble mailer.

### Current Record Details:
- **Artist**: %s
- **Title**: %s
- **Media Condition**: %s
- **Sleeve Condition**: %s
- **User Notes**: %s

### Final Description:`,
		req.GetArtist(),
		req.GetRecordTitle(),
		req.GetMediaCondition().String(),
		req.GetSleeveCondition().String(),
		req.GetUserNotes(),
	)
}

func appendShippingInfo(description string) string {
	trimmed := strings.TrimSpace(description)
	if trimmed != "" {
		return fmt.Sprintf("%s %s", trimmed, shippingInfo)
	}
	return ""
}
