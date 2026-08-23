package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TopicNotesFile is the append-only journal in a topic directory.
const TopicNotesFile = "notes.jsonl"

// TopicNote is one timestamped note line (topics and tasks).
type TopicNote struct {
	TS     string   `json:"ts"`
	Text   string   `json:"text"`
	Labels []string `json:"labels,omitempty"`
	Status string   `json:"status,omitempty"`
}

// NoteHasAllLabels reports whether n has every label in want (AND).
// Empty want matches every note.
func NoteHasAllLabels(n TopicNote, want []string) bool {
	if len(want) == 0 {
		return true
	}
	have := make(map[string]struct{}, len(n.Labels))
	for _, l := range n.Labels {
		have[l] = struct{}{}
	}
	for _, l := range want {
		if _, ok := have[l]; !ok {
			return false
		}
	}
	return true
}

// FilterNotes keeps notes that have every label in labels (AND).
func FilterNotes(notes []TopicNote, labels []string) []TopicNote {
	if len(labels) == 0 {
		return notes
	}
	var out []TopicNote
	for _, n := range notes {
		if NoteHasAllLabels(n, labels) {
			out = append(out, n)
		}
	}
	return out
}

// ApplyNoteLimit returns the last limit notes when limit > 0.
func ApplyNoteLimit(notes []TopicNote, limit int) []TopicNote {
	if limit > 0 && len(notes) > limit {
		return notes[len(notes)-limit:]
	}
	return notes
}

func topicNotesPath(topicDir string) string {
	return filepath.Join(topicDir, TopicNotesFile)
}

// ReadTopicNotes returns all notes in file order (oldest first).
func ReadTopicNotes(topicDir string) ([]TopicNote, error) {
	f, err := os.Open(topicNotesPath(topicDir))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []TopicNote
	sc := bufio.NewScanner(f)
	// notes can be long; raise the token size a bit
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var n TopicNote
		if err := json.Unmarshal([]byte(line), &n); err != nil {
			return nil, fmt.Errorf("parse %s line %d: %w", TopicNotesFile, lineNo, err)
		}
		out = append(out, n)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// AppendTopicNote appends one JSON line to notes.jsonl.
func AppendTopicNote(topicDir string, note TopicNote) error {
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		return err
	}
	line, err := json.Marshal(note)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(topicNotesPath(topicDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", TopicNotesFile, err)
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("write %s: %w", TopicNotesFile, err)
	}
	return nil
}

// RewriteNotes atomically replaces notes.jsonl with the given notes
// (one JSON object per line). If notes is empty, the file is truncated.
func RewriteNotes(topicDir string, notes []TopicNote) error {
	if err := os.MkdirAll(topicDir, 0o755); err != nil {
		return err
	}
	path := topicNotesPath(topicDir)
	tmp, err := os.CreateTemp(topicDir, "notes-*.jsonl.tmp")
	if err != nil {
		return fmt.Errorf("create temp: %w", err)
	}
	tmpName := tmp.Name()
	for _, n := range notes {
		line, err := json.Marshal(n)
		if err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
			return fmt.Errorf("marshal note: %w", err)
		}
		if _, err := tmp.Write(append(line, '\n')); err != nil {
			_ = tmp.Close()
			_ = os.Remove(tmpName)
			return fmt.Errorf("write temp: %w", err)
		}
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close temp: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename %s: %w", TopicNotesFile, err)
	}
	return nil
}

// MigrateLegacyTopicNotes copies topic.json "notes" into notes.jsonl once, then
// clears the blob so it is not a second source of truth.
func MigrateLegacyTopicNotes(topicDir string, parts []string) error {
	path := topicNotesPath(topicDir)
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	meta, err := ReadTopicMeta(topicDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	text := strings.TrimSpace(meta.Notes)
	if text == "" {
		return nil
	}
	ts := meta.UpdatedAt
	if ts == "" {
		ts = NowTimestamp(1)
	}
	if err := AppendTopicNote(topicDir, TopicNote{TS: ts, Text: text}); err != nil {
		return err
	}
	meta.Notes = ""
	meta.Path = append([]string(nil), parts...)
	return WriteTopicMeta(topicDir, meta)
}
