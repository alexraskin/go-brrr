# go-brrr

A Go SDK for [brrr.now](https://brrr.now) webhooks.

## Install

```bash
go get github.com/alexraskin/go-brrr
```

## Usage

```go
package main

import (
	"context"
	"log"

	brrr "github.com/alexraskin/go-brrr"
)

func main() {
	client := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")

	err := client.SendMessage(context.Background(), "Hello from Go!")
	if err != nil {
		log.Fatal(err)
	}
}
```

### With a title

```go
client.SendWithTitle(ctx, "Deploy", "v1.2.3 is live")
```

### Full notification

```go
exp := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
client.Send(ctx, brrr.Notification{
	Title:             "Coffee Machine Offline",
	Message:           "Morale is expected to drop.",
	Sound:             brrr.SoundUpbeatBells,
	InterruptionLevel: brrr.InterruptionTimeSensitive,
	ExpirationDate:    &exp,
})
```

### Multiple webhooks

brrr.now gives you a shared webhook (all devices) and per-device webhooks. Create a client for each:

```go
all    := brrr.New("br_usr_shared_secret")
iphone := brrr.New("br_dev_iphone_secret")
mac    := brrr.New("br_dev_mac_secret")

// Send to all devices
all.SendMessage(ctx, "Hello everyone!")

// Send to a specific device
iphone.SendMessage(ctx, "Just your phone")
```

### With slog

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
client := brrr.New("br_usr_...", brrr.WithLogger(logger))
```

## Sounds

`SoundDefault`, `SoundSystem`, `SoundBrrr`, `SoundBellRinging`, `SoundBubbleDing`, `SoundBubblySuccessDng`, `SoundCatMeow`, `SoundCalm1`, `SoundCalm2`, `SoundChaChing`, `SoundDogBarking`, `SoundDoorBell`, `SoundDuckQuack`, `SoundShortTripleBlink`, `SoundUpbeatBells`, `SoundWarmSoftError`

## License

unlicense
