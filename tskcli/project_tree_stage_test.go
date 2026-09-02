package tskcli

import "testing"

func TestProjectTreeStageAllowed(t *testing.T) {
	cases := []struct {
		name                     string
		stage, stageFlag         string
		done, archived, allFlag  bool
		want                     bool
	}{
		{"default hides done", "done", "", false, false, false, false},
		{"default hides archived", "archived", "", false, false, false, false},
		{"default keeps create", "create", "", false, false, false, true},
		{"done only accepts done", "done", "", true, false, false, true},
		{"done only rejects archived", "archived", "", true, false, false, false},
		{"done only rejects create", "create", "", true, false, false, false},
		{"archived only accepts archived", "archived", "", false, true, false, true},
		{"archived only rejects done", "done", "", false, true, false, false},
		{"both terminals accept done", "done", "", true, true, false, true},
		{"both terminals accept archived", "archived", "", true, true, false, true},
		{"both terminals reject create", "create", "", true, true, false, false},
		{"all accepts done", "done", "", false, false, true, true},
		{"all accepts create", "create", "", false, false, true, true},
		{"all+done narrows to done", "archived", "", true, false, true, false},
		{"all+done keeps done", "done", "", true, false, true, true},
		{"stage exact", "implementation", "implementation", false, false, false, true},
		{"stage mismatch", "create", "implementation", false, false, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := projectTreeStageAllowed(tc.stage, tc.stageFlag, tc.done, tc.archived, tc.allFlag)
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}
