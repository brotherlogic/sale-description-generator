package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pb "github.com/brotherlogic/sale-description-generator/api/gen/v1"
)

func TestLocalGenerator_Generate_Success(t *testing.T) {
	expectedModel := "test-model"
	expectedDesc := "Impeccably preserved copy with virtually no signs of play."

	testCases := []struct {
		name             string
		format           pb.MediaFormat
		expectedShipping string
	}{
		{
			name:             "Vinyl",
			format:           pb.MediaFormat_MEDIA_FORMAT_VINYL,
			expectedShipping: vinylShippingInfo,
		},
		{
			name:             "CD",
			format:           pb.MediaFormat_MEDIA_FORMAT_CD,
			expectedShipping: cdShippingInfo,
		},
		{
			name:             "Mixed",
			format:           pb.MediaFormat_MEDIA_FORMAT_MIXED,
			expectedShipping: genericShippingInfo,
		},
		{
			name:             "Unspecified",
			format:           pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			expectedShipping: genericShippingInfo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/chat/completions" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				if r.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", r.Method)
				}

				var reqBody chatCompletionRequest
				if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
					t.Errorf("failed to decode request body: %v", err)
				}

				if reqBody.Model != expectedModel {
					t.Errorf("expected model %s, got %s", expectedModel, reqBody.Model)
				}

				if len(reqBody.Messages) == 0 || !strings.Contains(reqBody.Messages[0].Content, "Wish You Were Here") {
					t.Errorf("expected prompt containing 'Wish You Were Here', got: %+v", reqBody.Messages)
				}

				resp := chatCompletionResponse{
					Choices: []struct {
						Message struct {
							Content string `json:"content"`
						} `json:"message"`
					}{
						{
							Message: struct {
								Content string `json:"content"`
							}{
								Content: expectedDesc,
							},
						},
					},
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			gen := NewLocalGenerator(server.URL, expectedModel)

			req := &pb.GenerateDescriptionRequest{
				RecordTitle:     "Wish You Were Here",
				Artist:          "Pink Floyd",
				MediaCondition:  pb.Grading_GRADING_NEAR_MINT,
				SleeveCondition: pb.Grading_GRADING_NEAR_MINT,
				UserNotes:       "Played once",
				MediaFormat:     tc.format,
			}

			desc, err := gen.Generate(context.Background(), req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !strings.HasPrefix(desc, expectedDesc) {
				t.Errorf("expected description to start with %q, got %q", expectedDesc, desc)
			}

			if !strings.HasSuffix(desc, tc.expectedShipping) {
				t.Errorf("expected description to end with %q, got %q", tc.expectedShipping, desc)
			}
		})
	}
}

func TestLocalGenerator_Generate_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal model failure"))
	}))
	defer server.Close()

	gen := NewLocalGenerator(server.URL, "test-model")

	req := &pb.GenerateDescriptionRequest{
		RecordTitle: "Animals",
		Artist:      "Pink Floyd",
	}

	_, err := gen.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error mentioning status 500, got: %v", err)
	}
}

func TestLocalGenerator_Generate_ModelErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := chatCompletionResponse{
			Error: &struct {
				Message string `json:"message"`
			}{
				Message: "model not loaded",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	gen := NewLocalGenerator(server.URL, "test-model")

	req := &pb.GenerateDescriptionRequest{
		RecordTitle: "Animals",
		Artist:      "Pink Floyd",
	}

	_, err := gen.Generate(context.Background(), req)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "model not loaded") {
		t.Errorf("expected error mentioning 'model not loaded', got: %v", err)
	}
}

func TestConstructPrompt(t *testing.T) {
	testCases := []struct {
		name           string
		format         pb.MediaFormat
		expectedPhrases []string
	}{
		{
			name:   "Vinyl",
			format: pb.MediaFormat_MEDIA_FORMAT_VINYL,
			expectedPhrases: []string{
				"You are an expert music record grader. Generate a concise, positive sale description for a vinyl record.",
				"- **Format**: MEDIA_FORMAT_VINYL",
			},
		},
		{
			name:   "CD",
			format: pb.MediaFormat_MEDIA_FORMAT_CD,
			expectedPhrases: []string{
				"You are an expert physical music grader. Generate a concise, positive sale description for a CD.",
				"Use general terms for the disc and packaging (e.g., CD, disc, packaging, case, insert) without assuming a specific case style (such as jewel case vs. digipak) unless mentioned in the user notes. Strictly avoid vinyl-specific terminology (such as vinyl, record, LP, wax, or grooves).",
				"- **Format**: MEDIA_FORMAT_CD",
			},
		},
		{
			name:   "Mixed",
			format: pb.MediaFormat_MEDIA_FORMAT_MIXED,
			expectedPhrases: []string{
				"You are an expert physical music grader. Generate a concise, positive sale description for a physical music release.",
				"Use general, media-neutral terms throughout for the media and packaging with no assumptions about the specific physical medium. Strictly avoid format-specific terminology (such as vinyl, record, LP, wax, grooves, or CD jewel case) unless explicitly mentioned in the user notes.",
				"- **Format**: MEDIA_FORMAT_MIXED",
			},
		},
		{
			name:   "Unspecified",
			format: pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			expectedPhrases: []string{
				"You are an expert physical music grader. Generate a concise, positive sale description for a physical music release.",
				"Use general, media-neutral terms throughout for the media and packaging with no assumptions about the specific physical medium. Strictly avoid format-specific terminology (such as vinyl, record, LP, wax, grooves, or CD jewel case) unless explicitly mentioned in the user notes.",
				"- **Format**: MEDIA_FORMAT_UNSPECIFIED",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.GenerateDescriptionRequest{
				RecordTitle:     "Meddle",
				Artist:          "Pink Floyd",
				MediaCondition:  pb.Grading_GRADING_VERY_GOOD_PLUS,
				SleeveCondition: pb.Grading_GRADING_VERY_GOOD,
				UserNotes:       "Original pressing",
				MediaFormat:     tc.format,
			}

			prompt := ConstructPrompt(req)
			if !strings.Contains(prompt, "Pink Floyd") || !strings.Contains(prompt, "Meddle") {
				t.Errorf("expected prompt to contain artist and title, got %s", prompt)
			}
			if !strings.Contains(prompt, "GRADING_VERY_GOOD_PLUS") {
				t.Errorf("expected prompt to contain media condition, got %s", prompt)
			}
			for _, phrase := range tc.expectedPhrases {
				if !strings.Contains(prompt, phrase) {
					t.Errorf("expected prompt for %s to contain %q, got:\n%s", tc.name, phrase, prompt)
				}
			}
		})
	}
}

func TestAppendShippingInfo(t *testing.T) {
	baseDesc := "Clean copy in great shape."

	testCases := []struct {
		name        string
		description string
		format      pb.MediaFormat
		expected    string
	}{
		{
			name:        "Vinyl",
			description: baseDesc,
			format:      pb.MediaFormat_MEDIA_FORMAT_VINYL,
			expected:    baseDesc + " " + vinylShippingInfo,
		},
		{
			name:        "CD",
			description: baseDesc,
			format:      pb.MediaFormat_MEDIA_FORMAT_CD,
			expected:    baseDesc + " " + cdShippingInfo,
		},
		{
			name:        "Mixed",
			description: baseDesc,
			format:      pb.MediaFormat_MEDIA_FORMAT_MIXED,
			expected:    baseDesc + " " + genericShippingInfo,
		},
		{
			name:        "Unspecified",
			description: baseDesc,
			format:      pb.MediaFormat_MEDIA_FORMAT_UNSPECIFIED,
			expected:    baseDesc + " " + genericShippingInfo,
		},
		{
			name:        "UnknownEnumValue",
			description: baseDesc,
			format:      pb.MediaFormat(99),
			expected:    baseDesc + " " + genericShippingInfo,
		},
		{
			name:        "EmptyDescription_Vinyl",
			description: "",
			format:      pb.MediaFormat_MEDIA_FORMAT_VINYL,
			expected:    "",
		},
		{
			name:        "WhitespaceOnly_CD",
			description: "   \t\n  ",
			format:      pb.MediaFormat_MEDIA_FORMAT_CD,
			expected:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := appendShippingInfo(tc.description, tc.format)
			if got != tc.expected {
				t.Errorf("appendShippingInfo(%q, %v) = %q; want %q", tc.description, tc.format, got, tc.expected)
			}
		})
	}
}
