package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// ResolveHome returns the TSK_HOME directory.
func ResolveHome() (string, error) {
	if v := os.Getenv("TSK_HOME"); v != "" {
		return filepath.Abs(v)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".tsk"), nil
}

// EnsureLayout creates the base storage directories under home.
func EnsureLayout(home string) error {
	for _, dir := range []string{"index", "inbox", "topics"} {
		if err := os.MkdirAll(filepath.Join(home, dir), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return nil
}