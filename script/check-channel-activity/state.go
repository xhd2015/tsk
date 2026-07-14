package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type notifyState struct {
	ChannelID      string `json:"channel_id"`
	LastActivityAt string `json:"last_activity_at"`
	LastNotifiedAt string `json:"last_notified_at"`
}

func stateFilePath(stateDir, channelID string) string {
	return filepath.Join(stateDir, channelID+".json")
}

func loadState(stateDir, channelID string) (*notifyState, error) {
	path := stateFilePath(stateDir, channelID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read state: %w", err)
	}
	var st notifyState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	return &st, nil
}

func writeState(stateDir, channelID, lastActivityAt string) error {
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	st := notifyState{
		ChannelID:      channelID,
		LastActivityAt: lastActivityAt,
		LastNotifiedAt: time.Now().UTC().Format(time.RFC3339),
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	data = append(data, '\n')
	path := stateFilePath(stateDir, channelID)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename state: %w", err)
	}
	return nil
}

func alreadyNotified(st *notifyState, lastActivityAt string) bool {
	if st == nil {
		return false
	}
	return st.LastActivityAt == lastActivityAt && st.LastNotifiedAt != ""
}