## Expected Output

Leading facts block (before blank line / art). Template uses runtime-exact `dir:`
(not `__DIR__ type=string` — assert-mod non-greedy string placeholder bug).

```
id: <number>
title: add dark mode
stage: create
terminal: false
topic: (not classified yet)
dir: <exact absolute path from stdout>
```

## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Leading facts order locked: `id` → `title` → `stage` → `terminal` → `topic` → `dir`.
- Inbox `topic: (not classified yet)` immediately above `dir:`.
- `dir:` value is an **absolute** filesystem path to the task directory.
- Path contains or ends with `inbox/<id>-create-add-dark-mode` (stage segment in dir name).
- Key is `dir:` only — no `path:` or `path_rel:`.
- Homes / temp roots vary: `assert.Output` uses exact stdout `dir:` as a literal line; Go checks abs + suffix + on-disk equality.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertAgentNoANSI(t, resp)
	assertAgentNoRectChrome(t, resp.Stdout)

	// Strict leading facts shape; dir: literal from stdout (Option A)
	assertAgentLeadingFactsShape(t, resp.Stdout, "add dark mode", "create", "false", agentInboxTopic)

	// Core fact values + locked order including topic before dir
	assertAgentCoreFacts(t, resp.Stdout, req.TaskID, req.Title, "create", "false")
	assertAgentDirFact(t, resp.Stdout, req.TaskID, "create", req.Title)

	// Absolute + relative segment without hardcoding home
	dirVal, ok := parseAgentFactValue(resp.Stdout, "dir")
	if !ok {
		t.Fatalf("dir: missing after core facts check")
	}
	if !filepath.IsAbs(dirVal) {
		t.Fatalf("dir: must be absolute, got %q", dirVal)
	}
	wantSuffix := taskDirName(req.TaskID, "create", req.Title) // e.g. 1-create-add-dark-mode
	rel := inboxTaskRel(req.TaskID, "create", req.Title)       // inbox/1-create-add-dark-mode
	dirSlash := filepath.ToSlash(dirVal)
	if !strings.HasSuffix(dirSlash, wantSuffix) && !strings.Contains(dirSlash, rel) {
		t.Fatalf("dir: %q must contain %q or suffix-match %q", dirVal, rel, wantSuffix)
	}
	// Align with on-disk task dir from index (strong equality without fixed home)
	wantAbs := findTaskDirByID(t, req, req.TaskID)
	if filepath.Clean(dirVal) != filepath.Clean(wantAbs) {
		t.Fatalf("dir: got %q want on-disk task dir %q", dirVal, wantAbs)
	}

	assertAgentNoAltPathKeys(t, resp.Stdout)
}
```
