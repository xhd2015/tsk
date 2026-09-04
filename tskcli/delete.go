package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runDelete(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("delete", args)

	var recursive, dryRun bool
	remaining, err := lessflags.
		Bool("--recursive", &recursive).
		Bool("--dry-run", &dryRun).
		Help("-h,--help", deleteHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk delete: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{TaskID: id})
	if dryRun {
		invk.setMutation(false)
	}

	plan, err := storage.PlanDelete(home, id, recursive)
	if err != nil {
		return fail(err)
	}
	if dryRun {
		for _, t := range plan.Tasks {
			fmt.Printf("[dry-run] would delete %d  [%d] %s\n", t.ID, t.ID, t.Title)
		}
		return nil
	}
	if err := storage.ApplyDelete(home, plan); err != nil {
		return fail(err)
	}
	fmt.Printf("deleted %d\n", id)
	return nil
}
