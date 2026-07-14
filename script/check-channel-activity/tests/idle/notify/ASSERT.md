## Expected

- Exit code 0; stderr empty.
- Stdout `status: notified`.
- Notify marker file created.
- State file records `last_activity_at` and `last_notified_at`.

## Side Effects

- `channels/state/eng-alerts.json` written.

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
	assertFileExists(t, statePath(req))

	data, err := os.ReadFile(statePath(req))
	if err != nil {
		t.Fatalf("read state: %v", err)
	}
	var st notifyState
	if err := json.Unmarshal(data, &st); err != nil {
		t.Fatalf("parse state: %v", err)
	}
	if st.ChannelID != req.ChannelID {
		t.Fatalf("state channel_id: got %q want %q", st.ChannelID, req.ChannelID)
	}
	if st.LastActivityAt != req.LastActivity {
		t.Fatalf("state last_activity_at: got %q want %q", st.LastActivityAt, req.LastActivity)
	}
	if st.LastNotifiedAt == "" {
		t.Fatal("state last_notified_at should be set")
	}
}
```