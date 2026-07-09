package tskcli

import (
	"errors"
	"fmt"
	"os"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/pipeline"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runStatus(home string, args []string) error {
	setCommand(currentCtx, "status", args)

	var colorFlag bool
	var plain bool
	var format string
	remaining, err := lessflags.
		String("--format", &format).
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
	if format == "" {
		format = "diagram"
	}
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

	task, _, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	if format == "agent" {
		// agent view: plain facts + 2-row art; never ANSI even with --color
		fmt.Print(pipeline.RenderAgent(task))
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

func isStdoutTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}