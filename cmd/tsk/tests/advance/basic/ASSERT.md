## Expected

- Exit code 0.
- Directory basename stays `inbox/[1]-add-dark-mode/` (no rename on advance).
- `index/1` still points at that path.
- `task.json` stage is `in_process` with `stage_history` entry `create` → `in_process`.

## Side Effects

- Stage and stage_history updated in `task.json` only.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	wantRel := inboxTaskRel(req.TaskID, req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTaskStage(t, req, req.TaskID, "in_process")

	task := readTaskJSON(t, findTaskDirByID(t, req, req.TaskID))
	if len(task.StageHistory) == 0 {
		t.Fatal("stage_history should not be empty")
	}
	last := task.StageHistory[len(task.StageHistory)-1]
	if last["from"] != "create" || last["to"] != "in_process" {
		t.Fatalf("stage_history last entry: got %v", last)
	}
}
```
