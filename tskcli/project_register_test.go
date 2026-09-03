package tskcli

import (
	"strings"
	"testing"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func TestProjectPathEqual(t *testing.T) {
	t.Parallel()
	home := t.TempDir()
	if !projectPathEqual(home+"/a", home+"/a") {
		t.Fatal("same abs")
	}
	if projectPathEqual(home+"/a", home+"/b") {
		t.Fatal("different abs")
	}
	if !projectPathEqual("", "") {
		t.Fatal("both empty")
	}
	if projectPathEqual(home+"/a", "") {
		t.Fatal("one empty")
	}
}

func TestMergeRegisterField(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		field     string
		existing  string
		incoming  string
		wantStore string
		wantMsg   string
		wantErr   string
	}{
		{
			name: "origin equal", field: "origin",
			existing: "github.com/x", incoming: "github.com/x",
			wantStore: "github.com/x",
		},
		{
			name: "fill location", field: "location",
			existing: "", incoming: "/tmp/proj",
			wantStore: "/tmp/proj", wantMsg: "updated location: (empty) -> /tmp/proj",
		},
		{
			name: "conflict location", field: "location",
			existing: "/tmp/a", incoming: "/tmp/b",
			wantErr: "different location",
		},
		{
			name: "clear forbidden", field: "origin",
			existing: "github.com/x", incoming: "",
			wantErr: "different origin",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, msg, err := mergeRegisterField("seatalk", tc.field, tc.existing, tc.incoming)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("err=%v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.wantStore {
				t.Fatalf("store=%q want %q", got, tc.wantStore)
			}
			if msg != tc.wantMsg {
				t.Fatalf("msg=%q want %q", msg, tc.wantMsg)
			}
		})
	}
}

func TestFindRegisterMatchIndex(t *testing.T) {
	t.Parallel()
	reg := storage.ProjectsFile{Projects: []storage.ProjectEntry{
		{Name: "by-loc", Location: "/main/repo"},
		{Name: "other", Location: "/only/other"},
	}}
	if i := findRegisterMatchIndex(reg, "/main/repo", "/main/repo"); i != 0 {
		t.Fatalf("step1/2 got %d", i)
	}
	if i := findRegisterMatchIndex(reg, "/probe", "/only/other"); i != 1 {
		t.Fatalf("incoming location got %d", i)
	}
	if i := findRegisterMatchIndex(reg, "/x", "/z"); i != -1 {
		t.Fatalf("miss got %d", i)
	}
}

func TestRegisterAutoName(t *testing.T) {
	t.Parallel()
	if got := registerAutoName("/tmp/foo/bar"); got != "bar" {
		t.Fatalf("got %q", got)
	}
	if got := registerAutoName(""); got != "" {
		t.Fatalf("got %q", got)
	}
}
