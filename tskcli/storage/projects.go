package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ProjectEntry is one row in TSK_HOME/projects.json.
type ProjectEntry struct {
	Origin   string `json:"origin,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"` // main checkout, tilde form
	// Cwd is legacy-only: accepted on read, migrated into Location, never written.
	Cwd string `json:"cwd,omitempty"`
}

// ProjectsFile is the on-disk registry document.
type ProjectsFile struct {
	Projects []ProjectEntry `json:"projects"`
}

func projectsPath(home string) string {
	return filepath.Join(home, "projects.json")
}

// ReadProjects loads projects.json (empty list if missing).
func ReadProjects(home string) (ProjectsFile, error) {
	data, err := os.ReadFile(projectsPath(home))
	if err != nil {
		if os.IsNotExist(err) {
			return ProjectsFile{Projects: []ProjectEntry{}}, nil
		}
		return ProjectsFile{}, fmt.Errorf("read projects.json: %w", err)
	}
	var f ProjectsFile
	if err := json.Unmarshal(data, &f); err != nil {
		return ProjectsFile{}, fmt.Errorf("parse projects.json: %w", err)
	}
	if f.Projects == nil {
		f.Projects = []ProjectEntry{}
	}
	for i := range f.Projects {
		normalizeProjectEntry(&f.Projects[i])
	}
	return f, nil
}

// WriteProjects writes projects.json atomically (location only; no cwd).
func WriteProjects(home string, f ProjectsFile) error {
	if f.Projects == nil {
		f.Projects = []ProjectEntry{}
	}
	for i := range f.Projects {
		normalizeProjectEntry(&f.Projects[i])
	}
	if err := os.MkdirAll(home, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(home, "projects-*.json.tmp")
	if err != nil {
		return fmt.Errorf("create projects temp: %w", err)
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
	dst := projectsPath(home)
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename projects.json: %w", err)
	}
	return nil
}

// FindProjectByName returns the entry with the given unique name.
func FindProjectByName(f ProjectsFile, name string) (ProjectEntry, bool) {
	name = strings.TrimSpace(name)
	for _, e := range f.Projects {
		if e.Name == name {
			return e, true
		}
	}
	return ProjectEntry{}, false
}

// FindProjectByOrigin returns the first entry with the given origin key.
func FindProjectByOrigin(f ProjectsFile, origin string) (ProjectEntry, bool) {
	origin = strings.TrimSpace(origin)
	for _, e := range f.Projects {
		if e.Origin == origin && e.Origin != "" {
			return e, true
		}
	}
	return ProjectEntry{}, false
}

// EffectiveLocation returns the project location (after normalize).
func (e ProjectEntry) EffectiveLocation() string {
	return e.Location
}

func normalizeProjectEntry(e *ProjectEntry) {
	if e.Location == "" && e.Cwd != "" {
		e.Location = e.Cwd
	}
	e.Cwd = ""
}

// SortedProjects returns a copy sorted by name then origin.
func SortedProjects(f ProjectsFile) []ProjectEntry {
	out := append([]ProjectEntry(nil), f.Projects...)
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
