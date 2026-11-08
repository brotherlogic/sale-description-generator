package service

import (
	"fmt"
	"strings"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

const (
	vinylShippingInfo   = "The record is packed outside the sleeve in a protective sleeve and a sturdy mailer."
	cdShippingInfo      = "The CD is packaged in a bubble mailer for protection."
	genericShippingInfo = "The item is carefully packaged for protection during shipping."
)

func formatPromptDetails(format pb.MediaFormat) (openingInstruction, formatConstraint string) {
	switch format {
	case pb.MediaFormat_MEDIA_FORMAT_VINYL:
		return "You are an expert music record grader. Generate a concise, positive sale description for a vinyl record.", ""
	case pb.MediaFormat_MEDIA_FORMAT_CD:
		return "You are an expert physical music grader. Generate a concise, positive sale description for a CD.",
			"\n5.  Use general terms for the disc and packaging (e.g., CD, disc, packaging, case, insert) without assuming a specific case style (such as jewel case vs. digipak) unless mentioned in the user notes. Strictly avoid vinyl-specific terminology (such as vinyl, record, LP, wax, or grooves)."
	default:
		return "You are an expert physical music grader. Generate a concise, positive sale description for a physical music release.",
			"\n5.  Use general, media-neutral terms throughout for the media and packaging with no assumptions about the specific physical medium. Strictly avoid format-specific terminology (such as vinyl, record, LP, wax, grooves, or CD jewel case) unless explicitly mentioned in the user notes."
	}
}

// ConstructPrompt formats the prompt for the LLM based on the request details.
func ConstructPrompt(req *pb.GenerateDescriptionRequest) string {
	openingInstruction, formatConstraint := formatPromptDetails(req.GetMediaFormat())

	return fmt.Sprintf(`%s

### Instructions:
1.  **Maximum of three sentences**.
2.  Focus ONLY on the condition and user notes.
3.  Include the artist and title if they fit naturally, but prioritize the condition.
4.  Maintain a high-quality, professional tone.%s

### Examples of Good Descriptions:
- One Owner, Played a few times with care.
- Played maybe once or twice. Fantastic condition.
- Still has the original hype sticker. The side of the cellophane was carefully cut to play the album once.
- Never played, well kept copy sold in a carefully packaged bubble mailer.

### Current Record Details:
- **Artist**: %s
- **Title**: %s
- **Format**: %s
- **Media Condition**: %s
- **Sleeve Condition**: %s
- **User Notes**: %s

### Final Description:`,
		openingInstruction,
		formatConstraint,
		req.GetArtist(),
		req.GetRecordTitle(),
		req.GetMediaFormat().String(),
		req.GetMediaCondition().String(),
		req.GetSleeveCondition().String(),
		req.GetUserNotes(),
	)
}

func appendShippingInfo(description string, format pb.MediaFormat) string {
	trimmed := strings.TrimSpace(description)
	if trimmed == "" {
		return ""
	}
	var suffix string
	switch format {
	case pb.MediaFormat_MEDIA_FORMAT_VINYL:
		suffix = vinylShippingInfo
	case pb.MediaFormat_MEDIA_FORMAT_CD:
		suffix = cdShippingInfo
	default:
		suffix = genericShippingInfo
	}
	return fmt.Sprintf("%s %s", trimmed, suffix)
}
