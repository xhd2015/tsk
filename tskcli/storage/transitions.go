package storage

import "fmt"

// AllStages is the ordered workflow for status display.
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
var AllowedStageTargets = map[string][]string{
	"create":         {"in_process"},
	"in_process":     {"clarification"},
	"clarification":  {"implementation"},
	"implementation": {"verification"},
	"verification":   {"summary"},
	"summary":        {"user_followup", "done"},
	"user_followup":  {"clarification", "done"},
	"done":           {},
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
	if from == "done" {
		return fmt.Errorf("invalid transition: task is already done")
	}
	if !CanStage(from, to) {
		return fmt.Errorf("invalid transition: %s -> %s", from, to)
	}
	return nil
}

// ValidateAdvance returns an error when advance is not allowed.
func ValidateAdvance(from string) error {
	if from == "done" {
		return fmt.Errorf("invalid transition: task is already done")
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