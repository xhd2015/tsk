## Expected

- Exit code 0; stderr empty.
- Stdout contains `[1]`, `[2]`, `alice`, `first`, `second` in order.
- Stdout ends with `\n`.

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
	if !strings.HasSuffix(resp.Stdout, "\n") {
		t.Fatalf("stdout should end with newline")
	}
	i1 := strings.Index(resp.Stdout, "first")
	i2 := strings.Index(resp.Stdout, "second")
	if i1 < 0 || i2 < 0 || i1 > i2 {
		t.Fatalf("expected first before second in stdout: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "[1]")
	assertContains(t, resp.Stdout, "[2]")
	assertContains(t, resp.Stdout, "alice")
}
```
