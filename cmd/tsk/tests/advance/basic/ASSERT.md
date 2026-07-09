## Expected

- Exit code 0.
- Old dir `inbox/1-create-add-dark-mode/` no longer exists.
- New dir `inbox/1-in_process-add-dark-mode/` exists.
- `index/1` updated to `inbox/1-in_process-add-dark-mode`.
- `task.json` stage is `in_process` with `stage_history` entry `create` → `in_process`.

## Side Effects

- Directory rename under inbox; index rewrite.

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

	oldRel := inboxTaskRel(req.TaskID, "create", req.Title)
	newRel := inboxTaskRel(req.TaskID, "in_process", req.Title)
	assertFileNotExists(t, taskAbs(req, oldRel))
	assertDirExists(t, taskAbs(req, newRel))
	assertIndexEquals(t, req, req.TaskID, newRel)
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