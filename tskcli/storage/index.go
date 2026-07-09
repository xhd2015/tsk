package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func indexDir(home string) string {
	return filepath.Join(home, "index")
}

func indexPath(home string, id int) string {
	return filepath.Join(indexDir(home), strconv.Itoa(id))
}

// WriteIndex atomically writes index/<id> with the relative task path.
func WriteIndex(home string, id int, relPath string) error {
	if err := os.MkdirAll(indexDir(home), 0o755); err != nil {
		return fmt.Errorf("create index dir: %w", err)
	}
	relPath = filepath.ToSlash(relPath)
	tmp, err := os.CreateTemp(indexDir(home), fmt.Sprintf("%d-*.tmp", id))
	if err != nil {
		return fmt.Errorf("create index temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.WriteString(relPath); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write index temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("close index temp: %w", err)
	}
	if err := os.Rename(tmpName, indexPath(home, id)); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename index: %w", err)
	}
	return nil
}

// ReadIndex returns the relative path for a task ID.
func ReadIndex(home string, id int) (string, error) {
	data, err := os.ReadFile(indexPath(home, id))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("task %d not found", id)
		}
		return "", fmt.Errorf("read index/%d: %w", id, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// TaskDir returns the absolute path to a task directory by ID.
func TaskDir(home string, id int) (string, error) {
	rel, err := ReadIndex(home, id)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, filepath.FromSlash(rel)), nil
}