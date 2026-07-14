package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
	"github.com/xhd2015/tsk/script/check-channel-activity/signals"
)

type runConfig struct {
	home      string
	stateDir  string
	channelID string
	execLine  string
	idle      time.Duration
	forever   bool
	interval  time.Duration
	maxTicks  int
	dryRun    bool
}

func runCheck(cfg *runConfig) error {
	if cfg.forever {
		return runForever(cfg)
	}
	return runOnce(cfg)
}

func stopForever() error {
	fmt.Fprint(os.Stderr, "stopped\n")
	return nil
}

func runForever(cfg *runConfig) error {
	signals.Refresh()
	ticks := 0
	for {
		if signals.Stopped() {
			return stopForever()
		}
		if err := runOnce(cfg); err != nil {
			return err
		}
		ticks++
		if cfg.maxTicks > 0 && ticks >= cfg.maxTicks {
			return nil
		}
		select {
		case <-signals.SigCh:
			return stopForever()
		case <-time.After(cfg.interval):
		}
	}
}

func runOnce(cfg *runConfig) error {
	activity, err := loadChannelActivity(cfg.home, cfg.channelID)
	if err != nil {
		return err
	}
	activity.checkIdle(cfg.idle)

	st, err := loadState(cfg.stateDir, cfg.channelID)
	if err != nil {
		return err
	}

	status, err := decideAndAct(cfg, activity, st)
	if err != nil {
		return err
	}
	printStatus(activity, status)
	return nil
}

func parseExecArgv(execLine string) ([]string, error) {
	argv, err := shellwords.Parse(execLine)
	if err != nil {
		return nil, fmt.Errorf("parse exec command: %w", err)
	}
	if len(argv) == 0 {
		return nil, fmt.Errorf("no command specified for --exec-if-idle-1h")
	}
	if len(argv) == 4 && argv[1] == "-c" && isShell(argv[0]) {
		argv = []string{argv[0], argv[1], argv[2], "sh", argv[3]}
	}
	return argv, nil
}

func isShell(path string) bool {
	base := strings.TrimSuffix(filepath.Base(path), ".exe")
	return base == "sh" || base == "bash" || base == "dash" || base == "zsh" || base == "ksh"
}

func decideAndAct(cfg *runConfig, activity *activityResult, st *notifyState) (string, error) {
	if !activity.IsIdle {
		return "active", nil
	}

	if cfg.dryRun {
		if alreadyNotified(st, activity.LastActivity) {
			return "already notified", nil
		}
		return "would notify (dry-run)", nil
	}

	if alreadyNotified(st, activity.LastActivity) {
		return "already notified", nil
	}

	if cfg.execLine == "" {
		return "", fmt.Errorf("no command specified for --exec-if-idle-1h")
	}
	argv, err := parseExecArgv(cfg.execLine)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec notify command: %w", err)
	}

	if err := writeState(cfg.stateDir, cfg.channelID, activity.LastActivity); err != nil {
		return "", err
	}
	return "notified", nil
}