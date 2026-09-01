package tskcli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runTopic(home string, args []string) error {
	setCommand(currentCtx, "topic", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(topicHelp())
		return nil
	}
	switch args[0] {
	case "set":
		return runTopicSet(home, args[1:])
	case "mkdir":
		return runTopicMkdir(home, args[1:])
	case "rm":
		return runTopicRm(home, args[1:])
	case "where":
		return runTopicWhere(home, args[1:])
	case "info":
		return runTopicInfo(home, args[1:])
	case "note":
		return runTopicNote(home, args[1:])
	case "notes":
		return runTopicNotes(home, args[1:])
	case "view":
		return runTopicView(home, args[1:])
	case "alias":
		return runTopicAlias(home, args[1:])
	default:
		return fail(fmt.Errorf("tsk topic: unknown subcommand %q", args[0]))
	}
}

func topicErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "Error:") {
		msg = "Error: " + msg
	}
	return fail(fmt.Errorf("%s", msg))
}

func resolveTopicInput(home, ref string) ([]string, error) {
	parts, err := storage.LookupTopicRef(home, ref)
	if err != nil {
		return nil, err
	}
	return parts, nil
}

func requireTopicDir(home, ref string) ([]string, string, error) {
	parts, err := resolveTopicInput(home, ref)
	if err != nil {
		return nil, "", err
	}
	if !storage.TopicDirExists(home, parts) {
		return nil, "", fmt.Errorf("topic not found: %s", ref)
	}
	return parts, storage.TopicAbs(home, parts), nil
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
		topicParts, err = resolveTopicInput(home, args[1])
		if err != nil {
			return fail(err)
		}
	}

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if task.ParentID != 0 {
		return fail(fmt.Errorf("tsk topic set: nested task %d: reparent before changing topic", id))
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
	if len(parts) == 0 {
		return fail(fmt.Errorf("tsk topic mkdir: path required"))
	}
	resolved, err := resolveTopicInput(home, args[0])
	if err != nil {
		return fail(err)
	}
	if !topicPartsEqual(resolved, parts) {
		return topicErr("%s is an alias for %s", args[0], storage.JoinTopicPath(resolved))
	}
	dir := filepath.Join(home, "topics", filepath.Join(parts...))
	return os.MkdirAll(dir, 0o755)
}

func runTopicRm(home string, args []string) error {
	setCommand(currentCtx, "topic", append([]string{"rm"}, args...))

	remaining, err := lessflags.
		Help("-h,--help", topicRmHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk topic rm: path required"))
	}
	parts, abs, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	n, err := storage.CountTasksUnderTopic(home, parts)
	if err != nil {
		return fail(err)
	}
	if n > 0 {
		return topicErr("tsk topic rm: topic %s still has %d task(s)", storage.JoinTopicPath(parts), n)
	}
	subs, err := storage.ListTopicChildNames(abs)
	if err != nil {
		return fail(err)
	}
	if len(subs) > 0 {
		return topicErr("tsk topic rm: topic %s has subtopics %v", storage.JoinTopicPath(parts), subs)
	}
	if err := os.RemoveAll(abs); err != nil {
		return fail(fmt.Errorf("tsk topic rm: %w", err))
	}
	fmt.Printf("removed topic %s\n", storage.JoinTopicPath(parts))
	return nil
}

func topicPartsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func splitTopic(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}