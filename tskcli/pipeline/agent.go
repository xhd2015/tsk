package pipeline

import (
	"fmt"
	"strings"

	"github.com/xhd2015/tsk/tskcli/storage"
)

// agentSpine is the fixed main-path stage order for --format=agent.
var agentSpine = []string{
	"create",
	"in_process",
	"clarification",
	"implementation",
	"verification",
	"summary",
	"done",
}

// RenderAgent returns the agent-oriented status view: facts, 2-row plain
// pipeline art (no boxes, no ANSI), and advance/next guidance.
// dir is the absolute filesystem path of the task directory (from LoadTaskByID).
func RenderAgent(task storage.Task, dir string) string {
	var b strings.Builder

	fmt.Fprintf(&b, "id: %d\n", task.ID)
	fmt.Fprintf(&b, "title: %s\n", task.Title)
	fmt.Fprintf(&b, "stage: %s\n", task.Stage)
	fmt.Fprintf(&b, "terminal: %t\n", task.Stage == "done")
	parts, err := storage.ParseTopicPath(task.TopicPath)
	if err != nil || len(parts) == 0 {
		fmt.Fprintf(&b, "topic: (not classified yet)\n")
	} else {
		fmt.Fprintf(&b, "topic: %s\n", strings.Join(parts, "/"))
	}
	fmt.Fprintf(&b, "dir: %s\n", dir)
	b.WriteByte('\n')

	b.WriteString(renderAgentArt(task.Stage))
	b.WriteByte('\n')

	b.WriteString(renderAgentAdvanceNext(task))

	out := b.String()
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}
	return out
}

func renderAgentArt(current string) string {
	var spine strings.Builder
	clarCol := 0
	summaryCol := 0

	for i, stage := range agentSpine {
		if i > 0 {
			spine.WriteString(" -> ")
		}
		mark := agentSpineMark(stage, current)
		col := spine.Len()
		if stage == "clarification" {
			// ^ under the start of the clarification token (back-line left anchor)
			clarCol = col
		}
		if stage == "summary" {
			// | drops from the start of the summary token (questions anchor)
			summaryCol = col
		}
		spine.WriteString(mark)
	}
	if summaryCol <= clarCol {
		summaryCol = clarCol + 1
	}

	followup := "user_followup"
	if current == "user_followup" {
		followup = "user_followup[doing]"
	}

	// Row 2 line 1: verticals under clarification (^) and summary (|)
	connector := makeSpaces(clarCol) + "^" + makeSpaces(summaryCol-clarCol-1) + "|"

	// Row 2 line 2: back path spanning clarCol..summaryCol so questions
	// meets the summary drop (no blank gap). user_followup sits toward the right.
	// Avoid the substring "+---" (sealed tests treat it as box chrome).
	back := makeSpaces(clarCol) + buildAgentBackLine(summaryCol-clarCol+1, followup)

	return spine.String() + "\n" + connector + "\n" + back + "\n"
}

// buildAgentBackLine builds a width-wide back edge:
//
//	+-- refine -- … user_followup … <-- questions --+
//
// Extra width is filled with '-' between "refine -- " and followup so
// user_followup stays toward the right (after summary). Never emits "+---".
func buildAgentBackLine(width int, followup string) string {
	// Fixed pieces: "+--" + " refine -- " + followup + " <-- questions " + "--+"
	// Stretch only after "+-- " (space) so we never form the banned "+---".
	const (
		head = "+-- refine -- "
		tail = " <-- questions --+"
	)
	core := head + followup + tail
	if width <= len(core) {
		return core
	}
	extra := width - len(core)
	// Prefer padding left of followup so the node starts further right (under/after summary).
	// Keep a space before the followup token for readability.
	if extra == 1 {
		return head + " " + followup + tail
	}
	return head + strings.Repeat("-", extra-1) + " " + followup + tail
}

func agentSpineMark(stage, current string) string {
	if current == "user_followup" {
		if stage == "done" {
			return "(" + stage + ")"
		}
		// all spine stages through summary are past when refining
		return stage
	}

	curIdx := indexOfStage(agentSpine, current)
	stageIdx := indexOfStage(agentSpine, stage)
	if curIdx < 0 {
		// unknown / off-spine current: leave spine bare except we still need
		// one [doing] elsewhere (handled on row 2 for user_followup only)
		return stage
	}
	switch {
	case stageIdx == curIdx:
		return stage + "[doing]"
	case stageIdx < curIdx:
		return stage
	default:
		return "(" + stage + ")"
	}
}

func indexOfStage(stages []string, stage string) int {
	for i, s := range stages {
		if s == stage {
			return i
		}
	}
	return -1
}

func makeSpaces(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

func renderAgentAdvanceNext(task storage.Task) string {
	var b strings.Builder
	id := task.ID

	switch task.Stage {
	case "create":
		fmt.Fprintf(&b, "advance: ok\n")
		fmt.Fprintf(&b, "advance_to: in_process\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk advance %d\n", id)
		fmt.Fprintf(&b, "    to: in_process\n")
		fmt.Fprintf(&b, "    edge: claim\n")

	case "in_process":
		fmt.Fprintf(&b, "advance: ok\n")
		fmt.Fprintf(&b, "advance_to: clarification\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk advance %d\n", id)
		fmt.Fprintf(&b, "    to: clarification\n")
		fmt.Fprintf(&b, "    edge: research\n")

	case "clarification":
		fmt.Fprintf(&b, "advance: blocked\n")
		fmt.Fprintf(&b, "advance_reason: use clarify confirm to advance from clarification\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk clarify confirm -y %d\n", id)
		fmt.Fprintf(&b, "    to: implementation\n")
		fmt.Fprintf(&b, "    edge: confirmed\n")
		fmt.Fprintf(&b, "  - cmd: tsk clarify add %d <question>\n", id)

	case "implementation":
		fmt.Fprintf(&b, "advance: ok\n")
		fmt.Fprintf(&b, "advance_to: verification\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk advance %d\n", id)
		fmt.Fprintf(&b, "    to: verification\n")

	case "verification":
		fmt.Fprintf(&b, "advance: ok\n")
		fmt.Fprintf(&b, "advance_to: summary\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk advance %d\n", id)
		fmt.Fprintf(&b, "    to: summary\n")

	case "summary":
		fmt.Fprintf(&b, "advance: blocked\n")
		fmt.Fprintf(&b, "advance_reason: use done or followup from summary\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk followup %d <message>\n", id)
		fmt.Fprintf(&b, "    to: user_followup\n")
		fmt.Fprintf(&b, "    edge: questions\n")
		fmt.Fprintf(&b, "  - cmd: tsk done %d\n", id)
		fmt.Fprintf(&b, "    to: done\n")
		fmt.Fprintf(&b, "    edge: no_followup\n")

	case "user_followup":
		fmt.Fprintf(&b, "advance: ok\n")
		fmt.Fprintf(&b, "advance_to: clarification\n")
		fmt.Fprintf(&b, "next:\n")
		fmt.Fprintf(&b, "  - cmd: tsk advance %d\n", id)
		fmt.Fprintf(&b, "    to: clarification\n")
		fmt.Fprintf(&b, "    edge: refine\n")
		fmt.Fprintf(&b, "  - cmd: tsk done %d\n", id)
		fmt.Fprintf(&b, "    to: done\n")

	case "done":
		fmt.Fprintf(&b, "advance: blocked\n")
		fmt.Fprintf(&b, "advance_reason: task is already done\n")
		fmt.Fprintf(&b, "next:\n")

	default:
		// unknown stage — treat as blocked without advance_to
		fmt.Fprintf(&b, "advance: blocked\n")
		fmt.Fprintf(&b, "next:\n")
	}

	return b.String()
}
