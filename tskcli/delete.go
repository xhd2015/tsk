package tskcli

import (
	"errors"
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runDelete(home string, args []string) error {
	setCommand(currentCtx, "delete", args)

	var recursive bool
	remaining, err := lessflags.
		Bool("--recursive", &recursive).
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

	if err := storage.DeleteTask(home, id, recursive); err != nil {
		return fail(err)
	}
	fmt.Printf("deleted %d\n", id)
	return nil
}
