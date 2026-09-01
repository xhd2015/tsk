package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ProjectAutoEntry is one row in TSK_HOME/projects-auto.json (add-/upsert-only).
type ProjectAutoEntry struct {
	Origin      string `json:"origin,omitempty"`
	Name        string `json:"name,omitempty"`
	Cwd         string `json:"cwd,omitempty"` // main repo, tilde form; set once
	FirstSeenAt string `json:"first_seen_at"`
	LastSeenAt  string `json:"last_seen_at"`
}

// ProjectsAutoFile is the on-disk auto-seen ledger.
type ProjectsAutoFile struct {
	Projects []ProjectAutoEntry `json:"projects"`
}

func projectsAutoPath(home string) string {
	return filepath.Join(home, "projects-auto.json")
}

// ReadProjectsAuto loads projects-auto.json (empty if missing).
func ReadProjectsAuto(home string) (ProjectsAutoFile, error) {
	data, err := os.ReadFile(projectsAutoPath(home))
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectsAutoFile{Projects: []ProjectAutoEntry{}}, nil
		}
		return ProjectsAutoFile{}, fmt.Errorf("read projects-auto.json: %w", err)
	}
	var f ProjectsAutoFile
	if err := json.Unmarshal(data, &f); err != nil {
		return ProjectsAutoFile{}, fmt.Errorf("parse projects-auto.json: %w", err)
	}
	if f.Projects == nil {
		f.Projects = []ProjectAutoEntry{}
	}
	return f, nil
}

// WriteProjectsAuto writes projects-auto.json atomically.
func WriteProjectsAuto(home string, f ProjectsAutoFile) error {
	if f.Projects == nil {
		f.Projects = []ProjectAutoEntry{}
	}
	if err := os.MkdirAll(home, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(home, "projects-auto-*.json.tmp")
	if err != nil {
		return fmt.Errorf("create projects-auto temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	dst := projectsAutoPath(home)
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename projects-auto.json: %w", err)
	}
	return nil
}

// NowLocalTimestamp returns RFC3339 with local timezone offset (e.g. +08:00).
// When TSK_DATE is set, returns a deterministic local-offset stamp for tests.
func NowLocalTimestamp() string {
	if date := os.Getenv("TSK_DATE"); date != "" {
		// Fixed +08:00 so doctests are stable across hosts.
		return date + "T12:00:00+08:00"
	}
	return time.Now().Format(time.RFC3339)
}

// UpsertProjectAuto inserts or updates an auto-seen project.
// New rows set cwd/first/last; existing rows only bump last_seen_at.
func UpsertProjectAuto(home string, ref ProjectRef, cwdTilde string) error {
	f, err := ReadProjectsAuto(home)
	if err != nil {
		return err
	}
	now := NowLocalTimestamp()
	for i := range f.Projects {
		e := &f.Projects[i]
		if autoEntryMatches(e, ref) {
			e.LastSeenAt = now
			return WriteProjectsAuto(home, f)
		}
	}
	entry := ProjectAutoEntry{
		Origin:      ref.Origin,
		Name:        ref.Name,
		Cwd:         cwdTilde,
		FirstSeenAt: now,
		LastSeenAt:  now,
	}
	f.Projects = append(f.Projects, entry)
	return WriteProjectsAuto(home, f)
}

func autoEntryMatches(e *ProjectAutoEntry, ref ProjectRef) bool {
	if ref.Origin != "" {
		return e.Origin == ref.Origin
	}
	if ref.Name != "" {
		return e.Name == ref.Name && e.Origin == ""
	}
	return false
}

// SortedProjectsAuto returns a copy sorted by name then origin.
func SortedProjectsAuto(f ProjectsAutoFile) []ProjectAutoEntry {
	out := append([]ProjectAutoEntry(nil), f.Projects...)
	sort.Slice(out, func(i, j int) bool {
		ni, nj := out[i].Name, out[j].Name
		if ni != nj {
			if ni == "" {
				return false
			}
			if nj == "" {
				return true
			}
			return ni < nj
		}
		return out[i].Origin < out[j].Origin
	})
	return out
}

// ProjectAutoKey returns a stable map key for an auto entry or task project ref.
func ProjectAutoKey(origin, name string) string {
	origin = strings.TrimSpace(origin)
	name = strings.TrimSpace(name)
	if origin != "" {
		return "origin:" + origin
	}
	if name != "" {
		return "name:" + name
	}
	return ""
}
