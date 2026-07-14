package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type channelJSON struct {
	CreatedAt string `json:"created_at"`
}

type channelMessage struct {
	CreatedAt string `json:"created_at"`
}

func resolveHome(tskHome string) (string, error) {
	if tskHome != "" {
		return filepath.Abs(tskHome)
	}
	if v := os.Getenv("TSK_HOME"); v != "" {
		return filepath.Abs(v)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".tsk"), nil
}

func channelIndexStatus(home, channelID string) (string, error) {
	path := filepath.Join(home, "channels", "index", channelID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	line := strings.TrimSpace(string(data))
	switch {
	case strings.HasPrefix(line, "active/"):
		return "active", nil
	case strings.HasPrefix(line, "archive/"):
		return "archive", nil
	default:
		return "", fmt.Errorf("invalid index entry for channel %q: %q", channelID, line)
	}
}

func channelDir(home, channelID, status string) string {
	return filepath.Join(home, "channels", status, channelID)
}

func readChannelCreatedAt(channelDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(channelDir, "channel.json"))
	if err != nil {
		return "", err
	}
	var ch channelJSON
	if err := json.Unmarshal(data, &ch); err != nil {
		return "", fmt.Errorf("parse channel.json: %w", err)
	}
	return ch.CreatedAt, nil
}

func readLastMessageCreatedAt(channelDir string) (string, bool, error) {
	path := filepath.Join(channelDir, "messages.jsonl")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	defer f.Close()

	var last string
	found := false
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var msg channelMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			return "", false, fmt.Errorf("parse messages.jsonl: %w", err)
		}
		last = msg.CreatedAt
		found = true
	}
	if err := sc.Err(); err != nil {
		return "", false, err
	}
	return last, found, nil
}