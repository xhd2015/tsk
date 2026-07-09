package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"
)

// StageHistoryEntry records one stage change.
type StageHistoryEntry struct {
	From string `json:"from"`
	To   string `json:"to"`
	At   string `json:"at"`
	Note string `json:"note,omitempty"`
}

// Task is the on-disk task.json schema.
type Task struct {
	ID           int                 `json:"id"`
	Title        string              `json:"title"`
	Slug         string              `json:"slug"`
	Labels       []string            `json:"labels"`
	TopicPath    json.RawMessage     `json:"topic_path"`
	Stage        string              `json:"stage"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
	StageHistory []StageHistoryEntry   `json:"stage_history"`
}

// ClarifyItem is one question in a clarify batch.
type ClarifyItem struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Status   string `json:"status"`
}

// ClarifyBatch is clarify/batch.json.
type ClarifyBatch struct {
	BatchID string        `json:"batch_id"`
	Status  string        `json:"status"`
	Items   []ClarifyItem `json:"items"`
}

// Slugify converts a title into a path-safe slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	s = b.String()
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}
	s = strings.Trim(s, "-")
	runes := []rune(s)
	if len(runes) > 64 {
		s = string(runes[:64])
		s = strings.Trim(s, "-")
	}
	return s
}

// TaskDirName returns the directory name for a task.
func TaskDirName(id int, stage, slug string) string {
	return fmt.Sprintf("%d-%s-%s", id, stage, slug)
}

// InboxRelPath returns the inbox-relative path for a task directory.
func InboxRelPath(id int, stage, title string) string {
	return filepath.ToSlash(filepath.Join("inbox", TaskDirName(id, stage, Slugify(title))))
}

// TopicRelPath returns the topic-relative path for a task directory.
func TopicRelPath(topic string, id int, stage, title string) string {
	parts := strings.Split(topic, "/")
	all := append(parts, TaskDirName(id, stage, Slugify(title)))
	return filepath.ToSlash(filepath.Join(append([]string{"topics"}, all...)...))
}

// NullTopicPath is the JSON encoding for inbox tasks.
var NullTopicPath = json.RawMessage("null")

// TopicPathJSON encodes a topic path slice for task.json.
func TopicPathJSON(parts []string) (json.RawMessage, error) {
	if len(parts) == 0 {
		return NullTopicPath, nil
	}
	data, err := json.Marshal(parts)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ParseTopicPath decodes topic_path from task.json.
func ParseTopicPath(raw json.RawMessage) ([]string, error) {
	s := strings.TrimSpace(string(raw))
	if s == "null" || s == "" {
		return nil, nil
	}
	var parts []string
	if err := json.Unmarshal(raw, &parts); err != nil {
		return nil, err
	}
	return parts, nil
}

// NowTimestamp returns a deterministic or current RFC3339 timestamp.
func NowTimestamp(seq int) string {
	if date := os.Getenv("TSK_DATE"); date != "" {
		return fmt.Sprintf("%sT%02d:00:00Z", date, seq%24)
	}
	return time.Now().UTC().Format(time.RFC3339)
}

// ReadTask loads task.json from a task directory.
func ReadTask(taskDir string) (Task, error) {
	data, err := os.ReadFile(filepath.Join(taskDir, "task.json"))
	if err != nil {
		return Task{}, fmt.Errorf("read task.json: %w", err)
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return Task{}, fmt.Errorf("parse task.json: %w", err)
	}
	return task, nil
}

// WriteTask writes task.json atomically.
func WriteTask(taskDir string, task Task) error {
	if task.Labels == nil {
		task.Labels = []string{}
	}
	if task.StageHistory == nil {
		task.StageHistory = []StageHistoryEntry{}
	}
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(taskDir, "task-*.json.tmp")
	if err != nil {
		return fmt.Errorf("create task temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write task temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	dst := filepath.Join(taskDir, "task.json")
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename task.json: %w", err)
	}
	return nil
}

// RenameTaskDir renames a task directory and updates the index.
func RenameTaskDir(home string, task *Task, oldDir, newStage string, note string) (string, error) {
	topicParts, err := ParseTopicPath(task.TopicPath)
	if err != nil {
		return "", err
	}
	var newRel string
	if len(topicParts) == 0 {
		newRel = InboxRelPath(task.ID, newStage, task.Title)
	} else {
		newRel = TopicRelPath(strings.Join(topicParts, "/"), task.ID, newStage, task.Title)
	}
	newAbs := filepath.Join(home, filepath.FromSlash(newRel))
	if err := os.MkdirAll(filepath.Dir(newAbs), 0o755); err != nil {
		return "", fmt.Errorf("mkdir parent: %w", err)
	}
	if err := os.Rename(oldDir, newAbs); err != nil {
		return "", fmt.Errorf("rename task dir: %w", err)
	}
	now := NowTimestamp(task.ID)
	from := task.Stage
	task.Stage = newStage
	task.UpdatedAt = now
	entry := StageHistoryEntry{From: from, To: newStage, At: now}
	if note != "" {
		entry.Note = note
	}
	task.StageHistory = append(task.StageHistory, entry)
	if err := WriteTask(newAbs, *task); err != nil {
		return "", err
	}
	if err := WriteIndex(home, task.ID, newRel); err != nil {
		return "", err
	}
	return newAbs, nil
}

// MoveTaskDir moves a task directory to a new location and updates index/topic_path.
func MoveTaskDir(home string, task *Task, oldDir string, topicParts []string) (string, error) {
	var newRel string
	if len(topicParts) == 0 {
		newRel = InboxRelPath(task.ID, task.Stage, task.Title)
	} else {
		newRel = TopicRelPath(strings.Join(topicParts, "/"), task.ID, task.Stage, task.Title)
	}
	newAbs := filepath.Join(home, filepath.FromSlash(newRel))
	if err := os.MkdirAll(filepath.Dir(newAbs), 0o755); err != nil {
		return "", fmt.Errorf("mkdir parent: %w", err)
	}
	if err := os.Rename(oldDir, newAbs); err != nil {
		return "", fmt.Errorf("move task dir: %w", err)
	}
	topicJSON, err := TopicPathJSON(topicParts)
	if err != nil {
		return "", err
	}
	task.TopicPath = topicJSON
	task.UpdatedAt = NowTimestamp(task.ID)
	if err := WriteTask(newAbs, *task); err != nil {
		return "", err
	}
	if err := WriteIndex(home, task.ID, newRel); err != nil {
		return "", err
	}
	return newAbs, nil
}

// ListTaskIDs returns all task IDs from the index directory.
func ListTaskIDs(home string) ([]int, error) {
	entries, err := os.ReadDir(indexDir(home))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []int
	for _, ent := range entries {
		var id int
		if _, err := fmt.Sscanf(ent.Name(), "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}
	sort.Ints(ids)
	return ids, nil
}

// LoadTaskByID loads a task by ID.
func LoadTaskByID(home string, id int) (Task, string, error) {
	dir, err := TaskDir(home, id)
	if err != nil {
		return Task{}, "", err
	}
	task, err := ReadTask(dir)
	if err != nil {
		return Task{}, "", err
	}
	return task, dir, nil
}

// ReadClarifyBatch loads clarify/batch.json.
func ReadClarifyBatch(taskDir string) (ClarifyBatch, error) {
	data, err := os.ReadFile(filepath.Join(taskDir, "clarify", "batch.json"))
	if err != nil {
		return ClarifyBatch{}, err
	}
	var batch ClarifyBatch
	if err := json.Unmarshal(data, &batch); err != nil {
		return ClarifyBatch{}, err
	}
	return batch, nil
}

// WriteClarifyBatch writes clarify/batch.json.
func WriteClarifyBatch(taskDir string, batch ClarifyBatch) error {
	clarifyDir := filepath.Join(taskDir, "clarify")
	if err := os.MkdirAll(clarifyDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(batch, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(clarifyDir, "batch.json"), data, 0o644)
}

// EnsureClarifyBatch returns an existing or new open clarify batch.
func EnsureClarifyBatch(taskDir string) (ClarifyBatch, error) {
	batch, err := ReadClarifyBatch(taskDir)
	if err == nil {
		return batch, nil
	}
	if !os.IsNotExist(err) {
		return ClarifyBatch{}, err
	}
	return ClarifyBatch{
		BatchID: "b1",
		Status:  "open",
		Items:   []ClarifyItem{},
	}, nil
}