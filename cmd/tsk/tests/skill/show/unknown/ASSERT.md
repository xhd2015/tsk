## Expected

- Non-zero exit; stderr mentions unknown/not found topic.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for unknown topic")
	}
	if !strings.Contains(resp.Stderr, "tsk skill:") {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
	lower := strings.ToLower(resp.Stderr)
	if !strings.Contains(lower, "does not exist") && !strings.Contains(lower, "not found") && !strings.Contains(lower, "unknown") {
		t.Fatalf("stderr should mention missing topic: %q", resp.Stderr)
	}
}
```
