# Scenario

**Feature**: `tsk status` renders hand-made compact pipeline or agent plain view (auto or explicit)

```
# flags --format=diagram|agent, --color, --plain; auto agent when host detected and no format flags
# diagram (default outside agents): compact art (~40 col, 3-line boxes) + optional ANSI
# agent: 2-row plain spine + back line + facts (id/title/stage/terminal/topic/dir); no ANSI, no boxes
# auto: bare status + detect/TSK_STATUS_FORMAT → agent; --format/--color/--plain always win
# geometry (diagram): left-rail refine; right-rail no-followup (┐/│/┘ same col); vertical satisfied; done→◉ dead end
# highlight: color only stage box span — left refine rail outside box SGR (see color-box-only)
# golden: status/diagram-golden + status/plain-golden expected.txt exact stdout
tsk status [--format=diagram|agent] [--color] [--plain] <id> -> pipeline view on stdout
```

## Context

Shared helpers for box-line assertions, width limits, ANSI checks, golden stdout files, no-followup rail column alignment, box-only highlight, and advancing to terminal `done`. Agent-format helpers live under `status/agent/SETUP.md`. Auto-format selection leaves live under `status/auto-format/` (inject host env via `ExtraEnv`; root `tskEnv` clears agent vars for stable diagram defaults). Exact diagram geometry is sealed by `diagram-golden` / `plain-golden` `expected.txt` fixtures.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/xhd2015/doctest/session"
)

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

func assertMaxWidth(t *testing.T, stdout string, max int) {
	t.Helper()
	for i, line := range strings.Split(stdout, "\n") {
		if line == "" {
			continue
		}
		w := len([]rune(stripANSI(line)))
		if w > max {
			t.Fatalf("line %d visible width %d exceeds %d: %q", i+1, w, max, stripANSI(line))
		}
	}
}

// assertMaxWidth42: new geometry is ~40 cols (was 36); soft cap only — goldens are source of truth.
func assertMaxWidth42(t *testing.T, stdout string) {
	t.Helper()
	assertMaxWidth(t, stdout, 42)
}

func isStageBoxMidRow(line, stage string) bool {
	plain := stripANSI(line)
	u := regexp.MustCompile(`[│┤]\s*` + regexp.QuoteMeta(stage) + `\s*[│├]`)
	a := regexp.MustCompile(`[|+]\s*` + regexp.QuoteMeta(stage) + `\s*[|+]`)
	return u.MatchString(plain) || a.MatchString(plain)
}

func assertBoxLineForStage(t *testing.T, stdout, stage string) {
	t.Helper()
	for _, line := range strings.Split(stdout, "\n") {
		if isStageBoxMidRow(line, stage) {
			return
		}
	}
	t.Fatalf("expected box mid-row for stage %q (│/┤ %s │/├ or ASCII) in:\n%s", stage, stage, stdout)
}

func assertStdoutEqualsFile(t *testing.T, d *session.Doctest, stdout, relPath string) {
	t.Helper()
	path := relPath
	if !filepath.IsAbs(path) {
		path = filepath.Join(d.DOCTEST_CASE, relPath)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	want := string(data)
	if stdout != want {
		t.Fatalf("stdout != %s\n--- got (%d bytes) ---\n%s\n--- want (%d bytes) ---\n%s",
			path, len(stdout), stdout, len(want), want)
	}
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

// lastRuneIndex returns the last index of r in s, or -1.
func lastRuneIndex(s string, r rune) int {
	runes := []rune(s)
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] == r {
			return i
		}
	}
	return -1
}

// assertNoFollowupRailAligned checks that the no-followup right rail is column-aligned:
// corner ┐ (or ASCII + after "no followup"), vertical │/| on intervening lines, and
// done-mid join ┘ (or ASCII +) share the same column.
func assertNoFollowupRailAligned(t *testing.T, stdout string) {
	t.Helper()
	lines := stdoutLines(stdout)
	cornerCol := -1
	joinCol := -1
	cornerLine := -1
	joinLine := -1

	for i, line := range lines {
		plain := stripANSI(line)
		if strings.Contains(plain, "no followup") {
			// unicode ┐ or plain trailing + on the horizontal branch
			col := lastRuneIndex(plain, '┐')
			if col < 0 {
				col = lastRuneIndex(plain, '+')
			}
			if col < 0 {
				t.Fatalf("no-followup line missing corner ┐/+: %q", plain)
			}
			cornerCol = col
			cornerLine = i
		}
		if isStageBoxMidRow(line, "done") {
			col := lastRuneIndex(plain, '┘')
			if col < 0 {
				col = lastRuneIndex(plain, '+')
			}
			if col < 0 {
				t.Fatalf("done mid missing join ┘/+: %q", plain)
			}
			joinCol = col
			joinLine = i
		}
	}
	if cornerCol < 0 {
		t.Fatalf("missing no-followup corner ┐/+ in:\n%s", stdout)
	}
	if joinCol < 0 {
		t.Fatalf("missing done mid join ┘/+ in:\n%s", stdout)
	}
	if cornerCol != joinCol {
		t.Fatalf("no-followup rail misaligned: corner col %d (line %d) != done join col %d (line %d):\n%s",
			cornerCol, cornerLine+1, joinCol, joinLine+1, stdout)
	}

	// vertical right-rail cells between corner and done mid
	for i := cornerLine + 1; i < joinLine; i++ {
		plain := stripANSI(lines[i])
		runes := []rune(plain)
		if len(runes) <= cornerCol {
			t.Fatalf("right-rail line %d shorter than corner col %d: %q", i+1, cornerCol, plain)
		}
		r := runes[cornerCol]
		if r != '│' && r != '|' {
			t.Fatalf("expected right rail │/| at col %d on line %d, got %q: %q",
				cornerCol, i+1, string(r), plain)
		}
	}
}

// firstVisibleUnderActiveSGR reports whether the first non-ANSI content rune of line
// is reached while an SGR color/style is already active (whole-line coloring).
func firstVisibleUnderActiveSGR(line string) bool {
	active := false
	for i := 0; i < len(line); {
		if line[i] == '\x1b' && i+1 < len(line) && line[i+1] == '[' {
			j := i + 2
			for j < len(line) && ((line[j] >= '0' && line[j] <= '9') || line[j] == ';') {
				j++
			}
			if j >= len(line) || line[j] != 'm' {
				return active // malformed; treat conservatively
			}
			code := line[i+2 : j]
			if code == "0" || code == "" {
				active = false
			} else {
				// any non-reset SGR (bold, green, grey, …)
				active = true
			}
			i = j + 1
			continue
		}
		// first visible content
		return active
	}
	return false
}

// stageBoxMidLine returns the raw (possibly ANSI) mid-row line for stage, or "".
func stageBoxMidLine(stdout, stage string) string {
	for _, line := range strings.Split(stdout, "\n") {
		if isStageBoxMidRow(line, stage) {
			return line
		}
	}
	return ""
}

// assertBoxColoredLeftRailClear: current stage box still has green on the stage
// mid-row, but the leading left refine-rail │ must not sit inside the same SGR
// span as the box (whole-line colorBoxLines is the bug).
func assertBoxColoredLeftRailClear(t *testing.T, stdout, stage string) {
	t.Helper()
	line := stageBoxMidLine(stdout, stage)
	if line == "" {
		t.Fatalf("expected box mid-row for stage %q in:\n%s", stage, stdout)
	}
	if !strings.Contains(line, "\x1b[32m") {
		t.Fatalf("expected green ANSI (\\x1b[32m) on %q box mid-row:\n%s", stage, line)
	}
	plain := stripANSI(line)
	// stages that share the left refine rail begin with │ (unicode) or | (plain)
	if !strings.HasPrefix(strings.TrimLeft(plain, " \t"), "│") &&
		!strings.HasPrefix(strings.TrimLeft(plain, " \t"), "|") {
		t.Fatalf("expected stage %q mid-row to start with left refine rail │/|, got %q", stage, plain)
	}
	if firstVisibleUnderActiveSGR(line) {
		t.Fatalf("left-rail │ must not be inside box SGR on stage %q mid-row (whole-line color is wrong):\n%q\nfull:\n%s",
			stage, line, stdout)
	}
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
	for i, line := range lines {
		if isStageBoxMidRow(line, stage) {
			return i
		}
	}
	return -1
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
	// vertical satisfied sits below user_followup (toward done), not on/above the followup box
	if satisfied <= followupBox {
		t.Fatalf("satisfied should appear below user_followup box (followup %d, satisfied %d):\n%s",
			followupBox+1, satisfied+1, stdout)
	}
}

func assertForkSemantics(t *testing.T, stdout string) {
	t.Helper()
	plain := stripANSI(stdout)
	if !strings.Contains(plain, "no followup") {
		t.Fatalf("expected no followup label on summary→done branch in:\n%s", stdout)
	}
	if !strings.Contains(plain, "refine") {
		t.Fatalf("expected refine label on user_followup→clarification loop in:\n%s", stdout)
	}
	for i, line := range stdoutLines(stdout) {
		p := stripANSI(line)
		if strings.Contains(p, "no followup") && strings.Contains(p, "questions") {
			t.Fatalf("no followup and questions must not share a line (line %d): %q", i+1, p)
		}
	}
	// satisfied is a vertical spine label (like claim), not a sideways satisfied► decoration
	if strings.Contains(plain, "satisfied►") || strings.Contains(plain, "satisfied>") ||
		strings.Contains(plain, "-satisfied>") || strings.Contains(plain, "─satisfied►") {
		t.Fatalf("satisfied must not use sideways branch decoration (satisfied►); got:\n%s", stdout)
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
	// vertical label sits on spine (│ or |) with the word satisfied — no horizontal arrow on that line
	if strings.Contains(satisfiedLine, "►") || strings.Contains(satisfiedLine, ">") {
		t.Fatalf("satisfied line should be vertical spine label without ►/>, got %q in:\n%s", satisfiedLine, stdout)
	}
	if !strings.Contains(satisfiedLine, "│") && !strings.Contains(satisfiedLine, "|") {
		t.Fatalf("satisfied should sit on vertical spine (│ or |), got %q in:\n%s", satisfiedLine, stdout)
	}
	if strings.Contains(plain, "╰──▼") {
		t.Fatalf("done box bottom must not embed ▼ (corrupts box) in:\n%s", stdout)
	}
	// done→◉ is a dead end: no refine rail continuing under terminal
	lines := stdoutLines(stdout)
	termIdx := lastLineIndexContaining(lines, "◉")
	if termIdx < 0 {
		termIdx = lastLineIndexContaining(lines, "@")
	}
	if termIdx >= 0 {
		for i := termIdx + 1; i < len(lines); i++ {
			p := strings.TrimSpace(stripANSI(lines[i]))
			if p == "" {
				continue
			}
			if strings.Contains(p, "refine") || strings.Contains(p, "clarification") {
				t.Fatalf("terminal is a dead end; unexpected content after ◉/@ at line %d: %q", i+1, p)
			}
		}
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
	_ = assertMaxWidth
	_ = assertMaxWidth42
	_ = isStageBoxMidRow
	_ = assertBoxLineForStage
	_ = assertStdoutEqualsFile
	_ = assertNoANSI
	_ = assertHasANSISGR
	_ = assertStageLineHasGreen
	_ = lastRuneIndex
	_ = assertNoFollowupRailAligned
	_ = firstVisibleUnderActiveSGR
	_ = stageBoxMidLine
	_ = assertBoxColoredLeftRailClear
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