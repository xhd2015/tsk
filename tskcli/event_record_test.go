package tskcli

import (
	"testing"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func TestDeriveAction(t *testing.T) {
	cases := []struct {
		cmd  string
		args []string
		want string
	}{
		{"add", []string{"ship logs"}, "add"},
		{"note", []string{"add", "--id", "42", "hello"}, "note.add"},
		{"note", []string{"list", "--id", "42"}, "note.list"},
		{"channel", []string{"participant", "add", "--channel-id", "x", "bob"}, "channel.participant.add"},
		{"project", []string{"notes", "add", "--name", "app", "hello"}, "project.notes.add"},
		{"project", []string{"tree", "--all"}, "project.tree"},
		{"skill", []string{"--install", "--dir", "~/s"}, "skill.install"},
		{"skill", []string{"--show", "overview"}, "skill.show"},
		{"topic", []string{"alias", "add", "eng", "知识库"}, "topic.alias.add"},
		{"", nil, ""},
	}
	for _, tc := range cases {
		got := deriveAction(tc.cmd, tc.args)
		if got != tc.want {
			t.Errorf("deriveAction(%q, %q)=%q want %q", tc.cmd, tc.args, got, tc.want)
		}
	}
}

func TestIsMutationAction(t *testing.T) {
	if !isMutationAction("note.add") {
		t.Fatal("note.add should be a mutation")
	}
	if isMutationAction("note.list") {
		t.Fatal("note.list should not be a mutation")
	}
	if !isMutationAction("add") {
		t.Fatal("add should be a mutation")
	}
	if isMutationAction("list") {
		t.Fatal("list should not be a mutation")
	}
}

func TestEventIsMutationLegacy(t *testing.T) {
	legacyAdd := storage.Event{Command: "add", Args: []string{"ship logs"}}
	if !eventIsMutation(legacyAdd) {
		t.Fatal("legacy add should classify as mutation")
	}
	legacyList := storage.Event{Command: "list", Args: []string{}}
	if eventIsMutation(legacyList) {
		t.Fatal("legacy list should not classify as mutation")
	}
	newRead := storage.Event{Command: "list", Action: "list", Mutation: false}
	if eventIsMutation(newRead) {
		t.Fatal("new list with mutation=false should not display as mutation")
	}
}

func TestFormatLogLine(t *testing.T) {
	got := formatLogLine(storage.Event{
		TS:       "2026-07-09T12:00:00+08:00",
		Command:  "note",
		Action:   "note.add",
		ExitCode: 0,
		Data:     &storage.EventData{TaskID: 42},
	})
	want := "2026-07-09T12:00:00+08:00  ok  note.add  task=42"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	failLine := formatLogLine(storage.Event{
		TS:       "2026-07-09T12:00:00+08:00",
		Action:   "done",
		ExitCode: 1,
		Data:     &storage.EventData{TaskID: 42},
	})
	if failLine != "2026-07-09T12:00:00+08:00  fail  done  task=42" {
		t.Fatalf("fail line %q", failLine)
	}
}

func TestLogsCountLabel(t *testing.T) {
	if logsCountLabel(0) != "0 logs" {
		t.Fatalf("0: %q", logsCountLabel(0))
	}
	if logsCountLabel(1) != "1 log" {
		t.Fatalf("1: %q", logsCountLabel(1))
	}
	if logsCountLabel(2) != "2 logs" {
		t.Fatalf("2: %q", logsCountLabel(2))
	}
}

func TestApplyEventLimit(t *testing.T) {
	evs := []storage.Event{{Action: "a"}, {Action: "b"}, {Action: "c"}}
	got := applyEventLimit(evs, 1)
	if len(got) != 1 || got[0].Action != "c" {
		t.Fatalf("limit 1: %+v", got)
	}
	if len(applyEventLimit(evs, 0)) != 3 {
		t.Fatal("limit 0 should keep all")
	}
}

func TestInvocationIsolation(t *testing.T) {
	a := &invocation{}
	b := &invocation{}
	a.setCommand("add", []string{"one"})
	a.setData(storage.EventData{TaskID: 1, Text: "a"})
	b.setCommand("list", nil)
	b.setData(storage.EventData{TaskID: 2})
	if a.command != "add" || b.command != "list" {
		t.Fatalf("commands mixed: %q %q", a.command, b.command)
	}
	if a.data == nil || b.data == nil || a.data.TaskID != 1 || b.data.TaskID != 2 {
		t.Fatalf("data mixed: %+v %+v", a.data, b.data)
	}
	if a.data.Text != "a" {
		t.Fatalf("a.text=%q", a.data.Text)
	}
	if !a.mutation || b.mutation {
		t.Fatalf("mutation mixed: %v %v", a.mutation, b.mutation)
	}
}
