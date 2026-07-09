## Expected

- Exit code 0.
- Task directory `inbox/1-create-add-dark-mode/` exists with `task.json` and `context/`.
- `index/1` contains `inbox/1-create-add-dark-mode`.
- `task.json` has `id: 1`, `title: "add dark mode"`, `slug: "add-dark-mode"`, `stage: "create"`, `topic_path: null`.
- `context/` exists and is empty.

## Side Effects

- `counter` incremented; `events.jsonl` gains a line (see `events/append`).

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	wantRel := inboxTaskRel(1, "create", req.Title)
	taskDir := taskAbs(req, wantRel)
	assertDirExists(t, taskDir)
	assertDirExists(t, filepath.Join(taskDir, "context"))
	assertFileExists(t, filepath.Join(taskDir, "task.json"))
	assertIndexEquals(t, req, 1, wantRel)

	task := readTaskJSON(t, taskDir)
	if task.ID != 1 {
		t.Fatalf("id: got %d want 1", task.ID)
	}
	if task.Title != req.Title {
		t.Fatalf("title: got %q want %q", task.Title, req.Title)
	}
	if task.Slug != "add-dark-mode" {
		t.Fatalf("slug: got %q want add-dark-mode", task.Slug)
	}
	if task.Stage != "create" {
		t.Fatalf("stage: got %q want create", task.Stage)
	}
	assertTopicPathNull(t, req, 1)

	entries, err := os.ReadDir(filepath.Join(taskDir, "context"))
	if err != nil {
		t.Fatalf("read context/: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("context/ should be empty, got %d entries", len(entries))
	}
}
```