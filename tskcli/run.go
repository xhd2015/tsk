package tskcli

import (
	"fmt"
	"os"
	"time"

	"github.com/xhd2015/tsk/tskcli/storage"
)

// Run executes tsk logic with args.
func Run(args []string) error {
	home, err := storage.ResolveHome()
	if err != nil {
		return err
	}
	ctx := &invocationContext{
		home:  home,
		args:  args,
		start: time.Now().UTC(),
	}
	initRunCtx(ctx)
	var runErr error
	defer func() {
		exitCode := 0
		if runErr != nil {
			exitCode = 1
		}
		ctx.finish(exitCode)
	}()
	runErr = dispatch(home, args)
	return runErr
}

type invocationContext struct {
	home    string
	args    []string
	command string
	eventArgs []string
	start   time.Time
}

func (ctx *invocationContext) finish(exitCode int) {
	ev := storage.Event{
		TS:       ctx.start.Format(time.RFC3339),
		Command:  ctx.command,
		Args:     ctx.eventArgs,
		ExitCode: exitCode,
	}
	_ = storage.AppendEvent(ctx.home, ev)
}

func dispatch(home string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("tsk: missing subcommand")
	}
	switch args[0] {
	case "create":
		return runCreate(home, args[1:])
	case "list":
		return runList(home, args[1:])
	case "show":
		return runShow(home, args[1:])
	case "status":
		return runStatus(home, args[1:])
	case "advance":
		return runAdvance(home, args[1:])
	case "stage":
		return runStage(home, args[1:])
	case "next":
		return runNext(home, args[1:])
	case "label":
		return runLabel(home, args[1:])
	case "topic":
		return runTopic(home, args[1:])
	case "clarify":
		return runClarify(home, args[1:])
	case "followup":
		return runFollowup(home, args[1:])
	case "done":
		return runDone(home, args[1:])
	default:
		return fmt.Errorf("tsk: unknown subcommand %q", args[0])
	}
}

func setCommand(ctx *invocationContext, command string, eventArgs []string) {
	ctx.command = command
	if eventArgs == nil {
		eventArgs = []string{}
	}
	ctx.eventArgs = eventArgs
}

func fail(err error) error {
	fmt.Fprintln(os.Stderr, err)
	return err
}