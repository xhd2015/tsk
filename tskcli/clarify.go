package tskcli

import (
	"errors"
	"fmt"
	"os"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runClarify(home string, args []string) error {
	setCommand(currentCtx, "clarify", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(clarifyHelp())
		return nil
	}
	sub := args[0]
	subArgs := args[1:]
	switch sub {
	case "add":
		return runClarifyAdd(home, subArgs)
	case "list":
		return runClarifyList(home, subArgs)
	case "confirm":
		return runClarifyConfirm(home, subArgs)
	default:
		return fail(fmt.Errorf("tsk clarify: unknown subcommand %q", sub))
	}
}

func runClarifyAdd(home string, args []string) error {
	setCommand(currentCtx, "clarify", append([]string{"add"}, args...))

	if len(args) < 2 {
		return fail(fmt.Errorf("tsk clarify add: task id and question required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	question := joinArgs(args[1:])

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if task.Stage != "clarification" {
		return fail(fmt.Errorf("tsk clarify add: task not in clarification stage"))
	}

	batch, err := storage.EnsureClarifyBatch(taskDir)
	if err != nil {
		return err
	}
	itemID := fmt.Sprintf("q%d", len(batch.Items)+1)
	batch.Items = append(batch.Items, storage.ClarifyItem{
		ID:       itemID,
		Question: question,
		Status:   "pending",
	})
	return storage.WriteClarifyBatch(taskDir, batch)
}

func runClarifyList(home string, args []string) error {
	setCommand(currentCtx, "clarify", append([]string{"list"}, args...))

	if len(args) != 1 {
		return fail(fmt.Errorf("tsk clarify list: task id required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	_, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	batch, err := storage.ReadClarifyBatch(taskDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, item := range batch.Items {
		fmt.Printf("%s: %s (%s)\n", item.ID, item.Question, item.Status)
	}
	return nil
}

func runClarifyConfirm(home string, args []string) error {
	setCommand(currentCtx, "clarify", append([]string{"confirm"}, args...))

	var assumeYes bool
	remaining, err := lessflags.
		Bool("-y,--yes", &assumeYes).
		Help("-h,--help", clarifyHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk clarify confirm: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if task.Stage != "clarification" {
		return fail(fmt.Errorf("tsk clarify confirm: task not in clarification stage"))
	}

	batch, err := storage.ReadClarifyBatch(taskDir)
	if err != nil {
		return fail(fmt.Errorf("tsk clarify confirm: no clarify batch"))
	}
	if len(batch.Items) == 0 {
		return fail(fmt.Errorf("tsk clarify confirm: no questions to confirm"))
	}
	if !assumeYes {
		return fail(fmt.Errorf("tsk clarify confirm: confirmation required (-y)"))
	}
	for i := range batch.Items {
		batch.Items[i].Status = "confirmed"
	}
	batch.Status = "confirmed"
	if err := storage.WriteClarifyBatch(taskDir, batch); err != nil {
		return err
	}

	_, err = storage.RenameTaskDir(home, &task, taskDir, "implementation", "")
	return err
}