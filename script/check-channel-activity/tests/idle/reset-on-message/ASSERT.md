## Expected

- Exit code 0 on second run.
- Stdout `status: notified` (new last_activity still idle).
- Marker touched again on second run.
- State `last_activity_at` updated to newer message timestamp.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"os"
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
	assertFileExists(t, req.MarkerPath)

	data, err := os.ReadFile(statePath(req))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	var st notifyState
	if err := json.Unmarshal(data, &st); err != nil {
		t.Fatalf("parse state: %v", err)
	}
	if st.LastActivityAt != req.LastActivity {
		t.Fatalf("state last_activity_at: got %q want %q", st.LastActivityAt, req.LastActivity)
	}
}
```