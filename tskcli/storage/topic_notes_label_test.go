package storage

import "testing"

func TestParseLabel(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in        string
		key, val  string
		hasValue  bool
		wantErr   bool
	}{
		{in: "progress", key: "progress"},
		{in: "session=abc", key: "session", val: "abc", hasValue: true},
		{in: "url=a=b", key: "url", val: "a=b", hasValue: true},
		{in: "empty=", key: "empty", val: "", hasValue: true},
		{in: "", wantErr: true},
		{in: "=x", wantErr: true},
	}
	for _, tc := range cases {
		key, val, hasValue, err := ParseLabel(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParseLabel(%q): want error", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseLabel(%q): %v", tc.in, err)
		}
		if key != tc.key || val != tc.val || hasValue != tc.hasValue {
			t.Fatalf("ParseLabel(%q)=(%q,%q,%v), want (%q,%q,%v)",
				tc.in, key, val, hasValue, tc.key, tc.val, tc.hasValue)
		}
	}
}

func TestLabelMatches(t *testing.T) {
	t.Parallel()
	cases := []struct {
		have, want string
		match      bool
	}{
		{"progress", "progress", true},
		{"session=abc", "session=abc", true},
		{"session=abc", "session", true},
		{"session", "session", true},
		{"session=abc", "session=other", false},
		{"progress=done", "progress", true},
		{"grok", "session", false},
		{"session=abc", "grok", false},
	}
	for _, tc := range cases {
		if got := LabelMatches(tc.have, tc.want); got != tc.match {
			t.Fatalf("LabelMatches(%q, %q)=%v, want %v", tc.have, tc.want, got, tc.match)
		}
	}
}

func TestLabelName(t *testing.T) {
	t.Parallel()
	cases := []struct {
		in, want string
		wantErr  bool
	}{
		{in: "progress", want: "progress"},
		{in: "grok-session-id=abc", want: "grok-session-id"},
		{in: "=bad", wantErr: true},
	}
	for _, tc := range cases {
		got, err := LabelName(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("LabelName(%q): want error", tc.in)
			}
			continue
		}
		if err != nil || got != tc.want {
			t.Fatalf("LabelName(%q)=(%q,%v), want %q", tc.in, got, err, tc.want)
		}
	}
}

func TestNoteHasAllLabelsKeyPresence(t *testing.T) {
	t.Parallel()
	n := TopicNote{Labels: []string{"grok", "session=abc"}}
	if !NoteHasAllLabels(n, []string{"session"}) {
		t.Fatal("bare session should match session=abc")
	}
	if !NoteHasAllLabels(n, []string{"grok", "session=abc"}) {
		t.Fatal("exact AND should match")
	}
	if NoteHasAllLabels(n, []string{"session=other"}) {
		t.Fatal("wrong value must not match")
	}
}
