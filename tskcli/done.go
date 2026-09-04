package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runDone(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("done", args)

	var force bool
	remaining, err := lessflags.
		Bool("--force", &force).
		Help("-h,--help", doneHelp()).
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
		return fail(fmt.Errorf("tsk done: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{TaskID: id})

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if storage.IsTerminal(task.Stage) {
		return fail(storage.TerminalTransitionError(task.Stage))
	}
	invk.setData(storage.EventData{Stage: "done"})
	return storage.SetTaskStage(&task, taskDir, "done", "")
}
