package tskcli

import (
	"strings"
	"testing"
)

func TestFormatProjectListTable_Basic(t *testing.T) {
	rows := []projectListJSONRow{
		{Name: "tsk", Origin: "github.com/xhd2015/tsk", Cwd: "~/Projects/xhd2015/tsk", Tasks: 3},
		{Name: "widget-cli", Origin: "git.example.com/…/widget-cli", Cwd: "~/Projects/…/widget-cli", Tasks: 25},
		{Name: "local-bot", Origin: "", Cwd: "~/local-bot", Tasks: 2},
	}
	got := formatProjectListTable(rows, true)
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	if len(lines) != 4 {
		t.Fatalf("lines=%d got:\n%s", len(lines), got)
	}
	if !strings.HasPrefix(lines[0], "NAME") || !strings.Contains(lines[0], "TASKS") {
		t.Fatalf("header=%q", lines[0])
	}
	// TASKS right-aligned: last non-space runes of data rows
	for i, want := range []string{"3", "25", "2"} {
		line := lines[i+1]
		trimmed := strings.TrimRight(line, " ")
		if !strings.HasSuffix(trimmed, want) {
			t.Fatalf("row %d want suffix %q got %q", i, want, line)
		}
	}
	// Columns align: find ORIGIN start from header
	originIdx := strings.Index(lines[0], "ORIGIN")
	if originIdx < 0 {
		t.Fatal("ORIGIN missing")
	}
	if !strings.HasPrefix(lines[1][originIdx:], "github.com") {
		t.Fatalf("origin column misaligned:\n%s", got)
	}
	cwdIdx := strings.Index(lines[0], "CWD")
	if cwdIdx < 0 {
		t.Fatal("CWD missing")
	}
	// local-bot: blank ORIGIN, cwd under CWD
	if strings.TrimSpace(lines[3][originIdx:cwdIdx]) != "" {
		t.Fatalf("expected blank ORIGIN for local-bot:\n%s", got)
	}
	if !strings.HasPrefix(lines[3][cwdIdx:], "~/local-bot") {
		t.Fatalf("cwd column misaligned:\n%s", got)
	}
}

func TestFormatProjectListTable_NoTasks(t *testing.T) {
	rows := []projectListJSONRow{
		{Name: "local-bot", Cwd: "~/local-bot"},
	}
	got := formatProjectListTable(rows, false)
	if strings.Contains(got, "TASKS") {
		t.Fatalf("unexpected TASKS:\n%s", got)
	}
	if !strings.Contains(got, "NAME") || !strings.Contains(got, "local-bot") {
		t.Fatalf("got:\n%s", got)
	}
}

func TestPadRightLeft(t *testing.T) {
	if got := padRight("ab", 4); got != "ab  " {
		t.Fatalf("padRight=%q", got)
	}
	if got := padLeft("9", 3); got != "  9" {
		t.Fatalf("padLeft=%q", got)
	}
	if got := padRight("世界", 4); runeWidth(got) != 4 {
		t.Fatalf("padRight multibyte width=%d got=%q", runeWidth(got), got)
	}
}
