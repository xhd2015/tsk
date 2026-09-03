package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeProjectEntry(t *testing.T) {
	t.Parallel()
	e := ProjectEntry{Cwd: "~/probe"}
	normalizeProjectEntry(&e)
	if e.Location != "~/probe" || e.Cwd != "" {
		t.Fatalf("%+v", e)
	}
	e2 := ProjectEntry{Location: "~/main", Cwd: "~/probe"}
	normalizeProjectEntry(&e2)
	if e2.Location != "~/main" || e2.Cwd != "" {
		t.Fatalf("%+v", e2)
	}
}

func TestReadProjectsMigratesLegacyCwd(t *testing.T) {
	home := t.TempDir()
	path := filepath.Join(home, "projects.json")
	if err := os.WriteFile(path, []byte(`{
  "projects": [
    { "name": "seatalk", "cwd": "~/seatalk-local-bot" }
  ]
}
`), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ReadProjects(home)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Projects) != 1 {
		t.Fatalf("len=%d", len(got.Projects))
	}
	e := got.Projects[0]
	if e.Location != "~/seatalk-local-bot" || e.Cwd != "" {
		t.Fatalf("%+v", e)
	}
	if err := WriteProjects(home, got); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), `"cwd"`) {
		t.Fatalf("wrote cwd: %s", raw)
	}
	if !strings.Contains(string(raw), `"location"`) {
		t.Fatalf("missing location: %s", raw)
	}
}

func TestUpsertProjectAutoFillsLocationFromLegacyCwd(t *testing.T) {
	home := t.TempDir()
	ref := ProjectRef{Origin: "github.com/xhd2015/dot-pkgs"}
	path := filepath.Join(home, "projects-auto.json")
	if err := os.WriteFile(path, []byte(`{
  "projects": [
    {
      "origin": "github.com/xhd2015/dot-pkgs",
      "cwd": "~/Projects/xhd2015/dot-pkgs",
      "first_seen_at": "2026-07-09T12:00:00+08:00",
      "last_seen_at": "2026-07-09T12:00:00+08:00"
    }
  ]
}
`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TSK_DATE", "2026-07-10")
	if err := UpsertProjectAuto(home, ref, "~/ignored-on-update"); err != nil {
		t.Fatal(err)
	}
	got, err := ReadProjectsAuto(home)
	if err != nil {
		t.Fatal(err)
	}
	e := got.Projects[0]
	if e.Location != "~/Projects/xhd2015/dot-pkgs" {
		t.Fatalf("location=%q", e.Location)
	}
	if e.Cwd != "" {
		t.Fatalf("cwd should be cleared, got %q", e.Cwd)
	}
	if e.LastSeenAt != "2026-07-10T12:00:00+08:00" {
		t.Fatalf("last_seen=%q", e.LastSeenAt)
	}
	raw, _ := os.ReadFile(path)
	if strings.Contains(string(raw), `"cwd"`) {
		t.Fatalf("auto still has cwd: %s", raw)
	}
}

func TestUpsertProjectAutoInsertSetsLocation(t *testing.T) {
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
	e := got.Projects[0]
	if e.Location != "~/seatalk-local-bot" || e.Cwd != "" {
		t.Fatalf("entry=%+v", e)
	}
}
