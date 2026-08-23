## Expected

- Exit code 0; stderr empty.
- `topics/knowledge-base/topic.json` has alias `知识库`.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	p := filepath.Join(req.TskHome, "topics", "knowledge-base", "topic.json")
	assertFileExists(t, p)
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read topic.json: %v", err)
	}
	var meta struct {
		Aliases []string `json:"aliases"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		t.Fatalf("parse topic.json: %v", err)
	}
	found := false
	for _, a := range meta.Aliases {
		if a == "知识库" {
			found = true
		}
	}
	if !found {
		t.Fatalf("aliases missing 知识库: %v", meta.Aliases)
	}
}
```
