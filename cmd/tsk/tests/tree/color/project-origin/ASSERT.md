## Expected

- Exit 0; stderr empty.
- Project short name `dot-pkgs` is not grey.
- Origin `github.com/xhd2015/dot-pkgs` is grey (`\x1b[90m` … `\x1b[0m`).

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}

	gray, reset := "\x1b[90m", "\x1b[0m"
	wantOrigin := gray + "github.com/xhd2015/dot-pkgs" + reset
	if !strings.Contains(resp.Stdout, "dot-pkgs  "+wantOrigin) {
		t.Fatalf("stdout=%q missing plain short name + grey origin", resp.Stdout)
	}
	// Short name must appear outside a grey wrap that includes it with the origin.
	if strings.Contains(resp.Stdout, gray+"dot-pkgs  github.com/xhd2015/dot-pkgs"+reset) {
		t.Fatalf("short name must not be grey with origin: %q", resp.Stdout)
	}
}
```
