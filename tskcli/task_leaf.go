package tskcli

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/xhd2015/tsk/tskcli/storage"
)

// ANSI colors for task-stage styling (color-on tree leaves).
// ansiGreen / ansiGray / ansiReset live in channel.go (same package).
const (
	ansiCyan    = "\x1b[36m"
	ansiYellow  = "\x1b[33m"
	ansiBlue    = "\x1b[34m"
	ansiMagenta = "\x1b[35m"
)

// taskTitleMaxRunes caps the title portion of a human tree leaf.
const taskTitleMaxRunes = 512

// treeColorEnabled resolves whether tree-like human output should use ANSI.
// --color forces on; --plain / --json force off; auto uses TTY unless NO_COLOR
// is a non-empty env value (go-best-practice/cli/color).
func treeColorEnabled(colorFlag, plain, asJSON bool) bool {
	if asJSON || plain {
		return false
	}
	if colorFlag {
		return true
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	return isStdoutTTY()
}

// formatTaskStageExtra returns the trailing stage marker for a task leaf.
// Color on: empty (style conveys stage). Color off: "  (stage)".
func formatTaskStageExtra(stage string, color bool) string {
	if color || stage == "" {
		return ""
	}
	return "  (" + stage + ")"
}

// taskStageStyle returns an ANSI SGR prefix for the task name when color is on.
// create is unstyled; done is gray+strikethrough; other stages use a small palette.
func taskStageStyle(stage string, color bool) string {
	if !color {
		return ""
	}
	switch stage {
	case "create":
		return ""
	case "done":
		return ansiGray + ansiStrikethrough
	case "in_process":
		return ansiCyan
	case "clarification", "user_followup":
		return ansiYellow
	case "implementation":
		return ansiGreen
	case "verification":
		return ansiBlue
	case "summary":
		return ansiMagenta
	default:
		return ""
	}
}

func formatTaskIDBracket(id int) string {
	return fmt.Sprintf("[%d]", id)
}

// truncateRunes caps s to max runes; if truncated and max > 1, appends "…".
func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max == 1 {
		return string(runes[:1])
	}
	return string(runes[:max-1]) + "…"
}

// formatTaskLeafName builds an aligned human leaf: padded [id] + space + title.
// idWidth is the max visible width of [id] tokens in the rendered tree.
func formatTaskLeafName(id int, title, slug string, idWidth int) string {
	idTok := formatTaskIDBracket(id)
	w := utf8.RuneCountInString(idTok)
	pad := idWidth - w
	if pad < 0 {
		pad = 0
	}
	text := strings.TrimSpace(title)
	if text == "" {
		text = slug
	}
	text = truncateRunes(text, taskTitleMaxRunes)
	return idTok + strings.Repeat(" ", pad) + " " + text
}

// maxTaskIDWidth returns the max rune width of [id] among tasks (nested).
func maxTaskIDWidth(tasks []storage.TopicTreeTask) int {
	max := 0
	var walk func([]storage.TopicTreeTask)
	walk = func(ts []storage.TopicTreeTask) {
		for _, t := range ts {
			if w := utf8.RuneCountInString(formatTaskIDBracket(t.ID)); w > max {
				max = w
			}
			if len(t.Tasks) > 0 {
				walk(t.Tasks)
			}
		}
	}
	walk(tasks)
	return max
}

func maxTaskIDWidthForest(inbox []storage.TopicTreeTask, forest []storage.TopicTree) int {
	max := maxTaskIDWidth(inbox)
	var walkTopic func(*storage.TopicTree)
	walkTopic = func(tree *storage.TopicTree) {
		if w := maxTaskIDWidth(tree.Tasks); w > max {
			max = w
		}
		for i := range tree.Projects {
			if w := maxTaskIDWidth(tree.Projects[i].Tasks); w > max {
				max = w
			}
		}
		for i := range tree.Subtopics {
			walkTopic(&tree.Subtopics[i])
		}
	}
	for i := range forest {
		walkTopic(&forest[i])
	}
	// Also count project-bucketed inbox tasks already in inbox list (ungrouped
	// + we partition later). Inbox slice is the full inbox forest roots.
	return max
}
