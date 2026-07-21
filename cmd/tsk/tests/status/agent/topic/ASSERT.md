## Expected Output

Leading facts block (before blank line / art). Template uses runtime-exact `dir:`
(not `__DIR__ type=string` — assert-mod non-greedy string placeholder bug).

```
id: <number>
title: topic status fact
stage: create
terminal: false
topic: eng/backend
dir: <exact absolute path from stdout>
```

## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Leading facts order locked: `id` → `title` → `stage` → `terminal` → `topic` → `dir`.
- `topic: eng/backend` (slash-joined `topic_path` segments) after `terminal:`, **above** `dir:`.
- Never omit `topic:`; never use `(not classified yet)` when a topic path is set.
- `dir:` absolute path containing `topics/eng/backend/` (and task dir name with id/stage/slug).
- Key is `dir:` only — no `path:` or `path_rel:`.
- No ANSI; no rectangle chrome.

## Exit Code

- 0

```go
import (
	"path/filepath"
	"strings"
	"fmt"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertAgentNoANSI(t, resp)
	assertAgentNoRectChrome(t, resp.Stdout)

	// Strict leading facts; dir: literal from stdout (Option A)
	assertAgentLeadingFactsShape(t, resp.Stdout, "topic status fact", "create", "false", "eng/backend")

	idStr := fmt.Sprintf("%d", req.TaskID)
	assertAgentFact(t, resp.Stdout, "id", idStr)
	assertAgentFact(t, resp.Stdout, "title", req.Title)
	assertAgentFact(t, resp.Stdout, "stage", "create")
	assertAgentFact(t, resp.Stdout, "terminal", "false")
	assertAgentTopicFact(t, resp.Stdout, "eng/backend")
	assertAgentFactKeyOrder(t, resp.Stdout, "id", "title", "stage", "terminal", "topic", "dir")

	wantRel := topicTaskRel(req.Topic, req.TaskID, "create", req.Title)
	assertAgentDirRel(t, resp.Stdout, wantRel)

	dirVal, ok := parseAgentFactValue(resp.Stdout, "dir")
	if !ok {
		t.Fatalf("dir: missing")
	}
	if !filepath.IsAbs(dirVal) {
		t.Fatalf("dir: must be absolute, got %q", dirVal)
	}
	dirSlash := filepath.ToSlash(dirVal)
	if !strings.Contains(dirSlash, "topics/eng/backend/") {
		t.Fatalf("dir: %q must contain topics/eng/backend/\nstdout:\n%s", dirVal, resp.Stdout)
	}
	wantAbs := findTaskDirByID(t, req, req.TaskID)
	if filepath.Clean(dirVal) != filepath.Clean(wantAbs) {
		t.Fatalf("dir: got %q want on-disk task dir %q", dirVal, wantAbs)
	}

	// Must not print inbox-style unclassified topic when topic_path is set
	if topicVal, ok := parseAgentFactValue(resp.Stdout, "topic"); ok && topicVal == agentInboxTopic {
		t.Fatalf("topic-placed task must not use %q; want eng/backend\nstdout:\n%s",
			agentInboxTopic, resp.Stdout)
	}

	assertAgentNoAltPathKeys(t, resp.Stdout)
}
```
