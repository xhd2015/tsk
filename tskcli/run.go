package tskcli

import (
	"fmt"

	"github.com/xhd2015/tsk/tskcli/storage"
)

// Run executes tsk logic with args.
func Run(args []string) error {
	home, err := storage.ResolveHome()
	if err != nil {
		return err
	}
	invk := &invocation{
		home: home,
		args: args,
	}
	var runErr error
	defer func() {
		exitCode := 0
		if runErr != nil {
			exitCode = 1
		}
		invk.finish(exitCode)
	}()
	runErr = dispatch(invk, args)
	return runErr
}

type invocation struct {
	home      string
	args      []string
	command   string
	action    string
	eventArgs []string
	mutation  bool
	data      *storage.EventData
}

func (invk *invocation) finish(exitCode int) {
	if invk.command == "" || invk.command == "logs" {
		return
	}
	ev := storage.Event{
		TS:       storage.NowLocalTimestamp(),
		Command:  invk.command,
		Action:   invk.action,
		Args:     invk.eventArgs,
		ExitCode: exitCode,
		User:     eventUser(),
		Mutation: invk.mutation,
		Data:     compactEventData(invk.data),
	}
	_ = storage.AppendEvent(invk.home, ev)
}

func dispatch(invk *invocation, args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(topHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runAdd(invk, args[1:])
	case "list":
		return runList(invk, args[1:])
	case "show":
		return runShow(invk, args[1:])
	case "status":
		return runStatus(invk, args[1:])
	case "advance":
		return runAdvance(invk, args[1:])
	case "stage":
		return runStage(invk, args[1:])
	case "next":
		return runNext(invk, args[1:])
	case "label":
		return runLabel(invk, args[1:])
	case "topic":
		return runTopic(invk, args[1:])
	case "clarify":
		return runClarify(invk, args[1:])
	case "followup":
		return runFollowup(invk, args[1:])
	case "done":
		return runDone(invk, args[1:])
	case "archive":
		return runArchive(invk, args[1:])
	case "delete":
		return runDelete(invk, args[1:])
	case "channel":
		return runChannel(invk, args[1:])
	case "note":
		return runNote(invk, args[1:])
	case "tree":
		return runTree(invk, args[1:])
	case "progress":
		return runProgress(invk, args[1:])
	case "search":
		return runSearch(invk, args[1:])
	case "logs":
		return runLogs(invk, args[1:])
	case "project":
		return runProject(invk, args[1:])
	case "update":
		return runUpdate(invk, args[1:])
	case "install":
		return runInstall(invk, args[1:])
	case "skill":
		return runSkill(invk, args[1:])
	default:
		return fmt.Errorf("tsk: unknown subcommand %q", args[0])
	}
}

func fail(err error) error {
	return err
}
