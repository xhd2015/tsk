# Scenario

**Feature**: re-register fills empty location from probe dir

```
register; strip location from projects.json; register again -> updated location
```

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)

	regPath := filepath.Join(req.TskHome, "projects.json")
	data, err := os.ReadFile(regPath)
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Projects []map[string]any `json:"projects"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	if len(f.Projects) != 1 {
		t.Fatalf("projects=%d", len(f.Projects))
	}
	delete(f.Projects[0], "location")
	out, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	out = append(out, '\n')
	if err := os.WriteFile(regPath, out, 0o600); err != nil {
		t.Fatal(err)
	}

	req.Args = []string{"project", "register", "--name", "seatalk", "--cwd", dir}
	return nil
}
```
