package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func projectCounterPath(home string) string {
	return filepath.Join(home, "project-counter")
}

// NextProjectID allocates the next monotonic project ID using flock on
// project-counter. Shared by projects.json and projects-auto.json.
func NextProjectID(home string) (int, error) {
	if err := EnsureLayout(home); err != nil {
		return 0, err
	}
	path := projectCounterPath(home)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return 0, fmt.Errorf("open project-counter: %w", err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("flock project-counter: %w", err)
	}
	defer func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN) }()

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read project-counter: %w", err)
	}
	cur := 0
	if len(data) > 0 {
		cur, err = strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			return 0, fmt.Errorf("parse project-counter: %w", err)
		}
	}
	next := cur + 1
	if err := os.WriteFile(path, []byte(strconv.Itoa(next)), 0o644); err != nil {
		return 0, fmt.Errorf("write project-counter: %w", err)
	}
	return next, nil
}
