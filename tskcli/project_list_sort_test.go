package tskcli

import "testing"

func TestSortProjectListByTasksDesc(t *testing.T) {
	rows := []projectListJSONRow{
		{Name: "a", Origin: "o/a", Tasks: 1},
		{Name: "c", Origin: "o/c", Tasks: 5},
		{Name: "b", Origin: "o/b", Tasks: 5},
		{Name: "d", Origin: "o/d", Tasks: 0},
	}
	sortProjectListByTasksDesc(rows)
	want := []struct {
		name  string
		tasks int
	}{
		{"b", 5},
		{"c", 5},
		{"a", 1},
		{"d", 0},
	}
	for i, w := range want {
		if rows[i].Name != w.name || rows[i].Tasks != w.tasks {
			t.Fatalf("row %d: got {%s,%d} want {%s,%d}", i, rows[i].Name, rows[i].Tasks, w.name, w.tasks)
		}
	}
}
