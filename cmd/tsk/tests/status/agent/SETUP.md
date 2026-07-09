# Scenario

**Feature**: `tsk status --format=agent` renders a strict 2-row plain pipeline diagram plus machine-readable facts

```
# --format=agent (or --format agent): 2 art rows + id/stage/terminal/advance/next; no ANSI, no boxes
tsk status --format=agent <id> -> spine with [doing]/(name)/bare marks + back line (refine, questions) + facts
```

## Context

Shared helpers for agent-format diagram structure, node marks, facts block, and no-rectangle chrome. Prefer structural token checks over pixel-perfect spacing on row 2.

```go
import "regexp"

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
	_ = assertAgentAdvanceOK
	_ = assertAgentAdvanceBlocked
	_ = assertAgentNextMentions
	_ = assertAgentArtHasBackLine
	_ = assertAgentArtNoSatisfied
	_ = assertAgentNoANSI
}
```
