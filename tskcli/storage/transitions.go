package storage

import "fmt"

// AllStages is the ordered workflow for status display.
// archived is terminal but off-spine (not listed here).
var AllStages = []string{
	"create",
	"in_process",
	"clarification",
	"implementation",
	"verification",
	"summary",
	"done",
}

// AdvanceNext maps a stage to the next stage for `tsk advance`.
var AdvanceNext = map[string]string{
	"create":         "in_process",
	"in_process":     "clarification",
	"implementation": "verification",
	"verification":   "summary",
	"user_followup":  "clarification",
}

// AllowedStageTargets lists valid direct targets for `tsk stage`.
// done remains reachable via stage from summary/user_followup only;
// archived is entered only via `tsk archive` (not listed as a target).
var AllowedStageTargets = map[string][]string{
	"create":         {"in_process"},
	"in_process":     {"clarification"},
	"clarification":  {"implementation"},
	"implementation": {"verification"},
	"verification":   {"summary"},
	"summary":        {"user_followup", "done"},
	"user_followup":  {"clarification", "done"},
	"done":           {},
	"archived":       {},
}

// IsTerminal reports whether stage is a finished terminal state.
func IsTerminal(stage string) bool {
	return stage == "done" || stage == "archived"
}

// TerminalTransitionError returns the error for mutating a terminal task.
func TerminalTransitionError(stage string) error {
	switch stage {
	case "archived":
		return fmt.Errorf("invalid transition: task is already archived")
	default:
		return fmt.Errorf("invalid transition: task is already done")
	}
}

// CanAdvance reports whether advance is allowed from the given stage.
func CanAdvance(from string) (string, bool) {
	to, ok := AdvanceNext[from]
	return to, ok
}

// CanStage reports whether a direct stage transition is allowed.
func CanStage(from, to string) bool {
	targets, ok := AllowedStageTargets[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}

// ValidateStageTransition returns an error for invalid transitions.
func ValidateStageTransition(from, to string) error {
	if IsTerminal(from) {
		return TerminalTransitionError(from)
	}
	if !CanStage(from, to) {
		return fmt.Errorf("invalid transition: %s -> %s", from, to)
	}
	return nil
}

// ValidateAdvance returns an error when advance is not allowed.
func ValidateAdvance(from string) error {
	if IsTerminal(from) {
		return TerminalTransitionError(from)
	}
	if from == "clarification" {
		return fmt.Errorf("invalid transition: use clarify confirm to advance from clarification")
	}
	if from == "summary" {
		return fmt.Errorf("invalid transition: use done or followup from summary")
	}
	if _, ok := CanAdvance(from); !ok {
		return fmt.Errorf("invalid transition: cannot advance from %s", from)
	}
	return nil
}
