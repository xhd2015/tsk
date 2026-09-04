package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Event is one append-only events.jsonl record.
type Event struct {
	TS       string     `json:"ts"`
	Command  string     `json:"command"`
	Action   string     `json:"action,omitempty"`
	Args     []string   `json:"args"`
	ExitCode int        `json:"exit_code"`
	User     string     `json:"user,omitempty"`
	Mutation bool       `json:"mutation"`
	Data     *EventData `json:"data,omitempty"`
}

// EventData holds command-specific resolved operands.
type EventData struct {
	TaskID    int      `json:"task_id,omitempty"`
	ParentID  int      `json:"parent_id,omitempty"`
	Topic     string   `json:"topic,omitempty"`
	ChannelID string   `json:"channel_id,omitempty"`
	Project   string   `json:"project,omitempty"`
	Title     string   `json:"title,omitempty"`
	Text      string   `json:"text,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Label     string   `json:"label,omitempty"`
	Stage     string   `json:"stage,omitempty"`
	Status    string   `json:"status,omitempty"`
	Index     int      `json:"index,omitempty"`
	Name      string   `json:"name,omitempty"`
	Handle    string   `json:"handle,omitempty"`
	Alias     string   `json:"alias,omitempty"`
	Query     string   `json:"query,omitempty"`
	Notes     []string `json:"notes,omitempty"`
	MessageID int      `json:"message_id,omitempty"`
}

func eventsPath(home string) string {
	return filepath.Join(home, "events.jsonl")
}

// AppendEvent appends one JSON line to events.jsonl.
func AppendEvent(home string, ev Event) error {
	if err := os.MkdirAll(home, 0o755); err != nil {
		return fmt.Errorf("create tsk home: %w", err)
	}
	if ev.Args == nil {
		ev.Args = []string{}
	}
	line, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(eventsPath(home), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(append(line, '\n'))
	return err
}

// ReadEvents returns events.jsonl in file order (oldest first).
// Missing file yields an empty slice.
func ReadEvents(home string) ([]Event, error) {
	data, err := os.ReadFile(eventsPath(home))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read events.jsonl: %w", err)
	}
	var out []Event
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var ev Event
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			return nil, fmt.Errorf("parse events.jsonl: %w", err)
		}
		if ev.Args == nil {
			ev.Args = []string{}
		}
		out = append(out, ev)
	}
	return out, nil
}

// ProjectRefKey is origin if set, else registered name.
func ProjectRefKey(ref *ProjectRef) string {
	if ref == nil {
		return ""
	}
	if ref.Origin != "" {
		return ref.Origin
	}
	return ref.Name
}
