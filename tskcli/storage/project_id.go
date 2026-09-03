package storage

import (
	"fmt"
	"path/filepath"
	"strconv"
)

// ProjectNotesDir returns TSK_HOME/projects/<id> for project-scoped notes.
func ProjectNotesDir(home string, id int) string {
	return filepath.Join(home, "projects", strconv.Itoa(id))
}

// LookupProjectID finds a shared project id for ref in auto or registered
// ledgers. ok is false when no row matches. Matching rows are expected to
// already have id > 0.
func LookupProjectID(home string, ref ProjectRef) (id int, ok bool, err error) {
	auto, err := ReadProjectsAuto(home)
	if err != nil {
		return 0, false, err
	}
	if e, found := FindProjectAuto(auto, ref); found {
		if e.ID <= 0 {
			return 0, false, fmt.Errorf("project %s has no id", ProjectAutoKey(e.Origin, e.Name))
		}
		return e.ID, true, nil
	}

	reg, err := ReadProjects(home)
	if err != nil {
		return 0, false, err
	}
	if e, found := findRegisteredByRef(reg, ref); found {
		if e.ID <= 0 {
			return 0, false, fmt.Errorf("project %s has no id", ProjectAutoKey(e.Origin, e.Name))
		}
		return e.ID, true, nil
	}
	return 0, false, nil
}

func findRegisteredByRef(reg ProjectsFile, ref ProjectRef) (ProjectEntry, bool) {
	if ref.Origin != "" {
		return FindProjectByOrigin(reg, ref.Origin)
	}
	if ref.Name != "" {
		return FindProjectByName(reg, ref.Name)
	}
	return ProjectEntry{}, false
}

// AllocateSharedProjectID returns an existing id for the same ProjectAutoKey
// in either ledger, or allocates a new id from project-counter.
func AllocateSharedProjectID(home string, ref ProjectRef) (int, error) {
	key := ProjectAutoKey(ref.Origin, ref.Name)
	if key == "" {
		return 0, fmt.Errorf("empty project ref")
	}

	auto, err := ReadProjectsAuto(home)
	if err != nil {
		return 0, err
	}
	for _, e := range auto.Projects {
		if ProjectAutoKey(e.Origin, e.Name) == key && e.ID > 0 {
			return e.ID, nil
		}
	}

	reg, err := ReadProjects(home)
	if err != nil {
		return 0, err
	}
	for _, e := range reg.Projects {
		if ProjectAutoKey(e.Origin, e.Name) == key && e.ID > 0 {
			return e.ID, nil
		}
	}

	return NextProjectID(home)
}

// EnsureProjectID returns the shared id for ref, creating an auto ledger row
// when the project is not yet recorded.
func EnsureProjectID(home string, ref ProjectRef, locationTilde string) (int, error) {
	if id, ok, err := LookupProjectID(home, ref); err != nil {
		return 0, err
	} else if ok {
		return id, nil
	}
	if err := UpsertProjectAuto(home, ref, locationTilde); err != nil {
		return 0, err
	}
	id, ok, err := LookupProjectID(home, ref)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("project id missing after upsert")
	}
	return id, nil
}
