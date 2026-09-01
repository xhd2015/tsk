## Expected Output

```
1
```

## Expected

- Exit code 0.
- Stdout is task id `1` plus trailing `\n` only.
- Stderr empty.
- Task directory `inbox/[1]-hello/` exists with `task.json` and `context/`.
- `index/1` contains `inbox/[1]-hello`.

## Side Effects

- Same filesystem layout as `create/no-topic` (inbox placement).

## Exit Code

- 0

```go
import (
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "1")

	wantRel := inboxTaskRel(1, req.Title)
	taskDir := taskAbs(req, wantRel)
	assertDirExists(t, taskDir)
	assertDirExists(t, filepath.Join(taskDir, "context"))
	assertFileExists(t, filepath.Join(taskDir, "task.json"))
	assertIndexEquals(t, req, 1, wantRel)
}
```