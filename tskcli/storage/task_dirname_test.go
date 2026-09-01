package storage

import "testing"

func TestTaskDirName(t *testing.T) {
	got := TaskDirName(100, "also-enhance-input")
	want := "[100]-also-enhance-input"
	if got != want {
		t.Fatalf("TaskDirName = %q want %q", got, want)
	}
}

func TestParseTaskDirName(t *testing.T) {
	cases := []struct {
		name    string
		wantID  int
		wantSlug string
		ok      bool
	}{
		{"[1]-add-dark-mode", 1, "add-dark-mode", true},
		{"[100]-create-agents-md", 100, "create-agents-md", true}, // slug may start with a stage word
		{"[2]-in_process-x", 2, "in_process-x", true},
		{"[0]-x", 0, "", false},
		{"[1]", 0, "", false},
		{"[1]-", 0, "", false},
		{"1-add-dark-mode", 0, "", false},
		{"topic-name", 0, "", false},
	}
	for _, tc := range cases {
		id, slug, ok := ParseTaskDirName(tc.name)
		if ok != tc.ok || id != tc.wantID || slug != tc.wantSlug {
			t.Fatalf("ParseTaskDirName(%q) = (%d, %q, %v) want (%d, %q, %v)",
				tc.name, id, slug, ok, tc.wantID, tc.wantSlug, tc.ok)
		}
	}
}

func TestInboxRelPathOmitsStage(t *testing.T) {
	got := InboxRelPath(1, "add dark mode")
	want := "inbox/[1]-add-dark-mode"
	if got != want {
		t.Fatalf("InboxRelPath = %q want %q", got, want)
	}
}
