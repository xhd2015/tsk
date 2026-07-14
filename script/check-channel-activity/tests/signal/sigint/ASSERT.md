---
label: slow
explanation: starts subprocess, waits for signal handling
---

## Expected

- Exit code 0.
- Stderr exactly `stopped\n`.
- No notify marker.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stdout=%q stderr=%q", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if resp.Stderr != "stopped\n" {
		t.Fatalf("stderr: got %q want %q", resp.Stderr, "stopped\n")
	}
	assertFileNotExists(t, req.MarkerPath)
}
```