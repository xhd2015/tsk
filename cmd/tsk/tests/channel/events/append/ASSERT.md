## Expected

- Exit code 0.
- `events.jsonl` last line has `command: channel`, `exit_code: 0`.

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
	"os"
	"encoding/json"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertEventsCountAtLeast(t, req, 1)
	assertChannelEventCommand(t, req)
	path := filepath.Join(req.TskHome, "events.jsonl")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read events.jsonl: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	var ev eventLine
	if err := json.Unmarshal([]byte(lines[len(lines)-1]), &ev); err != nil {
		t.Fatalf("parse event: %v", err)
	}
	if ev.ExitCode != 0 {
		t.Fatalf("event exit_code: got %d want 0", ev.ExitCode)
	}
}
```
