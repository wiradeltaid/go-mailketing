package mailketing_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/wiradeltaid/go-mailketing"
)

func TestSend_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/send" {
			t.Errorf("expected /send, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Api-Token") != "test-token" {
			t.Errorf("expected test-token, got %s", r.Header.Get("X-Api-Token"))
		}

		var req mailketing.SendEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Recipient != "user@example.com" {
			t.Errorf("expected recipient user@example.com, got %s", req.Recipient)
		}

		resp := mailketing.SendEmailResponse{
			Success: true,
			Data: mailketing.SendEmailData{
				MessageID: "<api.123@example.com>",
			},
			Message: "Email queued successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := mailketing.NewClient("test-token",
		mailketing.WithBaseURL(server.URL),
		mailketing.WithDefaultSender("Default Sender", "default@example.com"),
	)

	req := mailketing.SendEmailRequest{
		Recipient: "user@example.com",
		Subject:   "Welcome!",
		Content:   "Hello World",
	}

	res, err := client.Send(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}
	if res.Data.MessageID != "<api.123@example.com>" {
		t.Errorf("expected message_id <api.123@example.com>, got %s", res.Data.MessageID)
	}
}

func TestSend_ValidationErrors(t *testing.T) {
	client := mailketing.NewClient("test-token")

	// Missing recipient
	_, err := client.Send(context.Background(), mailketing.SendEmailRequest{
		FromEmail: "from@example.com",
		FromName:  "From",
		Subject:   "Subject",
		Content:   "Content",
	})
	if err == nil {
		t.Errorf("expected error for missing recipient")
	}

	// Missing subject
	_, err = client.Send(context.Background(), mailketing.SendEmailRequest{
		FromEmail: "from@example.com",
		FromName:  "From",
		Recipient: "to@example.com",
		Content:   "Content",
	})
	if err == nil {
		t.Errorf("expected error for missing subject")
	}
}

func TestSend_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		resp := mailketing.SendEmailResponse{
			Success: false,
			Message: "Validation failed",
			Errors: map[string]any{
				"from_email": "Domain not verified",
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := mailketing.NewClient("test-token", mailketing.WithBaseURL(server.URL))
	_, err := client.Send(context.Background(), mailketing.SendEmailRequest{
		FromEmail: "unverified@example.com",
		FromName:  "Sender",
		Recipient: "recipient@example.com",
		Subject:   "Test",
		Content:   "Hello",
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	apiErr, ok := err.(*mailketing.APIError)
	if !ok {
		t.Fatalf("expected *mailketing.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", apiErr.StatusCode)
	}
}

func TestGetCredits_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/credits" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		resp := mailketing.CreditsResponse{
			Success: true,
			Data: mailketing.CreditsData{
				Credits: 1500,
				Email:   "owner@wiradelta.com",
			},
			Message: "OK",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := mailketing.NewClient("test-token", mailketing.WithBaseURL(server.URL))
	credits, err := client.GetCredits(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if credits.Data.Credits != 1500 {
		t.Errorf("expected 1500 credits, got %d", credits.Data.Credits)
	}
}

func TestGetSenders_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/senders" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		resp := mailketing.SendersResponse{
			Success: true,
			Data:    []string{"notification@wiradelta.com", "admin@wiradelta.com"},
			Message: "OK",
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := mailketing.NewClient("test-token", mailketing.WithBaseURL(server.URL))
	senders, err := client.GetSenders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(senders.Data) != 2 {
		t.Errorf("expected 2 senders, got %d", len(senders.Data))
	}
}

func TestContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	client := mailketing.NewClient("test-token", mailketing.WithBaseURL(server.URL))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.GetCredits(ctx)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
