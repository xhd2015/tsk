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

func runTopic(invk *invocation, args []string) error {
	invk.setCommand("topic", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(topicHelp())
		return nil
	}
	switch args[0] {
	case "set":
		return runTopicSet(invk, args[1:])
	case "mkdir":
		return runTopicMkdir(invk, args[1:])
	case "rm":
		return runTopicRm(invk, args[1:])
	case "where":
		return runTopicWhere(invk, args[1:])
	case "info":
		return runTopicInfo(invk, args[1:])
	case "note":
		return runTopicNote(invk, args[1:])
	case "notes":
		return runTopicNotes(invk, args[1:])
	case "view":
		return runTopicView(invk, args[1:])
	case "alias":
		return runTopicAlias(invk, args[1:])
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

func runTopicSet(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"set"}, args...))

	if len(args) < 2 {
		return fail(fmt.Errorf("tsk topic set: task id and path required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{TaskID: id})

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
	if len(topicParts) > 0 {
		invk.setData(storage.EventData{Topic: storage.JoinTopicPath(topicParts)})
	}
	_, err = storage.MoveTaskDir(home, &task, taskDir, topicParts)
	return err
}

func runTopicMkdir(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"mkdir"}, args...))

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
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts)})
	return os.MkdirAll(dir, 0o755)
}

func runTopicRm(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"rm"}, args...))

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
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts)})
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
