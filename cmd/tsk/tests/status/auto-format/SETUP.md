# Scenario

**Feature**: bare `tsk status` auto-selects agent vs diagram from host detection and `TSK_STATUS_FORMAT`

```
# precedence: --format > --color|--plain > TSK_STATUS_FORMAT > agentrunner.Detect > diagram
# root tskEnv clears CODEX_THREAD_ID / PI_CODING_AGENT / TSK_STATUS_FORMAT; leaves set ExtraEnv
tsk status <id>  (+ env / flags) -> agent facts block  OR  diagram box art
```

## Preconditions

- Parent `tskEnv` strips host-agent and format-override env so bare diagram is stable.
- Leaves re-inject via `req.ExtraEnv` (`CODEX_THREAD_ID=…`, `PI_CODING_AGENT=…`, `TSK_STATUS_FORMAT=…`).
- No production implementation yet: expect RED until auto-format lands.

## Context

Helpers distinguish agent fact/spine output from diagram box art without depending on `status/agent/` helpers (sibling branch). All leaves create a task at `create` and run `status` with the scenario's env/flags.

```go
import (
	"regexp"
	"strings"
)
var (
	agentFactIDRE    = regexp.MustCompile(`(?m)^id\s*:`)
	agentFactTitleRE = regexp.MustCompile(`(?m)^title\s*:`)
	agentFactTopicRE = regexp.MustCompile(`(?m)^topic\s*:`)
	agentFactDirRE   = regexp.MustCompile(`(?m)^dir\s*:`)
)

func Setup(t *testing.T, req *Request) error {
	ensureAutoFormatHelpersUsed()
	return nil
}

// createForAutoFormat creates an inbox task and stores id/title on req.
func createForAutoFormat(t *testing.T, req *Request, title string) int {
	t.Helper()
	req.Title = title
	id := createTask(t, req, title, "", nil)
	req.TaskID = id
	return id
}

// setStatusEnv appends KEY=value entries for the child tsk process.
func setStatusEnv(req *Request, kvs ...string) {
	req.ExtraEnv = append(req.ExtraEnv, kvs...)
}

func assertAgentStatusFormat(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	if !agentFactIDRE.MatchString(plain) {
		t.Fatalf("expected agent format fact id: in:\n%s", stdout)
	}
	if !agentFactTitleRE.MatchString(plain) {
		t.Fatalf("expected agent format fact title: in:\n%s", stdout)
	}
	if !agentFactTopicRE.MatchString(plain) {
		t.Fatalf("expected agent format fact topic: in:\n%s", stdout)
	}
	if !agentFactDirRE.MatchString(plain) {
		t.Fatalf("expected agent format fact dir: in:\n%s", stdout)
	}
	// spine join marker present
	if !(strings.Contains(plain, "->") && strings.Contains(plain, "create") && strings.Contains(plain, "done")) {
		t.Fatalf("expected agent spine (create … done joined by ->) in:\n%s", stdout)
	}
	// no rectangle / stage box chrome
	for _, bad := range []string{"+---", "╭", "╰", "╮", "╯", "┌", "└", "┐", "┘"} {
		if strings.Contains(plain, bad) {
			t.Fatalf("agent format must not use box chrome %q in:\n%s", bad, stdout)
		}
	}
	for _, stage := range []string{"create", "done"} {
		unicodeBox := "│ " + stage + " │"
		asciiBox := "| " + stage + " |"
		if strings.Contains(plain, unicodeBox) || strings.Contains(plain, asciiBox) {
			t.Fatalf("agent format must not box stage %q in:\n%s", stage, stdout)
		}
	}
}

func assertDiagramStatusFormat(t *testing.T, stdout string) {
	t.Helper()
	assertCompactBoxArt(t, stdout)
	assertBoxLineForStage(t, stdout, "create")
	assertNotMermaidWide(t, stdout)
	// not the agent leading facts block
	assertNotAgentFactBlock(t, stdout)
}

func assertNotAgentFactBlock(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	// agent always prints whole-line title: and id: facts; diagram must not
	if agentFactIDRE.MatchString(plain) && agentFactTitleRE.MatchString(plain) {
		t.Fatalf("expected diagram path without agent id:/title: fact block, got:\n%s", stdout)
	}
}

func assertAutoFormatOK(t *testing.T, resp *Response, err error) {
	t.Helper()
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
}

func ensureAutoFormatHelpersUsed() {
	_ = agentFactIDRE
	_ = agentFactTitleRE
	_ = agentFactTopicRE
	_ = agentFactDirRE
	_ = createForAutoFormat
	_ = setStatusEnv
	_ = assertAgentStatusFormat
	_ = assertDiagramStatusFormat
	_ = assertNotAgentFactBlock
	_ = assertAutoFormatOK
}
```
