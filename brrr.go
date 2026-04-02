// Package brrr
// see more details https://brrr.now/docs/
package brrr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const defaultBaseURL = "https://api.brrr.now/v1/send"

// Sound represents a notification sound. Only supported on iPhone and iPad.
type Sound string

const (
	SoundDefault           Sound = "default"
	SoundSystem            Sound = "system"
	SoundBrrr              Sound = "brrr"
	SoundBellRinging       Sound = "bell_ringing"
	SoundBubbleDing        Sound = "bubble_ding"
	SoundBubblySuccessDing Sound = "bubbly_success_ding"
	SoundCatMeow           Sound = "cat_meow"
	SoundCalm1             Sound = "calm1"
	SoundCalm2             Sound = "calm2"
	SoundChaChing          Sound = "cha_ching"
	SoundDogBarking        Sound = "dog_barking"
	SoundDoorBell          Sound = "door_bell"
	SoundDuckQuack         Sound = "duck_quack"
	SoundShortTripleBlink  Sound = "short_triple_blink"
	SoundUpbeatBells       Sound = "upbeat_bells"
	SoundWarmSoftError     Sound = "warm_soft_error"
)

// Omit this field to use the system default behavior. 
// Passive adds the notification to the system's notification list without lighting up the screen or playing a sound.
// Active presents the notification, lights up the screen, and can play a sound.
// Time-sensitive presents the notification immediately, lights up the screen, can play a sound, and breaks through controls such as Notification Summary and Focus.
type InterruptionLevel string

const (
	InterruptionPassive       InterruptionLevel = "passive"
	InterruptionActive        InterruptionLevel = "active"
	InterruptionTimeSensitive InterruptionLevel = "time-sensitive"
)

// Notification the payload sent to the brrr API.
type Notification struct {
	Title             string            `json:"title,omitempty"`              // First line of the notification, if present.
	Subtitle          string            `json:"subtitle,omitempty"`           // Line under the title, if present.
	Message           string            `json:"message,omitempty"`            // Primary content of the notification.
	Sound             Sound             `json:"sound,omitempty"`              // Sound to play when the notification is delivered.
	OpenURL           string            `json:"open_url,omitempty"`           // Link to open when the notification is selected.
	ImageURL          string            `json:"image_url,omitempty"`          // Link to an image to be displayed in the notification.
	ExpirationDate    *time.Time        `json:"expiration_date,omitempty"`    // ISO 8601 date and time. In case the notification isn't successfully delivered when the webhook is invoked, Apple Push Notification Service may retry until the expiration date is reached.
	FilterCriteria    string            `json:"filter-criteria,omitempty"`    // Optional criterion used to match a Focus filter on the device. If the device’s active Focus allows this criterion, the notification can be shown. Focus filtering must be configured on the device.
	InterruptionLevel InterruptionLevel `json:"interruption-level,omitempty"` // see InterruptionLevel for more details.
	ThreadID          string            `json:"thread_id,omitempty"`          // Identifier used to group related notifications together in Notification Center.
}

// Client send notifications through the brrr API.
type Client struct {
	secret     string
	baseURL    string
	httpClient *http.Client
	logger     *slog.Logger
}

// Option configure a Client.
type Option func(*Client)

// WithHTTPClient set a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithLogger set a logger. By default no logging is performed.
func WithLogger(l *slog.Logger) Option {
	return func(cl *Client) { cl.logger = l }
}

// New create a Client for the given brrr secret (e.g. "br_usr_...").
func New(secret string, opts ...Option) (*Client, error) {
	if secret == "" {
		return nil, fmt.Errorf("brrr: secret must not be empty")
	}
	c := &Client{
		secret:     secret,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for _, o := range opts {
		o(c)
	}
	return c, nil
}

// Send send a fully customizable notification.
func (c *Client) Send(ctx context.Context, n Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("brrr: marshal notification: %w", err)
	}

	c.logger.DebugContext(ctx, "sending notification", "url", c.baseURL, "payload_size", len(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brrr: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.secret)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.ErrorContext(ctx, "request failed", "error", err)
		return fmt.Errorf("brrr: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		c.logger.ErrorContext(ctx, "unexpected status", "status", resp.StatusCode, "body", string(respBody))
		return fmt.Errorf("brrr: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	c.logger.InfoContext(ctx, "notification sent", "status", resp.StatusCode)
	return nil
}

// SendMessage send a simple message notification.
func (c *Client) SendMessage(ctx context.Context, message string) error {
	return c.Send(ctx, Notification{Message: message})
}

// SendWithTitle send a title and message notification.
func (c *Client) SendWithTitle(ctx context.Context, title string, message string) error {
	return c.Send(ctx, Notification{Title: title, Message: message})
}
