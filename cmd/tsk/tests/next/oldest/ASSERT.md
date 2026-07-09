## Expected

- Exit code 0.
- Stdout is `1\n` (older task id by `created_at`).
- Stderr empty.

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
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")

	// verify ordering assumption: task 1 created before task 2
	t1 := readTaskJSON(t, findTaskDirByID(t, req, 1))
	t2 := readTaskJSON(t, findTaskDirByID(t, req, 2))
	if !parseCreatedAt(t, t1.CreatedAt).Before(parseCreatedAt(t, t2.CreatedAt)) &&
		!parseCreatedAt(t, t1.CreatedAt).Equal(parseCreatedAt(t, t2.CreatedAt)) {
		t.Fatalf("task 1 should be older than task 2: %s vs %s", t1.CreatedAt, t2.CreatedAt)
	}
}
```