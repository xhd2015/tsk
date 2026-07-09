## Expected

- Exit code 0.
- `task.json` labels are sorted lexicographically: `["bug", "urgent"]`.

## Side Effects

- Task created in inbox at `inbox/1-create-x/`.

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
	assertDirExists(t, taskAbs(req, wantRel))
	assertLabelsSorted(t, req, 1, []string{"bug", "urgent"})
}
```