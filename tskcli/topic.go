package tskcli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/tsk/tskcli/storage"
)

func runTopic(home string, args []string) error {
	if len(args) == 0 {
		return fail(fmt.Errorf("tsk topic: subcommand required"))
	}
	switch args[0] {
	case "set":
		return runTopicSet(home, args[1:])
	case "mkdir":
		return runTopicMkdir(home, args[1:])
	default:
		return fail(fmt.Errorf("tsk topic: unknown subcommand %q", args[0]))
	}
}

func runTopicSet(home string, args []string) error {
	setCommand(currentCtx, "topic", append([]string{"set"}, args...))

	if len(args) < 2 {
		return fail(fmt.Errorf("tsk topic set: task id and path required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}

	var topicParts []string
	switch args[1] {
	case "--inbox", "":
		topicParts = nil
	default:
		topicParts = splitTopic(args[1])
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	_, err = storage.MoveTaskDir(home, &task, taskDir, topicParts)
	return err
}

func runTopicMkdir(home string, args []string) error {
	setCommand(currentCtx, "topic", append([]string{"mkdir"}, args...))

	if len(args) != 1 {
		return fail(fmt.Errorf("tsk topic mkdir: path required"))
	}
	parts := splitTopic(args[0])
	dir := filepath.Join(home, "topics", filepath.Join(parts...))
	return os.MkdirAll(dir, 0o755)
}

func splitTopic(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}