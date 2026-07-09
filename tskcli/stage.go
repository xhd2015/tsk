package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runStage(home string, args []string) error {
	setCommand(currentCtx, "stage", args)

	var note string
	remaining, err := lessflags.
		String("--note", &note).
		Help("-h,--help", stageHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 2 {
		return fail(fmt.Errorf("tsk stage: usage: tsk stage <id> <stage>"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	target := remaining[1]

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if err := storage.ValidateStageTransition(task.Stage, target); err != nil {
		return fail(err)
	}
	_, err = storage.RenameTaskDir(home, &task, taskDir, target, note)
	return err
}