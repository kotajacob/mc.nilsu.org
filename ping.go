package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
}

// ping a minecraft server for up to 5 seconds to populate a status variable. If
// the ping fails or cannot be parsed an error is returned.
func ping(addr string) (*status, error) {
	var s status
	resp, _, err := bot.PingAndListTimeout(addr, time.Second*10)
	if err != nil {
		return &s, fmt.Errorf("ping failed: %v\n", err)
	}

	err = json.Unmarshal(resp, &s)
	if err != nil {
		return &s, fmt.Errorf("error unmarshaling response: %v\n", err)
	}
	return &s, nil
}
