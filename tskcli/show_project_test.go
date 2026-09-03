package tskcli

import (
	"testing"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func TestFormatShowProject(t *testing.T) {
	home := t.TempDir()

	reg := storage.ProjectsFile{Projects: []storage.ProjectEntry{
		{Name: "knowledge-workspace", Location: "~/Projects/kw"},
		{Origin: "github.com/xhd2015/dot-pkgs", Location: "~/dot-pkgs"},
	}}
	if err := storage.WriteProjects(home, reg); err != nil {
		t.Fatal(err)
	}
	auto := storage.ProjectsAutoFile{Projects: []storage.ProjectAutoEntry{
		{Origin: "github.com/xhd2015/tsk", Location: "~/Projects/xhd2015/tsk"},
	}}
	if err := storage.WriteProjectsAuto(home, auto); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		ref  *storage.ProjectRef
		want string
	}{
		{
			name: "location from registry",
			ref:  &storage.ProjectRef{Name: "knowledge-workspace"},
			want: "~/Projects/kw",
		},
		{
			name: "location from registry origin",
			ref:  &storage.ProjectRef{Origin: "github.com/xhd2015/dot-pkgs"},
			want: "~/dot-pkgs",
		},
		{
			name: "location from auto",
			ref:  &storage.ProjectRef{Origin: "github.com/xhd2015/tsk"},
			want: "~/Projects/xhd2015/tsk",
		},
		{
			name: "name fallback",
			ref:  &storage.ProjectRef{Name: "orphan-name"},
			want: "orphan-name",
		},
		{
			name: "origin fallback",
			ref:  &storage.ProjectRef{Origin: "github.com/xhd2015/unknown"},
			want: "github.com/xhd2015/unknown",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := formatShowProject(home, tc.ref)
			if got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
