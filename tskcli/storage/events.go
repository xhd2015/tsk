package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Event is one append-only events.jsonl record.
type Event struct {
	TS       string   `json:"ts"`
	Command  string   `json:"command"`
	Args     []string `json:"args"`
	ExitCode int      `json:"exit_code"`
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