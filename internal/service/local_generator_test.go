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
	expectedDesc := "Super crisp vinyl with virtually no signs of play."

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
	}

	desc, err := gen.Generate(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(desc, expectedDesc) {
		t.Errorf("expected description to start with %q, got %q", expectedDesc, desc)
	}

	if !strings.Contains(desc, shippingInfo) {
		t.Errorf("expected description to contain shipping info, got %q", desc)
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
	req := &pb.GenerateDescriptionRequest{
		RecordTitle:     "Meddle",
		Artist:          "Pink Floyd",
		MediaCondition:  pb.Grading_GRADING_VERY_GOOD_PLUS,
		SleeveCondition: pb.Grading_GRADING_VERY_GOOD,
		UserNotes:       "Original harvest pressing",
	}

	prompt := ConstructPrompt(req)
	if !strings.Contains(prompt, "Pink Floyd") || !strings.Contains(prompt, "Meddle") {
		t.Errorf("expected prompt to contain artist and title, got %s", prompt)
	}
	if !strings.Contains(prompt, "GRADING_VERY_GOOD_PLUS") {
		t.Errorf("expected prompt to contain media condition, got %s", prompt)
	}
}
