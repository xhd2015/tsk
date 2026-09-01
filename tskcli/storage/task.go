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

// ProjectRef is the optional project identity on a task.
// Prefer Origin when the repo has a remote; Name is used only for
// registered non-git projects (origin XOR name on disk).
type ProjectRef struct {
	Origin string `json:"origin,omitempty"`
	Name   string `json:"name,omitempty"`
}

// Task is the on-disk task.json schema.
type Task struct {
	ID           int                 `json:"id"`
	Title        string              `json:"title"`
	Slug         string              `json:"slug"`
	Labels       []string            `json:"labels"`
	TopicPath    json.RawMessage     `json:"topic_path"`
	ParentID     int                 `json:"parent_id,omitempty"`
	Cwd          string              `json:"cwd,omitempty"`
	Project      *ProjectRef         `json:"project,omitempty"`
	Stage        string              `json:"stage"`
	CreatedAt    string              `json:"created_at"`
	UpdatedAt    string              `json:"updated_at"`
	StageHistory []StageHistoryEntry `json:"stage_history"`
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

// TaskDirName returns the directory name for a task: [id]-<stage>-<slug>.
func TaskDirName(id int, stage, slug string) string {
	return fmt.Sprintf("[%d]-%s-%s", id, stage, slug)
}

// InboxRelPath returns the inbox-relative path for a top-level task directory.
func InboxRelPath(id int, stage, title string) string {
	return filepath.ToSlash(filepath.Join("inbox", TaskDirName(id, stage, Slugify(title))))
}

// TopicRelPath returns the topic-relative path for a top-level task directory.
func TopicRelPath(topic string, id int, stage, title string) string {
	parts := strings.Split(topic, "/")
	all := append(parts, TaskDirName(id, stage, Slugify(title)))
	return filepath.ToSlash(filepath.Join(append([]string{"topics"}, all...)...))
}

// ChildRelPath returns the relative path for a nested sub-task under parentRel.
func ChildRelPath(parentRel string, id int, stage, slug string) string {
	return filepath.ToSlash(filepath.Join(parentRel, TaskDirName(id, stage, slug)))
}

// RelFromHome returns the slash-separated path of abs relative to home.
func RelFromHome(home, abs string) (string, error) {
	rel, err := filepath.Rel(home, abs)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

// RewriteIndexPrefix rewrites every index entry with path oldPrefix or under
// oldPrefix/ to the corresponding newPrefix path.
func RewriteIndexPrefix(home, oldPrefix, newPrefix string) error {
	oldPrefix = filepath.ToSlash(oldPrefix)
	newPrefix = filepath.ToSlash(newPrefix)
	ids, err := ListTaskIDs(home)
	if err != nil {
		return err
	}
	for _, id := range ids {
		rel, err := ReadIndex(home, id)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == oldPrefix {
			if err := WriteIndex(home, id, newPrefix); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(rel, oldPrefix+"/") {
			suffix := strings.TrimPrefix(rel, oldPrefix)
			if err := WriteIndex(home, id, newPrefix+suffix); err != nil {
				return err
			}
		}
	}
	return nil
}

// WalkTaskDirs visits taskDir and every nested task directory (skips context/clarify).
func WalkTaskDirs(taskDir string, fn func(abs string, task Task) error) error {
	task, err := ReadTask(taskDir)
	if err != nil {
		return err
	}
	if err := fn(taskDir, task); err != nil {
		return err
	}
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return err
	}
	for _, ent := range entries {
		if !ent.IsDir() || !IsTaskDirName(ent.Name()) {
			continue
		}
		if err := WalkTaskDirs(filepath.Join(taskDir, ent.Name()), fn); err != nil {
			return err
		}
	}
	return nil
}

// TopLevelTaskRel builds inbox/ or topics/… relative path for a root-level task.
func TopLevelTaskRel(topicParts []string, id int, stage, slug string) string {
	name := TaskDirName(id, stage, slug)
	if len(topicParts) == 0 {
		return filepath.ToSlash(filepath.Join("inbox", name))
	}
	elems := append([]string{"topics"}, topicParts...)
	elems = append(elems, name)
	return filepath.ToSlash(filepath.Join(elems...))
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

// RenameTaskDir renames a task directory in place (path-preserving) and
// cascades index updates for nested descendants.
func RenameTaskDir(home string, task *Task, oldDir, newStage string, note string) (string, error) {
	oldRel, err := RelFromHome(home, oldDir)
	if err != nil {
		return "", fmt.Errorf("rel task dir: %w", err)
	}
	parentRel := filepath.ToSlash(filepath.Dir(oldRel))
	newName := TaskDirName(task.ID, newStage, task.Slug)
	var newRel string
	if parentRel == "." || parentRel == "" {
		newRel = newName
	} else {
		newRel = filepath.ToSlash(filepath.Join(parentRel, newName))
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
	if err := RewriteIndexPrefix(home, oldRel, newRel); err != nil {
		return "", err
	}
	return newAbs, nil
}

// MoveTaskDir moves a root-level task directory (and nested children) to a new
// topic or inbox location and updates index/topic_path for the whole subtree.
func MoveTaskDir(home string, task *Task, oldDir string, topicParts []string) (string, error) {
	if task.ParentID != 0 {
		return "", fmt.Errorf("nested task %d: reparent before changing topic", task.ID)
	}
	oldRel, err := RelFromHome(home, oldDir)
	if err != nil {
		return "", fmt.Errorf("rel task dir: %w", err)
	}
	newRel := TopLevelTaskRel(topicParts, task.ID, task.Stage, task.Slug)
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
	now := NowTimestamp(task.ID)
	task.TopicPath = topicJSON
	task.UpdatedAt = now
	if err := WriteTask(newAbs, *task); err != nil {
		return "", err
	}
	if err := WalkTaskDirs(newAbs, func(abs string, t Task) error {
		if abs == newAbs {
			return nil
		}
		t.TopicPath = topicJSON
		t.UpdatedAt = now
		return WriteTask(abs, t)
	}); err != nil {
		return "", err
	}
	if err := RewriteIndexPrefix(home, oldRel, newRel); err != nil {
		return "", err
	}
	return newAbs, nil
}

// DeleteTask permanently removes a task directory and its index entry.
// Nested sub-tasks require recursive=true (whole subtree). Without recursive,
// presence of any immediate nested task directory is an error.
// Task IDs are never reused (monotonic counter), so no tombstone is written.
func DeleteTask(home string, id int, recursive bool) error {
	_, taskDir, err := LoadTaskByID(home, id)
	if err != nil {
		return err
	}

	var immediate []int
	entries, err := os.ReadDir(taskDir)
	if err != nil {
		return fmt.Errorf("read task dir: %w", err)
	}
	for _, ent := range entries {
		if !ent.IsDir() || !IsTaskDirName(ent.Name()) {
			continue
		}
		cid, _, _, ok := ParseTaskDirName(ent.Name())
		if !ok {
			continue
		}
		immediate = append(immediate, cid)
	}
	sort.Ints(immediate)
	if len(immediate) > 0 && !recursive {
		return fmt.Errorf("tsk delete: task %d has nested tasks %v; use --recursive", id, immediate)
	}

	var ids []int
	if err := WalkTaskDirs(taskDir, func(_ string, t Task) error {
		ids = append(ids, t.ID)
		return nil
	}); err != nil {
		return err
	}

	for _, rid := range ids {
		if err := os.Remove(indexPath(home, rid)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("remove index/%d: %w", rid, err)
		}
	}
	if err := os.RemoveAll(taskDir); err != nil {
		return fmt.Errorf("delete task dir: %w", err)
	}
	return nil
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
