package main

import (
	"fmt"
	"time"
)

type activityResult struct {
	ChannelID      string
	LastActivityAt time.Time
	LastActivity   string
	Idle           time.Duration
	IdleStr        string
	IsIdle         bool
}

func loadChannelActivity(home, channelID string) (*activityResult, error) {
	status, err := channelIndexStatus(home, channelID)
	if err != nil {
		return nil, err
	}
	if status == "" {
		return nil, fmt.Errorf("channel %q not found", channelID)
	}
	if status == "archive" {
		return nil, fmt.Errorf("channel %q is archived", channelID)
	}

	dir := channelDir(home, channelID, status)
	createdAt, err := readChannelCreatedAt(dir)
	if err != nil {
		return nil, err
	}

	lastTS := createdAt
	if msgTS, ok, err := readLastMessageCreatedAt(dir); err != nil {
		return nil, err
	} else if ok {
		lastTS = msgTS
	}

	lastActivityAt, err := time.Parse(time.RFC3339, lastTS)
	if err != nil {
		return nil, fmt.Errorf("parse last activity %q: %w", lastTS, err)
	}

	now := time.Now().UTC()
	idle := now.Sub(lastActivityAt)
	if idle < 0 {
		idle = 0
	}

	return &activityResult{
		ChannelID:      channelID,
		LastActivityAt: lastActivityAt,
		LastActivity:   lastActivityAt.UTC().Format(time.RFC3339),
		Idle:           idle,
		IdleStr:        humanDuration(idle),
	}, nil
}

func (a *activityResult) checkIdle(threshold time.Duration) {
	a.IsIdle = a.Idle >= threshold
}

func humanDuration(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Second {
		return "0s"
	}
	return d.String()
}

func printStatus(a *activityResult, status string) {
	fmt.Printf("channel: %s\n", a.ChannelID)
	fmt.Printf("last_activity: %s\n", a.LastActivity)
	fmt.Printf("idle: %s\n", a.IdleStr)
	fmt.Printf("status: %s\n", status)
}