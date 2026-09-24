// Package mailketing provides an idiomatic Go client for the Mailketing API v2.
package mailketing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the production Mailketing API v2 base URL.
	DefaultBaseURL = "https://api.mailketing.co.id/api/v2"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 15 * time.Second
)

// Client is a Mailketing API v2 client.
type Client struct {
	apiToken         string
	baseURL          string
	httpClient       *http.Client
	defaultFromName  string
	defaultFromEmail string
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the API base URL.
func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(rawURL, "/")
	}
}

// WithHTTPClient overrides the HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		if client != nil {
			c.httpClient = client
		}
	}
}

// WithDefaultSender sets default sender name and email.
func WithDefaultSender(fromName, fromEmail string) Option {
	return func(c *Client) {
		c.defaultFromName = fromName
		c.defaultFromEmail = fromEmail
	}
}

// NewClient returns a new Client configured with the provided API token and options.
func NewClient(apiToken string, opts ...Option) *Client {
	c := &Client{
		apiToken:   apiToken,
		baseURL:    DefaultBaseURL,
		httpClient: &http.Client{Timeout: DefaultTimeout},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// APIError represents an error returned by the Mailketing API.
type APIError struct {
	StatusCode int            `json:"status_code"`
	Message    string         `json:"message"`
	Errors     map[string]any `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("mailketing: status %d: %s (errors: %v)", e.StatusCode, e.Message, e.Errors)
	}
	return fmt.Sprintf("mailketing: status %d: %s", e.StatusCode, e.Message)
}

// SendEmailRequest represents parameters for sending an email.
type SendEmailRequest struct {
	FromName  string `json:"from_name"`
	FromEmail string `json:"from_email"`
	Subject   string `json:"subject"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	MessageID string `json:"message_id,omitempty"`
	Attach1   string `json:"attach1,omitempty"`
	Attach2   string `json:"attach2,omitempty"`
	Attach3   string `json:"attach3,omitempty"`
}

// SendEmailData contains data returned on successful send.
type SendEmailData struct {
	MessageID string `json:"message_id"`
}

// SendEmailResponse is the response returned by the send endpoint.
type SendEmailResponse struct {
	Success bool           `json:"success"`
	Data    SendEmailData  `json:"data"`
	Message string         `json:"message"`
	Errors  map[string]any `json:"errors,omitempty"`
}

// CreditsData holds credit information.
type CreditsData struct {
	Credits int    `json:"credits"`
	Email   string `json:"email"`
}

// CreditsResponse is the response returned by the credits endpoint.
type CreditsResponse struct {
	Success bool           `json:"success"`
	Data    CreditsData    `json:"data"`
	Message string         `json:"message"`
	Errors  map[string]any `json:"errors,omitempty"`
}

// SendersResponse is the response returned by the senders endpoint.
type SendersResponse struct {
	Success bool           `json:"success"`
	Data    []string       `json:"data"`
	Message string         `json:"message"`
	Errors  map[string]any `json:"errors,omitempty"`
}

// Send queues and sends an email through Mailketing API v2.
func (c *Client) Send(ctx context.Context, req SendEmailRequest) (*SendEmailResponse, error) {
	if req.FromName == "" {
		req.FromName = c.defaultFromName
	}
	if req.FromEmail == "" {
		req.FromEmail = c.defaultFromEmail
	}

	if req.FromEmail == "" {
		return nil, fmt.Errorf("mailketing: from_email is required")
	}
	if req.FromName == "" {
		return nil, fmt.Errorf("mailketing: from_name is required")
	}
	if req.Recipient == "" {
		return nil, fmt.Errorf("mailketing: recipient is required")
	}
	if req.Subject == "" {
		return nil, fmt.Errorf("mailketing: subject is required")
	}
	if req.Content == "" {
		return nil, fmt.Errorf("mailketing: content is required")
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("mailketing: marshal request: %w", err)
	}

	endpoint := c.baseURL + "/send"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("mailketing: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-Api-Token", c.apiToken)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mailketing: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("mailketing: read response body: %w", err)
	}

	var resp SendEmailResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("mailketing: unmarshal response (status %d): %w", httpResp.StatusCode, err)
	}

	if httpResp.StatusCode >= 400 || !resp.Success {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Message:    resp.Message,
			Errors:     resp.Errors,
		}
	}

	return &resp, nil
}

// GetCredits retrieves the remaining account credits.
func (c *Client) GetCredits(ctx context.Context) (*CreditsResponse, error) {
	endpoint := c.baseURL + "/credits"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("mailketing: create request: %w", err)
	}

	httpReq.Header.Set("X-Api-Token", c.apiToken)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mailketing: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("mailketing: read response body: %w", err)
	}

	var resp CreditsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("mailketing: unmarshal response (status %d): %w", httpResp.StatusCode, err)
	}

	if httpResp.StatusCode >= 400 || !resp.Success {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Message:    resp.Message,
			Errors:     resp.Errors,
		}
	}

	return &resp, nil
}

// GetSenders retrieves the list of verified sender email addresses.
func (c *Client) GetSenders(ctx context.Context) (*SendersResponse, error) {
	endpoint := c.baseURL + "/senders"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("mailketing: create request: %w", err)
	}

	httpReq.Header.Set("X-Api-Token", c.apiToken)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("mailketing: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("mailketing: read response body: %w", err)
	}

	var resp SendersResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("mailketing: unmarshal response (status %d): %w", httpResp.StatusCode, err)
	}

	if httpResp.StatusCode >= 400 || !resp.Success {
		return nil, &APIError{
			StatusCode: httpResp.StatusCode,
			Message:    resp.Message,
			Errors:     resp.Errors,
		}
	}

	return &resp, nil
}

// ParseURL ensures a valid URL string helper.
func parseURL(raw string) (*url.URL, error) {
	return url.Parse(raw)
}
