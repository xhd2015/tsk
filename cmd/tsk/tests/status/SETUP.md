# Scenario

**Feature**: `tsk status` renders hand-made compact pipeline (default diagram) or agent plain view

```
# flags --format=diagram|agent, --color, --plain
# diagram (default): compact art (~34 col max, 3-line boxes) + optional ANSI
# agent: 2-row plain spine + back line + facts; no ANSI, no boxes
tsk status [--format=diagram|agent] [--color] [--plain] <id> -> pipeline view on stdout
```

## Context

Shared helpers for box-line assertions, width limits, ANSI checks, and advancing to terminal `done`. Agent-format helpers live under `status/agent/SETUP.md`.

```go
import "regexp"

var ansiEscapeRE = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiEscapeRE.ReplaceAllString(s, "")
}

func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	ensureStatusHelpersUsed()
	return nil
}

func maxLineWidth(s string) int {
	max := 0
	for _, line := range strings.Split(s, "\n") {
		if line == "" {
			continue
		}
		w := len([]rune(stripANSI(line)))
		if w > max {
			max = w
		}
	}
	return max
}

func assertMaxWidth36(t *testing.T, stdout string) {
	t.Helper()
	for i, line := range strings.Split(stdout, "\n") {
		if line == "" {
			continue
		}
		w := len([]rune(stripANSI(line)))
		if w > 36 {
			t.Fatalf("line %d visible width %d exceeds 36: %q", i+1, w, stripANSI(line))
		}
	}
}

func assertBoxLineForStage(t *testing.T, stdout, stage string) {
	t.Helper()
	unicodeBox := "│ " + stage + " │"
	asciiBox := "| " + stage + " |"
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, unicodeBox) || strings.Contains(line, asciiBox) {
			return
		}
	}
	t.Fatalf("expected box line for stage %q (│ %s │ or | %s |) in:\n%s", stage, stage, stage, stdout)
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("expected no ANSI escapes in output, got:\n%s", s)
	}
}

func assertHasANSISGR(t *testing.T, s, sgr string) {
	t.Helper()
	want := "\x1b[" + sgr
	if !strings.Contains(s, want) {
		t.Fatalf("expected ANSI SGR %q in output, got:\n%s", want, s)
	}
}

func assertStageLineHasGreen(t *testing.T, stdout, stage string) {
	t.Helper()
	for _, line := range strings.Split(stdout, "\n") {
		if strings.Contains(line, stage) && strings.Contains(line, "\x1b[32m") {
			return
		}
	}
	t.Fatalf("expected green ANSI (\\x1b[32m) on line containing %q:\n%s", stage, stdout)
}

func assertCompactBoxArt(t *testing.T, stdout string) {
	t.Helper()
	if !strings.Contains(stdout, "╭") && !strings.Contains(stdout, "+") {
		t.Fatalf("expected compact box art (╭ or +) in stdout:\n%s", stdout)
	}
}

func assertNotMermaidWide(t *testing.T, stdout string) {
	t.Helper()
	for _, bad := range []string{"stateDiagram", "stateDiagram-v2", "[*]"} {
		if strings.Contains(stdout, bad) {
			t.Fatalf("expected hand-made compact pipeline, not mermaid output containing %q:\n%s", bad, stdout)
		}
	}
}

func assertStdoutHasStages(t *testing.T, stdout string, stages ...string) {
	t.Helper()
	for _, stage := range stages {
		assertContains(t, stdout, stage)
	}
}

func assertStdoutEndsWithNewline(t *testing.T, stdout string) {
	t.Helper()
	if stdout == "" {
		t.Fatal("stdout should not be empty")
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout should end with newline, got %q", stdout)
	}
}

func assertStatusOK(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)
}

func advanceToDone(t *testing.T, req *Request, id int) {
	t.Helper()
	advanceTo(t, req, id, "summary")
	runTskOK(t, req, "done", fmt.Sprintf("%d", id))
}

func statusArgs(id int, extra ...string) []string {
	args := append([]string{"status"}, extra...)
	args = append(args, fmt.Sprintf("%d", id))
	return args
}

var pipelineStages = []string{
	"create", "in_process", "clarification", "implementation",
	"verification", "summary", "user_followup", "done",
}

func assertAllStagesBoxed(t *testing.T, stdout string) {
	t.Helper()
	for _, stage := range pipelineStages {
		assertBoxLineForStage(t, stdout, stage)
	}
}

func countSubstring(s, substr string) int {
	return strings.Count(s, substr)
}

func assertContainsArrowDown(t *testing.T, stdout string, minCount int) {
	t.Helper()
	got := countSubstring(stdout, "▼")
	if got < minCount {
		t.Fatalf("expected at least %d ▼ arrows, got %d in:\n%s", minCount, got, stdout)
	}
}

func stdoutLines(stdout string) []string {
	return strings.Split(stdout, "\n")
}

func firstLineIndexContaining(lines []string, substrings ...string) int {
	for i, line := range lines {
		plain := stripANSI(line)
		for _, sub := range substrings {
			if strings.Contains(plain, sub) {
				return i
			}
		}
	}
	return -1
}

func lastLineIndexContaining(lines []string, substr string) int {
	last := -1
	for i, line := range lines {
		if strings.Contains(stripANSI(line), substr) {
			last = i
		}
	}
	return last
}

func boxLineIndex(lines []string, stage string) int {
	unicodeBox := "│ " + stage + " │"
	asciiBox := "| " + stage + " |"
	return firstLineIndexContaining(lines, unicodeBox, asciiBox)
}

func assertEdgeLabelOrder(t *testing.T, stdout string) {
	t.Helper()
	lines := stdoutLines(stdout)

	createBox := boxLineIndex(lines, "create")
	claim := firstLineIndexContaining(lines, "claim")
	if createBox < 0 || claim < 0 || claim <= createBox {
		t.Fatalf("claim should follow create box (create line %d, claim line %d):\n%s",
			createBox+1, claim+1, stdout)
	}

	inProcessBox := boxLineIndex(lines, "in_process")
	clarificationBox := boxLineIndex(lines, "clarification")
	research := firstLineIndexContaining(lines, "research")
	if inProcessBox < 0 || clarificationBox < 0 || research < 0 ||
		research <= inProcessBox || research >= clarificationBox {
		t.Fatalf("research should appear after in_process and before clarification (in_process %d, research %d, clarification %d):\n%s",
			inProcessBox+1, research+1, clarificationBox+1, stdout)
	}

	implementationBox := boxLineIndex(lines, "implementation")
	confirmed := firstLineIndexContaining(lines, "confirmed")
	if clarificationBox < 0 || implementationBox < 0 || confirmed < 0 ||
		confirmed <= clarificationBox || confirmed >= implementationBox {
		t.Fatalf("confirmed should appear after clarification and before implementation (clarification %d, confirmed %d, implementation %d):\n%s",
			clarificationBox+1, confirmed+1, implementationBox+1, stdout)
	}

	summaryBox := boxLineIndex(lines, "summary")
	questions := firstLineIndexContaining(lines, "questions")
	if summaryBox < 0 || questions < 0 || questions <= summaryBox {
		t.Fatalf("questions should follow summary box (summary %d, questions %d):\n%s",
			summaryBox+1, questions+1, stdout)
	}

	verificationBox := boxLineIndex(lines, "verification")
	followupBox := boxLineIndex(lines, "user_followup")
	satisfied := firstLineIndexContaining(lines, "satisfied")
	if summaryBox < 0 || verificationBox < 0 || followupBox < 0 || satisfied < 0 {
		t.Fatalf("missing stage box or satisfied label in:\n%s", stdout)
	}
	if satisfied <= summaryBox {
		t.Fatalf("satisfied should follow summary fork (summary %d, satisfied %d):\n%s",
			summaryBox+1, satisfied+1, stdout)
	}
	if satisfied < verificationBox {
		t.Fatalf("satisfied should not appear before verification (verification %d, satisfied %d):\n%s",
			verificationBox+1, satisfied+1, stdout)
	}
	if satisfied < followupBox-1 {
		t.Fatalf("satisfied should appear near user_followup→done merge (followup %d, satisfied %d):\n%s",
			followupBox+1, satisfied+1, stdout)
	}
}

func assertForkSemantics(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	if !strings.Contains(plain, "no followup") {
		t.Fatalf("expected no followup label on summary→done branch in:\n%s", stdout)
	}
	for i, line := range stdoutLines(stdout) {
		p := stripANSI(line)
		if strings.Contains(p, "no followup") && strings.Contains(p, "questions") {
			t.Fatalf("no followup and questions must not share a line (line %d): %q", i+1, p)
		}
	}
	satisfiedLine := ""
	for _, line := range stdoutLines(stdout) {
		if strings.Contains(stripANSI(line), "satisfied") {
			satisfiedLine = stripANSI(line)
			break
		}
	}
	if satisfiedLine == "" {
		t.Fatalf("missing satisfied label in:\n%s", stdout)
	}
	if !strings.Contains(satisfiedLine, "►") && !strings.Contains(satisfiedLine, ">") {
		t.Fatalf("satisfied merge should include branch arrow (► or >), got %q in:\n%s", satisfiedLine, stdout)
	}
	if strings.Contains(plain, "╰──▼") {
		t.Fatalf("done box bottom must not embed ▼ (corrupts box) in:\n%s", stdout)
	}
}

func assertFollowupBeforeTerminal(t *testing.T, stdout string) {
	t.Helper()
	lines := stdoutLines(stdout)
	followupIdx := boxLineIndex(lines, "user_followup")
	if followupIdx < 0 {
		followupIdx = firstLineIndexContaining(lines, "user_followup")
	}
	terminalIdx := lastLineIndexContaining(lines, "◉")
	if followupIdx < 0 {
		t.Fatalf("stdout missing user_followup in:\n%s", stdout)
	}
	if terminalIdx < 0 {
		t.Fatalf("stdout missing terminal ◉ in:\n%s", stdout)
	}
	if followupIdx > terminalIdx {
		t.Fatalf("user_followup (line %d) must appear before terminal ◉ (line %d):\n%s",
			followupIdx+1, terminalIdx+1, stdout)
	}
	for i := terminalIdx + 1; i < len(lines); i++ {
		plain := strings.TrimSpace(stripANSI(lines[i]))
		if plain == "" {
			continue
		}
		if strings.Contains(plain, "user_followup") {
			t.Fatalf("orphan user_followup after terminal ◉ at line %d: %q", i+1, plain)
		}
	}
}

func ensureStatusHelpersUsed() {
	_ = ansiEscapeRE
	_ = stripANSI
	_ = maxLineWidth
	_ = assertMaxWidth36
	_ = assertBoxLineForStage
	_ = assertNoANSI
	_ = assertHasANSISGR
	_ = assertStageLineHasGreen
	_ = assertCompactBoxArt
	_ = assertNotMermaidWide
	_ = assertStdoutHasStages
	_ = assertStdoutEndsWithNewline
	_ = assertStatusOK
	_ = advanceToDone
	_ = statusArgs
	_ = pipelineStages
	_ = assertAllStagesBoxed
	_ = countSubstring
	_ = assertContainsArrowDown
	_ = stdoutLines
	_ = firstLineIndexContaining
	_ = lastLineIndexContaining
	_ = boxLineIndex
	_ = assertEdgeLabelOrder
	_ = assertFollowupBeforeTerminal
	_ = assertForkSemantics
}
```