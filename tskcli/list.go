package tskcli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runList(home string, args []string) error {
	setCommand(currentCtx, "list", args)

	var stage, label, topicPrefix string
	remaining, err := lessflags.
		String("--stage", &stage).
		String("--label", &label).
		String("--topic", &topicPrefix).
		Help("-h,--help", listHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return fail(fmt.Errorf("tsk list: unexpected arguments"))
	}

	ids, err := storage.ListTaskIDs(home)
	if err != nil {
		return err
	}
	for _, id := range ids {
		task, _, err := storage.LoadTaskByID(home, id)
		if err != nil {
			return err
		}
		if stage != "" && task.Stage != stage {
			continue
		}
		if label != "" && !containsLabel(task.Labels, label) {
			continue
		}
		if topicPrefix != "" {
			parts, err := storage.ParseTopicPath(task.TopicPath)
			if err != nil {
				return err
			}
			topicStr := strings.Join(parts, "/")
			if !strings.HasPrefix(topicStr, strings.Trim(topicPrefix, "/")) {
				continue
			}
		}
		fmt.Println(id)
	}
	return nil
}

func containsLabel(labels []string, want string) bool {
	for _, l := range labels {
		if l == want {
			return true
		}
	}
	return false
}

func parseCreatedAt(raw string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if ts, err := time.Parse(f, raw); err == nil {
			return ts, nil
		}
	}
	return time.Time{}, fmt.Errorf("parse created_at %q", raw)
}