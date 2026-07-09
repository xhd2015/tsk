## Expected

- Exit code 0.
- Stage is `implementation`; directory renamed to `*-implementation-*`.
- `clarify/batch.json` status is `confirmed` (or all items confirmed).
- `index/<id>` updated to new relative path.

## Side Effects

- Auto-advance from clarification to implementation on full confirm.

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

	wantRel := inboxTaskRel(req.TaskID, "implementation", req.Title)
	assertDirExists(t, taskAbs(req, wantRel))
	assertIndexEquals(t, req, req.TaskID, wantRel)
	assertTaskStage(t, req, req.TaskID, "implementation")

	batch := readClarifyBatch(t, findTaskDirByID(t, req, req.TaskID))
	if batch.Status != "confirmed" && batch.Status != "closed" {
		// allow either confirmed or all items confirmed
		allConfirmed := len(batch.Items) == 2
		for _, item := range batch.Items {
			if item.Status != "confirmed" {
				allConfirmed = false
				break
			}
		}
		if !allConfirmed {
			t.Fatalf("batch status %q; items not all confirmed: %+v", batch.Status, batch.Items)
		}
	}
}
```