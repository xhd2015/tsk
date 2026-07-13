package pipeline

import (
	"strings"

	"github.com/xhd2015/tsk/tskcli/storage"
)

const (
	ansiReset     = "\x1b[0m"
	ansiGreenBold = "\x1b[1m\x1b[32m"
	ansiGrey      = "\x1b[90m"
	ansiOrange    = "\x1b[33m"
)

var pipelineStages = []string{
	"user_followup",
	"clarification",
	"implementation",
	"verification",
	"in_process",
	"summary",
	"create",
	"done",
}

// Highlight applies semantic ANSI colors to a rendered pipeline diagram.
func Highlight(rendered string, task storage.Task, color bool) string {
	lines := splitLines(rendered)
	if !color {
		return joinLines(lines)
	}

	visited := visitedStages(task)
	boxes := findStageBoxes(lines)

	for stage, box := range boxes {
		switch {
		case stage == task.Stage:
			colorBoxLines(lines, box.start, box.end, ansiGreenBold)
		case visited[stage]:
			colorBoxLines(lines, box.start, box.end, ansiGrey)
		}
	}

	if box, ok := boxes[task.Stage]; ok && box.start > 0 {
		colorEdgeIntoCurrent(lines, edgeLineAboveBox(lines, box.start))
	}

	if task.Stage == "done" {
		colorDoneMarker(lines)
	}

	return joinLines(lines)
}

func splitLines(rendered string) []string {
	lines := strings.Split(rendered, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func joinLines(lines []string) string {
	if len(lines) == 0 {
		return "\n"
	}
	return strings.Join(lines, "\n") + "\n"
}

func visitedStages(task storage.Task) map[string]bool {
	visited := make(map[string]bool)
	for _, entry := range task.StageHistory {
		if entry.To != task.Stage {
			visited[entry.To] = true
		}
	}
	return visited
}

type stageBox struct {
	start int
	end   int
}

func findStageBoxes(lines []string) map[string]stageBox {
	boxes := make(map[string]stageBox)
	for i, line := range lines {
		stage := stageNameInBoxLine(line)
		if stage == "" {
			continue
		}
		start, end := boxBounds(lines, i)
		boxes[stage] = stageBox{start: start, end: end}
	}
	return boxes
}

func boxBounds(lines []string, labelIdx int) (start, end int) {
	start = labelIdx
	for start > 0 && !isBoxTop(lines[start]) {
		start--
	}
	end = labelIdx
	for end < len(lines)-1 && !isBoxBottom(lines[end]) {
		end++
	}
	return start, end
}

func stageNameInBoxLine(line string) string {
	if !strings.Contains(line, "│") && !strings.Contains(line, "|") {
		return ""
	}

	segments := strings.Split(line, "│")
	if len(segments) == 1 {
		segments = strings.Split(line, "|")
	}
	for i := len(segments) - 1; i >= 0; i-- {
		seg := strings.TrimSpace(segments[i])
		if seg == "" || isConnectorOnly(seg) {
			continue
		}
		for _, stage := range pipelineStages {
			if segmentMatchesStage(seg, stage) {
				return stage
			}
		}
	}
	return ""
}

func segmentMatchesStage(seg, stage string) bool {
	if seg == stage {
		return true
	}
	return strings.Contains(" "+seg+" ", " "+stage+" ")
}

func isBoxTop(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.Contains(line, "╭") ||
		(strings.HasPrefix(trimmed, "+") && strings.Contains(line, "-"))
}

func isBoxBottom(line string) bool {
	return strings.Contains(line, "╰") ||
		(strings.Contains(line, "+") && strings.Contains(line, "-") && strings.Count(line, "+") >= 2)
}

func colorBoxLines(lines []string, start, end int, prefix string) {
	for i := start; i <= end && i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "" {
			continue
		}
		spanStart, spanEnd, ok := stageBoxSpan(lines[i])
		if !ok {
			// Fallback: avoid coloring if we cannot isolate the box.
			continue
		}
		line := lines[i]
		lines[i] = line[:spanStart] + prefix + line[spanStart:spanEnd] + ansiReset + line[spanEnd:]
	}
}

// stageBoxSpan returns the byte range [start, end) of the stage box on a line,
// excluding leading left-rail and trailing right-rail connectors.
func stageBoxSpan(line string) (start, end int, ok bool) {
	if i := strings.Index(line, "╭"); i >= 0 {
		if j := strings.Index(line[i:], "╮"); j >= 0 {
			return i, i + j + len("╮"), true
		}
	}
	if i := strings.Index(line, "╰"); i >= 0 {
		if j := strings.Index(line[i:], "╯"); j >= 0 {
			return i, i + j + len("╯"), true
		}
	}

	if stage := stageNameInBoxLine(line); stage != "" {
		if s, e, found := midBoxSpan(line, stage); found {
			return s, e, true
		}
	}

	if s, e, found := asciiBoxChromeSpan(line); found {
		return s, e, true
	}
	return 0, 0, false
}

func midBoxSpan(line, stage string) (start, end int, ok bool) {
	runes := []rune(line)
	stageRunes := []rune(stage)
	nameAt := -1
	for i := 0; i <= len(runes)-len(stageRunes); i++ {
		match := true
		for k := 0; k < len(stageRunes); k++ {
			if runes[i+k] != stageRunes[k] {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		// Prefer a boxed occurrence (spaces/walls around the stage name).
		leftOK := i == 0 || runes[i-1] == ' ' || isBoxWall(runes[i-1])
		right := i + len(stageRunes)
		rightOK := right >= len(runes) || runes[right] == ' ' || isBoxWall(runes[right])
		if leftOK && rightOK {
			nameAt = i
			// Keep scanning so the last boxed match wins (mirrors stageNameInBoxLine).
		}
	}
	if nameAt < 0 {
		return 0, 0, false
	}

	left := nameAt - 1
	for left >= 0 && runes[left] == ' ' {
		left--
	}
	if left < 0 || !isBoxWall(runes[left]) {
		return 0, 0, false
	}

	right := nameAt + len(stageRunes)
	for right < len(runes) && runes[right] == ' ' {
		right++
	}
	if right >= len(runes) || !isBoxWall(runes[right]) {
		return 0, 0, false
	}

	start = len(string(runes[:left]))
	end = len(string(runes[:right+1]))
	return start, end, true
}

func isBoxWall(r rune) bool {
	switch r {
	case '│', '┤', '├', '|', '+':
		return true
	default:
		return false
	}
}

// asciiBoxChromeSpan finds a top/bottom ASCII box border (+-+ / +--+--+).
func asciiBoxChromeSpan(line string) (start, end int, ok bool) {
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		if runes[i] != '+' {
			continue
		}
		if i+1 >= len(runes) || runes[i+1] != '-' {
			continue
		}
		j := i + 1
		for j < len(runes) && (runes[j] == '-' || runes[j] == '+') {
			j++
		}
		if j-1 > i && runes[j-1] == '+' {
			start = len(string(runes[:i]))
			end = len(string(runes[:j]))
			return start, end, true
		}
	}
	return 0, 0, false
}

func edgeLineAboveBox(lines []string, boxStart int) int {
	idx := boxStart - 1
	for idx >= 0 && isConnectorOnly(lines[idx]) {
		idx--
	}
	return idx
}

func colorEdgeIntoCurrent(lines []string, idx int) {
	if idx < 0 || idx >= len(lines) {
		return
	}
	line := lines[idx]
	if strings.TrimSpace(line) == "" || isConnectorOnly(line) {
		return
	}
	lines[idx] = ansiOrange + line + ansiReset
}

func isConnectorOnly(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	for _, r := range trimmed {
		switch r {
		case ' ', '│', '|', '─', '-', '┼', '├', '┤', '┬', '┴', '▼', '▲', '◄', '►', '╯', '╰', '╭', '┌', '┐', '└', '┘', '*', 'v', '^', '<', '>', '+':
			continue
		default:
			return false
		}
	}
	return true
}

func colorDoneMarker(lines []string) {
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "◉") || strings.Contains(lines[i], "@") {
			lines[i] = ansiGreenBold + lines[i] + ansiReset
			return
		}
	}
}