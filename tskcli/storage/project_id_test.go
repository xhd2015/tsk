package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNextProjectIDMonotonic(t *testing.T) {
	home := t.TempDir()
	a, err := NextProjectID(home)
	if err != nil {
		t.Fatal(err)
	}
	b, err := NextProjectID(home)
	if err != nil {
		t.Fatal(err)
	}
	if a != 1 || b != 2 {
		t.Fatalf("got %d %d want 1 2", a, b)
	}
	raw, err := os.ReadFile(filepath.Join(home, "project-counter"))
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "2" {
		t.Fatalf("counter=%q", raw)
	}
}

func TestAllocateSharedReusesAcrossLedgers(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TSK_DATE", "2026-07-09")

	ref := ProjectRef{Origin: "github.com/x/y"}
	if err := UpsertProjectAuto(home, ref, "~/y"); err != nil {
		t.Fatal(err)
	}
	auto, err := ReadProjectsAuto(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(auto.Projects) != 1 || auto.Projects[0].ID != 1 {
		t.Fatalf("auto=%+v", auto.Projects)
	}

	id, err := AllocateSharedProjectID(home, ref)
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("reuse id=%d want 1", id)
	}

	reg := ProjectsFile{Projects: []ProjectEntry{{
		ID: id, Origin: ref.Origin, Name: "y", Location: "~/y",
	}}}
	if err := WriteProjects(home, reg); err != nil {
		t.Fatal(err)
	}

	id2, err := AllocateSharedProjectID(home, ProjectRef{Origin: ref.Origin, Name: "y"})
	if err != nil {
		t.Fatal(err)
	}
	if id2 != 1 {
		t.Fatalf("still want 1, got %d", id2)
	}
}

func TestRegisterThenAutoSharesID(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TSK_DATE", "2026-07-09")

	ref := ProjectRef{Name: "seatalk"}
	id, err := AllocateSharedProjectID(home, ref)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteProjects(home, ProjectsFile{Projects: []ProjectEntry{{
		ID: id, Name: "seatalk", Location: "~/seatalk",
	}}}); err != nil {
		t.Fatal(err)
	}

	if err := UpsertProjectAuto(home, ref, "~/seatalk"); err != nil {
		t.Fatal(err)
	}
	auto, err := ReadProjectsAuto(home)
	if err != nil {
		t.Fatal(err)
	}
	if auto.Projects[0].ID != id {
		t.Fatalf("auto id=%d want %d", auto.Projects[0].ID, id)
	}
}

func TestLookupProjectIDRequiresID(t *testing.T) {
	home := t.TempDir()
	if err := WriteProjects(home, ProjectsFile{Projects: []ProjectEntry{{
		Name: "broken", Location: "~/broken",
	}}}); err != nil {
		t.Fatal(err)
	}
	_, ok, err := LookupProjectID(home, ProjectRef{Name: "broken"})
	if ok || err == nil {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
}

func TestEnsureProjectIDCreatesAutoRow(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TSK_DATE", "2026-07-09")

	ref := ProjectRef{Origin: "github.com/a/b"}
	id, err := EnsureProjectID(home, ref, "~/b")
	if err != nil {
		t.Fatal(err)
	}
	if id != 1 {
		t.Fatalf("id=%d", id)
	}
	dir := ProjectNotesDir(home, id)
	want := filepath.Join(home, "projects", "1")
	if dir != want {
		t.Fatalf("dir=%q want %q", dir, want)
	}
	id2, err := EnsureProjectID(home, ref, "~/ignored")
	if err != nil {
		t.Fatal(err)
	}
	if id2 != id {
		t.Fatalf("second ensure %d != %d", id2, id)
	}
}

func TestUpsertProjectAutoInsertAssignsID(t *testing.T) {
	home := t.TempDir()
	t.Setenv("TSK_DATE", "2026-07-09")

	ref := ProjectRef{Name: "seatalk-local-bot"}
	if err := UpsertProjectAuto(home, ref, "~/seatalk-local-bot"); err != nil {
		t.Fatal(err)
	}
	got, err := ReadProjectsAuto(home)
	if err != nil {
		t.Fatal(err)
	}
	if got.Projects[0].ID != 1 {
		t.Fatalf("id=%d want 1", got.Projects[0].ID)
	}
}
