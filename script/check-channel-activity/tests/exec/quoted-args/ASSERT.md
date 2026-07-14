## Expected

- Exit code 0; stderr empty.
- Stdout `status: notified`.
- `argv.txt` contains `hello world` (quoted spaces preserved as one argument).

## Side Effects

- `argv.txt` written under WorkRoot.

## Exit Code

- 0

```go
import (
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)
	assertStatusBlock(t, resp.Stdout, "notified", req.LastActivity)

	data, err := os.ReadFile(req.ArgvPath)
	if err != nil {
		t.Fatalf("read argv file: %v", err)
	}
	got := strings.TrimSpace(string(data))
	if got != "hello world" {
		t.Fatalf("argv file: got %q want %q", got, "hello world")
	}
}
```