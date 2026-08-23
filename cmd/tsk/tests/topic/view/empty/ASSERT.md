## Expected

- Exit code 0; stderr empty.
- First line is `knowledge-base`.
- Second line is `(empty)`.
- No `topic.json` / `notes.jsonl` nodes.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if resp.Stdout != "knowledge-base\n(empty)\n" {
		t.Fatalf("stdout=%q", resp.Stdout)
	}
}
```
