package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runDone(home string, args []string) error {
	setCommand(currentCtx, "done", args)

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
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk done: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if task.Stage == "done" {
		return fail(fmt.Errorf("invalid transition: task is already done"))
	}
	if task.Stage != "summary" && task.Stage != "user_followup" && !force {
		return fail(fmt.Errorf("invalid transition: done only from summary or user_followup"))
	}
	return storage.SetTaskStage(&task, taskDir, "done", "")
}
