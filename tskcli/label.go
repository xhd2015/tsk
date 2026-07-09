package tskcli

import (
	"fmt"
	"sort"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func runLabel(home string, args []string) error {
	setCommand(currentCtx, "label", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(labelHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runLabelAdd(home, args[1:])
	case "rm":
		return runLabelRm(home, args[1:])
	default:
		return fail(fmt.Errorf("tsk label: unknown subcommand %q", args[0]))
	}
}

func runLabelAdd(home string, args []string) error {
	setCommand(currentCtx, "label", append([]string{"add"}, args...))

	if len(args) != 2 {
		return fail(fmt.Errorf("tsk label add: task id and label required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	label := args[1]

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if !containsLabel(task.Labels, label) {
		task.Labels = append(task.Labels, label)
		sort.Strings(task.Labels)
		task.UpdatedAt = storage.NowTimestamp(task.ID)
		return storage.WriteTask(taskDir, task)
	}
	return nil
}

func runLabelRm(home string, args []string) error {
	setCommand(currentCtx, "label", append([]string{"rm"}, args...))

	if len(args) != 2 {
		return fail(fmt.Errorf("tsk label rm: task id and label required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	label := args[1]

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	var kept []string
	for _, l := range task.Labels {
		if l != label {
			kept = append(kept, l)
		}
	}
	task.Labels = kept
	task.UpdatedAt = storage.NowTimestamp(task.ID)
	return storage.WriteTask(taskDir, task)
}