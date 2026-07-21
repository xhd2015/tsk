## Expected

- Exit code 1; stderr `Error:` naming actor handle and channel.
- `participants.jsonl` unchanged (`alice` only).

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
	want := `"charlie" is not a participant in channel "closed-ch"`
	if !strings.Contains(resp.Stderr, want) {
		t.Fatalf("stderr: got %q want substring %q", resp.Stderr, want)
	}
	assertParticipantHandlesSorted(t, activeChannelDir(req, req.ChannelID), []string{"alice"})
}
```
