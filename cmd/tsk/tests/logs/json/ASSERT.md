## Expected

- Exit 0; JSON array of two mutations (`add`, `note.add`).
- `ts` is `2026-07-09T12:00:00+08:00`; `user` is `alice`.
- `note.add` `data.text` is `hello`; `data.task_id` matches; no ANSI.

## Exit Code

- 0

```go
import (
	"encoding/json"
	"strings"
)

type logEventJSON struct {
	TS       string `json:"ts"`
	Command  string `json:"command"`
	Action   string `json:"action"`
	Mutation bool   `json:"mutation"`
	User     string `json:"user"`
	Data     *struct {
		TaskID int    `json:"task_id"`
		Title  string `json:"title"`
		Text   string `json:"text"`
	} `json:"data"`
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr=%q", resp.Stderr)
	}
	if strings.ContainsRune(resp.Stdout, '\x1b') {
		t.Fatalf("ANSI in json: %q", resp.Stdout)
	}
	var arr []logEventJSON
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse json: %v stdout=%q", err, resp.Stdout)
	}
	if len(arr) != 2 {
		t.Fatalf("len=%d stdout=%q", len(arr), resp.Stdout)
	}
	const ts = "2026-07-09T12:00:00+08:00"
	if arr[0].TS != ts || arr[0].Action != "add" || !arr[0].Mutation {
		t.Fatalf("add event: %+v", arr[0])
	}
	if arr[0].User != "alice" {
		t.Fatalf("user=%q", arr[0].User)
	}
	if arr[0].Data == nil || arr[0].Data.TaskID != req.TaskID || arr[0].Data.Title != "ship logs" {
		t.Fatalf("add data: %+v", arr[0].Data)
	}
	if arr[1].Action != "note.add" || !arr[1].Mutation || arr[1].Data == nil {
		t.Fatalf("note event: %+v", arr[1])
	}
	if arr[1].Data.TaskID != req.TaskID || arr[1].Data.Text != "hello" {
		t.Fatalf("note data: %+v", arr[1].Data)
	}
}
```
