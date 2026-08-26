## Expected

- Exit code 0; stderr empty.
- Stdout starts with `Usage: tsk search`.
- Mentions `--task`, `--note`, `--progress`, `--topic`, `--all`, `--json`.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	if !strings.HasPrefix(resp.Stdout, "Usage: tsk search") {
		t.Fatalf("stdout=%q missing usage prefix", resp.Stdout)
	}
	for _, flag := range []string{"--task", "--note", "--progress", "--topic", "--all", "--color", "--no-color", "--json"} {
		if !strings.Contains(resp.Stdout, flag) {
			t.Fatalf("stdout=%q missing %q", resp.Stdout, flag)
		}
	}
}
```
