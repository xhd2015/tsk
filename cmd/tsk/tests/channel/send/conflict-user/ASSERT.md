## Expected

- Exit code ≠ 0; stderr has exactly one `Error:` prefix.
- Stderr mentions conflict and both user values (`alice`, `bob`).
- No messages written under `team-ch`.

## Errors

- Conflicting parent vs leaf `--user` (not silent leaf-wins).

## Exit Code

- non-zero

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for conflicting --user")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	low := strings.ToLower(resp.Stderr)
	if !strings.Contains(low, "conflict") {
		t.Fatalf("expected conflict in stderr, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stderr, "alice")
	assertContains(t, resp.Stderr, "bob")

	msgs := readMessagesJSONL(t, activeChannelDir(req, "team-ch"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages on conflict, got %d", len(msgs))
	}
}
```
