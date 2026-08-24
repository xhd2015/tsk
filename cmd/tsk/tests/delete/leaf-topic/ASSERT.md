## Expected

- Exit code 0; stdout `deleted <id>`.
- Topic task dir and index gone; topic directory may remain.

## Exit Code

- 0

```go
import (
	"fmt"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 || resp.Stderr != "" {
		t.Fatalf("exit=%d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, fmt.Sprintf("deleted %d", req.TaskID))

	wantRel := topicTaskRel(req.Topic, req.TaskID, "create", req.Title)
	assertFileNotExists(t, taskAbs(req, wantRel))
	assertFileNotExists(t, filepath.Join(req.TskHome, "index", fmt.Sprintf("%d", req.TaskID)))
}
```
