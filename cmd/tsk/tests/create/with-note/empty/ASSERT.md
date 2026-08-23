## Expected

- Non-zero exit; stderr mentions `--note text required`.

## Exit Code

- 1

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for empty --note")
	}
	if !strings.Contains(resp.Stderr, "--note text required") {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
}
```
