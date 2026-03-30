package brrr_test

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	brrr "github.com/alexraskin/go-brrr"
)

func Example() {
	client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SendMessage(context.Background(), "Hello world!")
	if err != nil {
		fmt.Println(err)
	}
}

func Example_withLogger() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0", brrr.WithLogger(logger))
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.SendWithTitle(context.Background(), "Alert", "Something happened!")
	if err != nil {
		fmt.Println(err)
	}
}

func Example_customNotification() {
	client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")
	if err != nil {
		fmt.Println(err)
		return
	}

	exp := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
	err = client.Send(context.Background(), brrr.Notification{
		Title:             "Coffee Machine Offline",
		Message:           "The coffee machine is currently unreachable.",
		Sound:             brrr.SoundUpbeatBells,
		ExpirationDate:    &exp,
		InterruptionLevel: brrr.InterruptionTimeSensitive,
	})
	if err != nil {
		fmt.Println(err)
	}
}
