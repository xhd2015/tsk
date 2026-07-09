package tskcli

import (
	"fmt"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func runStatus(home string, args []string) error {
	setCommand(currentCtx, "status", args)

	if len(args) != 1 {
		return fail(fmt.Errorf("tsk status: task id required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}

	task, _, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	stages := storage.AllStages
	for i, stage := range stages {
		line := stage
		if stage == task.Stage {
			line = "> " + stage + " <"
		} else {
			line = "  " + stage
		}
		fmt.Println(line)
		if i < len(stages)-1 {
			fmt.Println("  |")
		}
	}
	return nil
}