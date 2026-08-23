## Expected

- Non-zero exit; stderr mentions conflict.

## Exit Code

- 1

```go
import (
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for --parent/--topic conflict")
	}
	if !strings.Contains(resp.Stderr, "--topic conflicts with --parent") {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
}
```
