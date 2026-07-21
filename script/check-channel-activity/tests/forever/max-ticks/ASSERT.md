## Expected

- Exit code 0; stderr empty.
- Stdout contains two status blocks: first `notified`, second `already notified`.
- Marker touched once.

## Exit Code

- 0

```go
import (
	"strings"
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

	blocks := strings.Count(resp.Stdout, "channel: "+channelID)
	if blocks != 2 {
		t.Fatalf("expected 2 status blocks, got %d; stdout=%q", blocks, resp.Stdout)
	}
	assertContains(t, resp.Stdout, "status: notified")
	assertContains(t, resp.Stdout, "status: already notified")
	assertFileExists(t, req.MarkerPath)
}
```