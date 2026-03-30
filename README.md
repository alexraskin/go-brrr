# go-brrr

A Go SDK for [brrr.now](https://brrr.now) webhooks.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/alexraskin/go-brrr) 
[![Go Reference](https://pkg.go.dev/badge/github.com/alexraskin/go-brrr.svg)](https://pkg.go.dev/github.com/alexraskin/go-brrr)

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
	client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendMessage(context.Background(), "Hello from Go!")
	if err != nil {
		log.Fatal(err)
	}
}
```

### With a title

```go
client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")
if err != nil {
	log.Fatal(err)
}

err = client.SendWithTitle(ctx, "Deploy", "v1.2.3 is live")
if err != nil {
	log.Fatal(err)
}
```

### Full notification

```go
client, err := brrr.New("br_usr_a1b2c3d4e5f6g7h8i9j0")
if err != nil {
	log.Fatal(err)
}

exp := time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
err = client.Send(ctx, brrr.Notification{
	Title:             "Coffee Machine Offline",
	Message:           "Morale is expected to drop.",
	Sound:             brrr.SoundUpbeatBells,
	InterruptionLevel: brrr.InterruptionTimeSensitive,
	ExpirationDate:    &exp,
})
if err != nil {
	log.Fatal(err)
}
```

### Multiple webhooks

brrr.now gives you a shared webhook (all devices) and per-device webhooks. Create a client for each:

```go
all, err := brrr.New("br_usr_shared_secret")
if err != nil {
	log.Fatal(err)
}

iphone, err := brrr.New("br_dev_iphone_secret")
if err != nil {
	log.Fatal(err)
}

mac, err := brrr.New("br_dev_mac_secret")
if err != nil {
	log.Fatal(err)
}

// Send to all devices
all.SendMessage(ctx, "Hello everyone!")

// Send to a specific device
iphone.SendMessage(ctx, "Just your phone")
```

### With slog

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

client, err := brrr.New("br_usr_...", brrr.WithLogger(logger))
if err != nil {
	log.Fatal(err)
}
```

## Sounds
> [!NOTE]  
> Only supported on iPhone and iPad.

`SoundDefault`, `SoundSystem`, `SoundBrrr`, `SoundBellRinging`, `SoundBubbleDing`, `SoundBubblySuccessDing`, `SoundCatMeow`, `SoundCalm1`, `SoundCalm2`, `SoundChaChing`, `SoundDogBarking`, `SoundDoorBell`, `SoundDuckQuack`, `SoundShortTripleBlink`, `SoundUpbeatBells`, `SoundWarmSoftError` 

## License

unlicense