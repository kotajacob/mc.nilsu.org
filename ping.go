package main

import (
	"encoding/json"
	"fmt"

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

func ping(addr string) (*status, error) {
	var s status
	resp, _, err := bot.PingAndList(addr)
	if err != nil {
		return &s, fmt.Errorf("ping failed: %v\n", err)
	}

	err = json.Unmarshal(resp, &s)
	if err != nil {
		return &s, fmt.Errorf("error unmarshaling response: %v\n", err)
	}
	return &s, nil
}
