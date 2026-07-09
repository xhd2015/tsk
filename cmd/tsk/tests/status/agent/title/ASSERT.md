## Expected Output

Leading facts block (before blank line / art):

```
---
version: 2
__ID__: type=number, example=1, task id from create
__DIR__: type=string, example=/tmp/work/.tsk/inbox/1-create-add-dark-mode, absolute task dir
---
id: __ID__
title: add dark mode
stage: create
terminal: false
topic: (not classified yet)
dir: __DIR__
```

## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Leading facts order locked: `id` → `title` → `stage` → `terminal` → `topic` → `dir`.
- `title:` value is the exact create title (`add dark mode`), same key as `tsk show`.
- Inbox `topic: (not classified yet)` after `terminal:`, before `dir:` (not `topic: inbox` as in `tsk show`).
- `dir:` absolute task directory path after `topic:` (see `agent/dir`).
- No slug/labels keys required in agent facts for this change.
- Existing agent chrome: spine still present; no ANSI; no rectangle boxes.

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
	assertAgentSpineOrder(t, resp.Stdout)

	// Strict leading facts block: id → title → stage → terminal → topic → dir
	assert.Output(t, agentLeadingFacts(resp.Stdout), `---
version: 2
__ID__: type=number, example=1, task id from create
__DIR__: type=string, example=/tmp/work/.tsk/inbox/1-create-add-dark-mode, absolute task dir
---
id: __ID__
title: add dark mode
stage: create
terminal: false
topic: (not classified yet)
dir: __DIR__
`)

	// Cross-check helpers (value + key order) against full stdout
	assertAgentCoreFacts(t, resp.Stdout, req.TaskID, req.Title, "create", "false")
	assertAgentFact(t, resp.Stdout, "title", "add dark mode")
}
```
