package brrr

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestSendMessage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected application/json, got %s", ct)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{
		secret:     "test_secret",
		httpClient: srv.Client(),
		logger:     slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})),
	}
	// Override base URL by pointing at the test server.
	origURL := apiBaseURL
	defer func() { setBaseURL(origURL) }()
	setBaseURL(srv.URL + "/")

	err := c.SendMessage(context.Background(), "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendFull(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := &Client{
		secret:     "test_secret",
		httpClient: srv.Client(),
		logger:     slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	origURL := apiBaseURL
	defer func() { setBaseURL(origURL) }()
	setBaseURL(srv.URL + "/")

	exp := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	err := c.Send(context.Background(), Notification{
		Title:             "Test",
		Subtitle:          "Sub",
		Message:           "Body",
		Sound:             SoundBrrr,
		OpenURL:           "https://example.com",
		ImageURL:          "https://example.com/img.png",
		ExpirationDate:    &exp,
		FilterCriteria:    "work",
		InterruptionLevel: InterruptionTimeSensitive,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad secret"))
	}))
	defer srv.Close()

	c := &Client{
		secret:     "bad",
		httpClient: srv.Client(),
		logger:     slog.New(slog.NewTextHandler(os.Stderr, nil)),
	}
	origURL := apiBaseURL
	defer func() { setBaseURL(origURL) }()
	setBaseURL(srv.URL + "/")

	err := c.SendMessage(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}

