# Scenario

**Feature**: `tsk status --format=agent` renders a strict 2-row plain pipeline diagram plus machine-readable facts

```
# --format=agent (or --format agent): 2 art rows + id/title/stage/terminal/topic/dir/advance/next; no ANSI, no boxes
tsk status --format=agent <id> -> spine with [doing]/(name)/bare marks + back line (refine, questions) + facts
```

## Context

Shared helpers for agent-format diagram structure, node marks, facts block (id → title → stage → terminal → topic → dir), and no-rectangle chrome. Prefer structural token checks over pixel-perfect spacing on row 2.

```go
import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xhd2015/doctest/assert"
)

var agentSpineStages = []string{
	"create", "in_process", "clarification", "implementation",
	"verification", "summary", "done",
}

var agentStageTokenRE = regexp.MustCompile(
	`\(?\b(create|in_process|clarification|implementation|verification|summary|done)\b(?:\[doing\])?\)?`,
)

func Setup(t *testing.T, req *Request) error {
	ensureAgentHelpersUsed()
	return nil
}

func agentStatusArgs(id int, extra ...string) []string {
	args := []string{"--format=agent"}
	args = append(args, extra...)
	return statusArgs(id, args...)
}

func findAgentSpineLine(stdout string) string {
	for _, line := range stdoutLines(stdout) {
		plain := stripANSI(line)
		if strings.Contains(plain, "->") &&
			strings.Contains(plain, "create") &&
			strings.Contains(plain, "done") {
			return plain
		}
	}
	return ""
}

func agentArtText(stdout string) string {
	var b strings.Builder
	for _, line := range stdoutLines(stdout) {
		plain := stripANSI(line)
		if strings.Contains(plain, "->") ||
			strings.Contains(plain, "user_followup") ||
			strings.Contains(plain, "refine") ||
			strings.Contains(plain, "questions") ||
			strings.Contains(plain, "^") ||
			(strings.Contains(plain, "+") && strings.Contains(plain, "-")) {
			b.WriteString(plain)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func assertAgentNoRectChrome(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	for _, bad := range []string{"+---", "╭", "╰", "╮", "╯", "┌", "└", "┐", "┘"} {
		if strings.Contains(plain, bad) {
			t.Fatalf("agent format must not use rectangle/box chrome %q in:\n%s", bad, stdout)
		}
	}
	for _, stage := range pipelineStages {
		unicodeBox := "│ " + stage + " │"
		asciiBox := "| " + stage + " |"
		if strings.Contains(plain, unicodeBox) || strings.Contains(plain, asciiBox) {
			t.Fatalf("agent format must not box stage %q as %q / %q in:\n%s",
				stage, unicodeBox, asciiBox, stdout)
		}
	}
}

func assertAgentSpineOrder(t *testing.T, stdout string) {
	t.Helper()
	line := findAgentSpineLine(stdout)
	if line == "" {
		t.Fatalf("expected agent spine row with create … done joined by -> in:\n%s", stdout)
	}
	pos := 0
	for _, stage := range agentSpineStages {
		idx := strings.Index(line[pos:], stage)
		if idx < 0 {
			t.Fatalf("spine missing stage %q after pos %d on line %q\nstdout:\n%s",
				stage, pos, line, stdout)
		}
		pos += idx + len(stage)
	}
	// stages joined by ->
	if !strings.Contains(line, "->") {
		t.Fatalf("spine must join stages with -> : %q", line)
	}
	// rough check: count arrows between 7 spine stages
	arrows := strings.Count(line, "->")
	if arrows < 6 {
		t.Fatalf("expected at least 6 -> arrows on spine, got %d on %q", arrows, line)
	}
}

func assertAgentDoing(t *testing.T, stdout, stage string) {
	t.Helper()
	mark := stage + "[doing]"
	if !strings.Contains(stripANSI(stdout), mark) {
		t.Fatalf("expected current mark %q in:\n%s", mark, stdout)
	}
	if n := strings.Count(stripANSI(stdout), "[doing]"); n != 1 {
		t.Fatalf("expected exactly one [doing] mark, got %d in:\n%s", n, stdout)
	}
}

func assertAgentFuture(t *testing.T, stdout, stage string) {
	t.Helper()
	want := "(" + stage + ")"
	if !strings.Contains(stripANSI(stdout), want) {
		t.Fatalf("expected future mark %q in:\n%s", want, stdout)
	}
}

func assertAgentPastBare(t *testing.T, stdout, stage string) {
	t.Helper()
	line := findAgentSpineLine(stdout)
	if line == "" {
		t.Fatalf("missing spine line for past-bare check of %q:\n%s", stage, stdout)
	}
	if strings.Contains(line, stage+"[doing]") {
		t.Fatalf("past stage %q must not be [doing] on spine: %q", stage, line)
	}
	if strings.Contains(line, "("+stage+")") {
		t.Fatalf("past stage %q must not be future (name) on spine: %q", stage, line)
	}
	if !strings.Contains(line, stage) {
		t.Fatalf("past stage %q missing from spine: %q", stage, line)
	}
}

func assertAgentFact(t *testing.T, stdout, key, value string) {
	t.Helper()
	// allow flexible spacing after colon
	re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `\s*:\s*` + regexp.QuoteMeta(value) + `\s*$`)
	if !re.MatchString(stripANSI(stdout)) {
		// also accept inline without requiring whole-line strictness
		loose := key + ": " + value
		if !strings.Contains(stripANSI(stdout), loose) {
			t.Fatalf("expected fact %q (or line match) in:\n%s", loose, stdout)
		}
	}
}

func assertAgentHasFactKeys(t *testing.T, stdout string, keys ...string) {
	t.Helper()
	plain := stripANSI(stdout)
	for _, key := range keys {
		re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `\s*:`)
		if !re.MatchString(plain) && !strings.Contains(plain, key+":") {
			t.Fatalf("expected fact key %q in:\n%s", key, stdout)
		}
	}
}

// agentLeadingFacts returns the leading key: value facts block (lines before the
// first blank line). Agent format prints id/title/stage/terminal/topic/dir here, then art.
func agentLeadingFacts(stdout string) string {
	plain := stripANSI(stdout)
	var b strings.Builder
	for _, line := range strings.Split(plain, "\n") {
		if line == "" {
			break
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// assertAgentLeadingFactsShape matches the leading facts block via assert.Output.
// The dir: line is a runtime-exact literal (value taken from stdout), not __DIR__ type=string,
// because doctest's assert-mod matches type=string as non-greedy [^\n]*? which can consume
// empty and leave the path as "unparsed remainder". Use the exact stdout dir string — do not
// Clean/EvalSymlinks it before templating (macOS /var vs /private/var).
func assertAgentLeadingFactsShape(t *testing.T, stdout, title, stage, terminal, topic string) {
	t.Helper()
	facts := agentLeadingFacts(stdout)
	dirVal, ok := parseAgentFactValue(stdout, "dir")
	if !ok || dirVal == "" {
		t.Fatalf("dir: missing for leading facts template\n%s", stdout)
	}
	tmpl := fmt.Sprintf(`---
version: 2
__ID__: type=number, example=1
---
id: __ID__
title: %s
stage: %s
terminal: %s
topic: %s
dir: %s
`, title, stage, terminal, topic, dirVal)
	assert.Output(t, facts, tmpl)
}

// parseAgentFactValue returns the value for a leading whole-line key: value fact.
func parseAgentFactValue(stdout, key string) (string, bool) {
	plain := stripANSI(stdout)
	re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `\s*:\s*(.*?)\s*$`)
	m := re.FindStringSubmatch(plain)
	if m == nil {
		return "", false
	}
	return m[1], true
}

// assertAgentFactKeyOrder checks fact keys appear in order as whole-line keys.
func assertAgentFactKeyOrder(t *testing.T, stdout string, keys ...string) {
	t.Helper()
	plain := stripANSI(stdout)
	pos := 0
	for _, key := range keys {
		re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(key) + `\s*:`)
		loc := re.FindStringIndex(plain[pos:])
		if loc == nil {
			t.Fatalf("expected fact key %q after earlier keys in order %v; stdout:\n%s",
				key, keys, stdout)
		}
		pos += loc[1]
	}
}

// assertAgentNoAltPathKeys fails if agent facts use path: or path_rel: instead of dir:.
func assertAgentNoAltPathKeys(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	for _, bad := range []string{"path", "path_rel"} {
		re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(bad) + `\s*:`)
		if re.MatchString(plain) {
			t.Fatalf("agent facts must use dir: only, not %s:; stdout:\n%s", bad, stdout)
		}
	}
}

// agentInboxTopic is the locked topic: value for tasks with null topic_path (inbox).
// Differs from tsk show, which prints "inbox" for uncategorized tasks.
const agentInboxTopic = "(not classified yet)"

// assertAgentTopicFact checks topic: is always present after terminal and before
// dir, with the exact value (slash-joined path segments, or agentInboxTopic).
func assertAgentTopicFact(t *testing.T, stdout, wantTopic string) {
	t.Helper()
	assertAgentHasFactKeys(t, stdout, "topic")
	assertAgentFact(t, stdout, "topic", wantTopic)
	assertAgentFactKeyOrder(t, stdout, "terminal", "topic", "dir")
}

// assertAgentDirFact checks dir: is present after topic in key order, value is
// an absolute path, and it contains/suffix-matches the expected inbox relative
// segment for id/stage/title (homes and temp roots vary — no full-string hardcode).
// For topic-placed tasks, use assertAgentDirRel instead.
func assertAgentDirFact(t *testing.T, stdout string, id int, stage, title string) {
	t.Helper()
	assertAgentDirRel(t, stdout, inboxTaskRel(id, stage, title))
}

// assertAgentDirRel checks dir: after topic, absolute, matching a relative segment
// under TSK_HOME (e.g. inbox/… or topics/eng/backend/…).
func assertAgentDirRel(t *testing.T, stdout, wantRelSlash string) {
	t.Helper()
	assertAgentHasFactKeys(t, stdout, "dir")
	assertAgentFactKeyOrder(t, stdout, "topic", "dir")
	assertAgentNoAltPathKeys(t, stdout)

	dirVal, ok := parseAgentFactValue(stdout, "dir")
	if !ok || dirVal == "" {
		t.Fatalf("expected non-empty dir: fact in:\n%s", stdout)
	}
	if !filepath.IsAbs(dirVal) {
		t.Fatalf("dir: value must be absolute path, got %q", dirVal)
	}
	dirSlash := filepath.ToSlash(dirVal)
	if !strings.Contains(dirSlash, wantRelSlash) && !strings.HasSuffix(dirSlash, wantRelSlash) {
		t.Fatalf("dir: %q must contain or suffix-match relative task path %q\nstdout:\n%s",
			dirVal, wantRelSlash, stdout)
	}
}

// assertAgentCoreFacts requires leading facts id, title, stage, terminal, topic, dir
// with locked order id → title → stage → terminal → topic → dir.
// Topic value is agentInboxTopic (inbox / null topic_path). Dir checked via inbox path.
func assertAgentCoreFacts(t *testing.T, stdout string, id int, title, stage, terminal string) {
	t.Helper()
	idStr := fmt.Sprintf("%d", id)
	assertAgentFact(t, stdout, "id", idStr)
	assertAgentFact(t, stdout, "title", title)
	assertAgentFact(t, stdout, "stage", stage)
	assertAgentFact(t, stdout, "terminal", terminal)
	assertAgentTopicFact(t, stdout, agentInboxTopic)
	assertAgentFactKeyOrder(t, stdout, "id", "title", "stage", "terminal", "topic", "dir")
	assertAgentDirFact(t, stdout, id, stage, title)
}

func assertAgentAdvanceOK(t *testing.T, stdout, advanceTo string) {
	t.Helper()
	assertAgentFact(t, stdout, "advance", "ok")
	if advanceTo != "" {
		assertAgentFact(t, stdout, "advance_to", advanceTo)
	}
}

func assertAgentAdvanceBlocked(t *testing.T, stdout string) {
	t.Helper()
	assertAgentFact(t, stdout, "advance", "blocked")
}

func assertAgentNextMentions(t *testing.T, stdout string, subs ...string) {
	t.Helper()
	plain := stripANSI(stdout)
	// next: block should exist
	if !strings.Contains(plain, "next:") {
		t.Fatalf("expected next: block in:\n%s", stdout)
	}
	for _, sub := range subs {
		if !strings.Contains(plain, sub) {
			t.Fatalf("expected next guidance to mention %q in:\n%s", sub, stdout)
		}
	}
}

func assertAgentArtHasBackLine(t *testing.T, stdout string) {
	t.Helper()
	art := agentArtText(stdout)
	for _, need := range []string{"user_followup", "refine", "questions"} {
		if !strings.Contains(art, need) {
			t.Fatalf("expected agent art to contain %q (back line), art:\n%s\nfull:\n%s",
				need, art, stdout)
		}
	}
}

func assertAgentArtNoSatisfied(t *testing.T, stdout string) {
	t.Helper()
	art := agentArtText(stdout)
	if strings.Contains(art, "satisfied") {
		t.Fatalf("agent art must not draw satisfied edge; art:\n%s\nfull:\n%s", art, stdout)
	}
}

func assertAgentNoANSI(t *testing.T, resp *Response) {
	t.Helper()
	assertNoANSI(t, resp.Stdout)
	assertNoANSI(t, resp.Stderr)
}

func ensureAgentHelpersUsed() {
	_ = agentSpineStages
	_ = agentStageTokenRE
	_ = agentStatusArgs
	_ = findAgentSpineLine
	_ = agentArtText
	_ = assertAgentNoRectChrome
	_ = assertAgentSpineOrder
	_ = assertAgentDoing
	_ = assertAgentFuture
	_ = assertAgentPastBare
	_ = assertAgentFact
	_ = assertAgentHasFactKeys
	_ = agentLeadingFacts
	_ = assertAgentLeadingFactsShape
	_ = parseAgentFactValue
	_ = assertAgentFactKeyOrder
	_ = assertAgentNoAltPathKeys
	_ = agentInboxTopic
	_ = assertAgentTopicFact
	_ = assertAgentDirFact
	_ = assertAgentDirRel
	_ = assertAgentCoreFacts
	_ = assertAgentAdvanceOK
	_ = assertAgentAdvanceBlocked
	_ = assertAgentNextMentions
	_ = assertAgentArtHasBackLine
	_ = assertAgentArtNoSatisfied
	_ = assertAgentNoANSI
}
```
