## Expected

- Exit code ≠ 0; stderr has exactly one `Error:` prefix.
- Stderr indicates parent `--channel-id` is not accepted for `list`.

## Errors

- Parent `--channel-id` on `list` is a hard reject (not ignored).

## Exit Code

- non-zero

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit when list gets parent --channel-id")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	assertContains(t, resp.Stderr, "--channel-id")
	low := strings.ToLower(resp.Stderr)
	if !strings.Contains(low, "not accepted") && !strings.Contains(low, "not allowed") && !strings.Contains(low, "cannot") {
		t.Fatalf("expected rejection wording for parent --channel-id on list, got %q", resp.Stderr)
	}
}
```
