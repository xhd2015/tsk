# Scenario

**Feature**: `--color` wins over `NO_COLOR=1`

```
NO_COLOR=1 + search --color token -> ANSI present
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "env color", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "env-color-token")
	req.ExtraEnv = []string{"NO_COLOR=1"}
	req.Args = []string{"search", "--color", "env-color-token"}
	return nil
}
```
