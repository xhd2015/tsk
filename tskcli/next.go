package tskcli

import (
	"fmt"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func runNext(home string, args []string) error {
	setCommand(currentCtx, "next", args)

	if len(args) != 0 {
		return fail(fmt.Errorf("tsk next: unexpected arguments"))
	}

	ids, err := storage.ListTaskIDs(home)
	if err != nil {
		return err
	}

	var bestID int
	var bestTime string
	found := false
	for _, id := range ids {
		task, _, err := storage.LoadTaskByID(home, id)
		if err != nil {
			return err
		}
		if task.Stage != "in_process" {
			continue
		}
		if !found {
			bestID = id
			bestTime = task.CreatedAt
			found = true
			continue
		}
		t1, err1 := parseCreatedAt(bestTime)
		t2, err2 := parseCreatedAt(task.CreatedAt)
		if err1 != nil || err2 != nil {
			if task.CreatedAt < bestTime || (task.CreatedAt == bestTime && id < bestID) {
				bestID = id
				bestTime = task.CreatedAt
			}
			continue
		}
		if t2.Before(t1) || (t2.Equal(t1) && id < bestID) {
			bestID = id
			bestTime = task.CreatedAt
		}
	}
	if found {
		fmt.Println(bestID)
	}
	return nil
}