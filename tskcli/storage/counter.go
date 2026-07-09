package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

func counterPath(home string) string {
	return filepath.Join(home, "counter")
}

// NextID allocates the next monotonic task ID using flock on the counter file.
func NextID(home string) (int, error) {
	if err := EnsureLayout(home); err != nil {
		return 0, err
	}
	path := counterPath(home)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return 0, fmt.Errorf("open counter: %w", err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("flock counter: %w", err)
	}
	defer func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN) }()

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read counter: %w", err)
	}
	cur := 0
	if len(data) > 0 {
		cur, err = strconv.Atoi(string(data))
		if err != nil {
			return 0, fmt.Errorf("parse counter: %w", err)
		}
	}
	next := cur + 1
	if err := os.WriteFile(path, []byte(strconv.Itoa(next)), 0o644); err != nil {
		return 0, fmt.Errorf("write counter: %w", err)
	}
	return next, nil
}