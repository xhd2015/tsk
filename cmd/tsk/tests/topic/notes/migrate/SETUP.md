# Scenario

**Feature**: legacy topic.json notes blob is copied into notes.jsonl once

```
topic.json notes blob -> topic notes lists it and clears the blob
```

## Steps

1. mkdir `knowledge-base`.
2. Write `topic.json` with a `notes` string.
3. `tsk topic notes knowledge-base`.

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "knowledge-base")
	dir := filepath.Join(req.TskHome, "topics", "knowledge-base")
	body := `{
  "path": ["knowledge-base"],
  "title": "knowledge-base",
  "aliases": [],
  "notes": "legacy blob",
  "updated_at": "2026-07-09T03:00:00Z"
}
`
	if err := os.WriteFile(filepath.Join(dir, "topic.json"), []byte(body), 0o644); err != nil {
		return err
	}
	req.Args = []string{"topic", "notes", "knowledge-base"}
	return nil
}
```
