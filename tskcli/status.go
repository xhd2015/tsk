package tskcli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/agent-pro/agent/agentrunner"
	"github.com/xhd2015/tsk/tskcli/pipeline"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runStatus(home string, args []string) error {
	setCommand(currentCtx, "status", args)

	// Use *string / presence-sensitive bools so we can distinguish
	// "flag absent" from "flag present" for auto format defaulting.
	var formatPtr *string
	var colorFlag bool
	var plain bool
	remaining, err := lessflags.
		String("--format", &formatPtr).
		Bool("--color", &colorFlag).
		Bool("--plain", &plain).
		Help("-h,--help", statusHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}

	format := resolveStatusFormat(formatPtr, colorFlag, plain)
	switch format {
	case "diagram", "agent":
		// ok
	default:
		return fail(fmt.Errorf("tsk status: invalid --format %q (allowed: diagram, agent)", format))
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk status: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	if format == "agent" {
		// agent view: plain facts + 2-row art; never ANSI even with --color
		fmt.Print(pipeline.RenderAgent(task, taskDir))
		return nil
	}

	color := colorFlag
	if !plain && !colorFlag {
		color = isStdoutTTY()
	}

	rendered := pipeline.Render(plain)
	out := pipeline.Highlight(rendered, task, color && !plain)
	fmt.Print(out)
	return nil
}

// resolveStatusFormat applies auto-format precedence (highest first):
//  1. --format present → use its value
//  2. --color or --plain present → diagram
//  3. TSK_STATUS_FORMAT=agent|diagram → that (invalid/empty ignored)
//  4. agentrunner.Detect ok → agent
//  5. else diagram
func resolveStatusFormat(formatPtr *string, colorFlag, plain bool) string {
	if formatPtr != nil {
		return *formatPtr
	}
	if colorFlag || plain {
		return "diagram"
	}
	if v := strings.TrimSpace(os.Getenv("TSK_STATUS_FORMAT")); v == "agent" || v == "diagram" {
		return v
	}
	if _, ok := agentrunner.Detect(agentrunner.Options{}); ok {
		return "agent"
	}
	return "diagram"
}

func isStdoutTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
