package tskcli

import (
	"testing"

	"github.com/xhd2015/tsk/tskcli/storage"
)

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
