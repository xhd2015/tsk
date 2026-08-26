# Scenario

**Feature**: `--json` with `--color` stays plain (no ANSI)

```
note -> search --json --color token -> JSON array, no escapes
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "json color", "eng/backend", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "json-color-token")
	req.Args = []string{"search", "--json", "--color", "json-color-token"}
	return nil
}
```
