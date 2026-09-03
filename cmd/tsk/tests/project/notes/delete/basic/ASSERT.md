## Expected

- Exit 0; only `keep` remains; footer `1 notes`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "keep")
	if strings.Contains(resp.Stdout, "drop") {
		t.Fatalf("deleted note still present: %q", resp.Stdout)
	}
	assertContains(t, resp.Stdout, "1 notes\n")
}
```

