package tskcli

import (
	"testing"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func TestOrderedProjectGroupsRootFirst(t *testing.T) {
	groups := map[string]*projectGroup{
		"github.com/example/aaa": {Key: "github.com/example/aaa", Label: "aaa  github.com/example/aaa"},
		"github.com/example/zzz": {Key: "github.com/example/zzz", Label: "zzz  github.com/example/zzz"},
		"github.com/example/mmm": {Key: "github.com/example/mmm", Label: "mmm  github.com/example/mmm"},
	}
	ordered := orderedProjectGroups(groups, "github.com/example/zzz")
	if len(ordered) != 3 {
		t.Fatalf("len=%d", len(ordered))
	}
	if ordered[0].Key != "github.com/example/zzz" {
		t.Fatalf("root first: got %q", ordered[0].Key)
	}
	if ordered[1].Key != "github.com/example/aaa" || ordered[2].Key != "github.com/example/mmm" {
		t.Fatalf("rest not label-sorted: %q %q", ordered[1].Key, ordered[2].Key)
	}
}

func TestTaskMatchesProjectKeySet(t *testing.T) {
	keys := map[string]struct{}{
		"github.com/xhd2015/dot-pkgs": {},
		"name:local-bot":              {},
	}
	cases := []struct {
		name string
		task storage.Task
		want bool
	}{
		{
			name: "nil project",
			task: storage.Task{},
			want: false,
		},
		{
			name: "origin hit",
			task: storage.Task{Project: &storage.ProjectRef{Origin: "github.com/xhd2015/dot-pkgs"}},
			want: true,
		},
		{
			name: "origin miss",
			task: storage.Task{Project: &storage.ProjectRef{Origin: "github.com/xhd2015/wrk"}},
			want: false,
		},
		{
			name: "name hit",
			task: storage.Task{Project: &storage.ProjectRef{Name: "local-bot"}},
			want: true,
		},
		{
			name: "name miss",
			task: storage.Task{Project: &storage.ProjectRef{Name: "other"}},
			want: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := taskMatchesProjectKeySet(tc.task, keys)
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
	if taskMatchesProjectKeySet(storage.Task{Project: &storage.ProjectRef{Origin: "github.com/xhd2015/dot-pkgs"}}, nil) {
		t.Fatal("nil keys should not match")
	}
}
