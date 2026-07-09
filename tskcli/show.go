package tskcli

import (
	"errors"
	"fmt"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runShow(home string, args []string) error {
	setCommand(currentCtx, "show", args)

	remaining, err := lessflags.
		Help("-h,--help", showHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk show: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}

	task, _, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	topicParts, err := storage.ParseTopicPath(task.TopicPath)
	if err != nil {
		return err
	}
	var topicStr string
	if len(topicParts) == 0 {
		topicStr = "inbox"
	} else {
		topicStr = strings.Join(topicParts, "/")
	}

	fmt.Printf("id: %d\n", task.ID)
	fmt.Printf("title: %s\n", task.Title)
	fmt.Printf("slug: %s\n", task.Slug)
	fmt.Printf("stage: %s\n", task.Stage)
	fmt.Printf("topic: %s\n", topicStr)
	if len(task.Labels) == 0 {
		fmt.Println("labels:")
	} else {
		fmt.Printf("labels: %s\n", strings.Join(task.Labels, ", "))
	}
	fmt.Printf("created_at: %s\n", task.CreatedAt)
	fmt.Printf("updated_at: %s\n", task.UpdatedAt)
	return nil
}