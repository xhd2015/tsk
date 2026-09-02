package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runArchive(home string, args []string) error {
	setCommand(currentCtx, "archive", args)

	var force bool
	remaining, err := lessflags.
		Bool("--force", &force).
		Help("-h,--help", archiveHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	_ = force // accepted for compatibility; no extra effect
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk archive: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if storage.IsTerminal(task.Stage) {
		return fail(storage.TerminalTransitionError(task.Stage))
	}
	return storage.SetTaskStage(&task, taskDir, "archived", "")
}
