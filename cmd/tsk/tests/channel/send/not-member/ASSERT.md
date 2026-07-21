## Expected

- Exit code 1; stderr `Error:` naming actor handle and channel.
- No messages written.

## Exit Code

- 1

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	want := `"bob" is not a participant in channel "private-ch"`
	if !strings.Contains(resp.Stderr, want) {
		t.Fatalf("stderr: got %q want substring %q", resp.Stderr, want)
	}
	msgs := readMessagesJSONL(t, activeChannelDir(req, "private-ch"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages, got %d", len(msgs))
	}
}
```
