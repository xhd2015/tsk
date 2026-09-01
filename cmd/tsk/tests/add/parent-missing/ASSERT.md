## Expected

- Non-zero exit; stderr mentions parent not found.

## Exit Code

- 1

```go
import (
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for missing parent")
	}
	if !strings.Contains(resp.Stderr, "parent task not found") {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
}
```
