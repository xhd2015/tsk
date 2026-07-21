## Expected

- Exit code 1; stderr contains `Error:`.
- No channel directory created.

## Exit Code

- 1

```go
import (
	"os"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	_, err = os.Stat(channelAbs(req, "index"))
	if err == nil {
		entries, _ := os.ReadDir(channelAbs(req, "active"))
		if len(entries) > 0 {
			t.Fatalf("expected no active channels")
		}
	}
}
```
