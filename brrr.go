// Package brrr sends push notifications to iOS devices via the brrr.now webhook API.
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

var apiBaseURL = "https://api.brrr.now/v1/"

// Sound represents a notification sound.
type Sound string

const (
	SoundDefault          Sound = "default"
	SoundSystem           Sound = "system"
	SoundBrrr             Sound = "brrr"
	SoundBellRinging      Sound = "bell_ringing"
	SoundBubbleDing       Sound = "bubble_ding"
	SoundBubblySuccessDng Sound = "bubbly_success_ding"
	SoundCatMeow          Sound = "cat_meow"
	SoundCalm1            Sound = "calm1"
	SoundCalm2            Sound = "calm2"
	SoundChaChing         Sound = "cha_ching"
	SoundDogBarking       Sound = "dog_barking"
	SoundDoorBell         Sound = "door_bell"
	SoundDuckQuack        Sound = "duck_quack"
	SoundShortTripleBlink Sound = "short_triple_blink"
	SoundUpbeatBells      Sound = "upbeat_bells"
	SoundWarmSoftError    Sound = "warm_soft_error"
)

// InterruptionLevel controls how a notification is presented on the device.
type InterruptionLevel string

const (
	InterruptionPassive       InterruptionLevel = "passive"
	InterruptionActive        InterruptionLevel = "active"
	InterruptionTimeSensitive InterruptionLevel = "time-sensitive"
)

// Notification is the payload sent to the brrr.now API.
type Notification struct {
	Title             string            `json:"title,omitempty"`
	Subtitle          string            `json:"subtitle,omitempty"`
	Message           string            `json:"message,omitempty"`
	Sound             Sound             `json:"sound,omitempty"`
	OpenURL           string            `json:"open_url,omitempty"`
	ImageURL          string            `json:"image_url,omitempty"`
	ExpirationDate    *time.Time        `json:"expiration_date,omitempty"`
	FilterCriteria    string            `json:"filter-criteria,omitempty"`
	InterruptionLevel InterruptionLevel `json:"interruption-level,omitempty"`
}

// Client sends notifications through the brrr.now webhook API.
type Client struct {
	secret     string
	httpClient *http.Client
	logger     *slog.Logger
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithLogger sets a structured logger. By default no logging is performed.
func WithLogger(l *slog.Logger) Option {
	return func(cl *Client) { cl.logger = l }
}

// New creates a Client for the given webhook secret (e.g. "br_usr_a1b2c3d4e5f6g7h8i9j0").
func New(secret string, opts ...Option) *Client {
	c := &Client{
		secret:     secret,
		httpClient: http.DefaultClient,
		logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Send sends a fully customizable notification.
func (c *Client) Send(ctx context.Context, n Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("brrr: marshal notification: %w", err)
	}

	url := apiBaseURL + c.secret
	c.logger.DebugContext(ctx, "sending notification", "url", url, "payload_size", len(body))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brrr: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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

// SendMessage is a shortcut to send a plain text notification.
func (c *Client) SendMessage(ctx context.Context, message string) error {
	return c.Send(ctx, Notification{Message: message})
}

// SendWithTitle is a shortcut to send a notification with a title and message.
func (c *Client) SendWithTitle(ctx context.Context, title, message string) error {
	return c.Send(ctx, Notification{Title: title, Message: message})
}
