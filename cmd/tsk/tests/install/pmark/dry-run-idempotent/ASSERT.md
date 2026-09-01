## Expected

- Exit 0.
- Stdout skips wrapper and `.zshrc` checker.
- No `would write` / `would create` / `would overwrite`.
- Wrapper body still the forwarder.

## Exit Code

- 0

```go
import "strings"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q stdout=%q", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	assertContains(t, resp.Stdout, "[dry-run] skip: ~/.local/bin/pmark (already identical)")
	assertContains(t, resp.Stdout, "[dry-run] skip: ~/.zshrc (PATH checker already present)")
	for _, bad := range []string{"would write", "would create", "would overwrite", "would append", "would replace"} {
		if strings.Contains(resp.Stdout, bad) {
			t.Fatalf("unexpected %q in stdout:\n%s", bad, resp.Stdout)
		}
	}
	want := "#!/bin/sh\nexec tsk project add \"$@\"\n"
	if got := readLocalBin(t, req, "pmark"); got != want {
		t.Fatalf("wrapper mutated: %q", got)
	}
}
```
