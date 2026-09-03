## Expected

- Exit 0.
- Stdout includes `updated location: (empty) ->`.
- `projects.json` has a non-empty `location`.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertContains(t, resp.Stdout, "updated location: (empty) ->")

	data, err := os.ReadFile(filepath.Join(req.TskHome, "projects.json"))
	if err != nil {
		t.Fatal(err)
	}
	assertContains(t, string(data), `"location":`)
}
```
