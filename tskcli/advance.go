package tskcli

import (
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runAdvance(home string, args []string) error {
	setCommand(currentCtx, "advance", args)

	var note string
	remaining, err := lessflags.String("--note", &note).Parse(args)
	if err != nil {
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk advance: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if err := storage.ValidateAdvance(task.Stage); err != nil {
		return fail(err)
	}
	to, _ := storage.CanAdvance(task.Stage)
	_, err = storage.RenameTaskDir(home, &task, taskDir, to, note)
	return err
}