package storage

import "testing"

func TestIsTerminal(t *testing.T) {
	cases := []struct {
		stage string
		want  bool
	}{
		{"create", false},
		{"in_process", false},
		{"summary", false},
		{"user_followup", false},
		{"done", true},
		{"archived", true},
		{"", false},
	}
	for _, tc := range cases {
		if got := IsTerminal(tc.stage); got != tc.want {
			t.Fatalf("IsTerminal(%q)=%v want %v", tc.stage, got, tc.want)
		}
	}
}

func TestTerminalTransitionError(t *testing.T) {
	if err := TerminalTransitionError("done"); err == nil || err.Error() != "invalid transition: task is already done" {
		t.Fatalf("done: %v", err)
	}
	if err := TerminalTransitionError("archived"); err == nil || err.Error() != "invalid transition: task is already archived" {
		t.Fatalf("archived: %v", err)
	}
}
