package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/xhd2015/tsk/script/check-channel-activity/signals"
	"github.com/xhd2015/tsk/script/check-channel-activity/signals"
	lessflags "github.com/xhd2015/less-flags"
)

const helpText = `Usage: check-channel-activity --channel-id ID --exec-if-idle-1h LINE [options]

Monitor a tsk channel's last message activity and run a notify command when idle.

Options:
  --channel-id ID           channel to monitor (required)
  --exec-if-idle-1h LINE    shell command line to run when idle (required; quote if spaces)
  --idle DURATION           idle threshold (default: 1h)
  --forever                 loop until SIGINT/SIGTERM
  --interval DURATION       sleep between checks with --forever (default: 1m)
  --tsk-home PATH           tsk storage root (default: $TSK_HOME or ~/.tsk)
  --state-dir PATH          notification state directory (default: $TSK_HOME/channels/state)
  --dry-run                 print status without executing notify command
  -h, --help                show help
`

func main() {
	if hasFlag(os.Args[1:], "--forever") {
		signals.Refresh()
		if signals.WaitDuringStartup() {
			fmt.Fprint(os.Stderr, "stopped\n")
			return
		}
	}
	signals.Refresh()
	if signals.Stopped() {
		fmt.Fprint(os.Stderr, "stopped\n")
		return
	}
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", strings.TrimPrefix(err.Error(), "Error: "))
		os.Exit(1)
	}
}

func run(args []string) error {
	var (
		channelID   string
		execLine    string
		idle        = time.Hour
		forever     bool
		interval    = time.Minute
		maxTicks    int
		tskHome     string
		stateDir    string
		dryRun      bool
		hasExecFlag bool
	)

	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--exec-if-idle-1h", &execLine).
		Duration("--idle", &idle).
		Bool("--forever", &forever).
		Duration("--interval", &interval).
		Int("--max-ticks", &maxTicks).
		String("--tsk-home", &tskHome).
		String("--state-dir", &stateDir).
		Bool("--dry-run", &dryRun).
		Help("-h,--help", helpText).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}
	if len(remaining) != 0 {
		return fmt.Errorf("unexpected arguments: %v", remaining)
	}

	for _, a := range args {
		if a == "--exec-if-idle-1h" {
			hasExecFlag = true
			break
		}
	}

	if channelID == "" {
		return fmt.Errorf("--channel-id is required")
	}
	if !hasExecFlag {
		return fmt.Errorf("--exec-if-idle-1h is required")
	}
	if execLine == "" {
		return fmt.Errorf("--exec-if-idle-1h requires a command")
	}

	home, err := resolveHome(tskHome)
	if err != nil {
		return err
	}
	if stateDir == "" {
		stateDir = filepath.Join(home, "channels", "state")
	}

	cfg := &runConfig{
		home:      home,
		stateDir:  stateDir,
		channelID: channelID,
		execLine:  execLine,
		idle:      idle,
		forever:   forever,
		interval:  interval,
		maxTicks:  maxTicks,
		dryRun:    dryRun,
	}
	return runCheck(cfg)
}

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}