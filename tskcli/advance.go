package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runAdvance(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("advance", args)

	var note string
	remaining, err := lessflags.
		String("--note", &note).
		Help("-h,--help", advanceHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk advance: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{TaskID: id, Text: note})

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if err := storage.ValidateAdvance(task.Stage); err != nil {
		return fail(err)
	}
	to, _ := storage.CanAdvance(task.Stage)
	invk.setData(storage.EventData{Stage: to})
	return storage.SetTaskStage(&task, taskDir, to, note)
}
