package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runLogs(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("logs", args)

	var all, asJSON bool
	var limit int
	remaining, err := lessflags.
		Bool("--all", &all).
		Bool("--json", &asJSON).
		Int("--limit", &limit).
		Help("-h,--help", logsHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return logsErr("%s", err.Error())
	}
	if len(remaining) != 0 {
		return logsErr("tsk logs: unexpected args")
	}
	if limit < 0 {
		return logsErr("tsk logs: --limit must be >= 0")
	}

	events, err := storage.ReadEvents(home)
	if err != nil {
		return logsErr("%s", err.Error())
	}
	if !all {
		filtered := events[:0]
		for _, ev := range events {
			if eventIsMutation(ev) {
				filtered = append(filtered, ev)
			}
		}
		events = filtered
	}
	events = applyEventLimit(events, limit)

	if asJSON {
		if events == nil {
			events = []storage.Event{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(events)
	}

	for _, ev := range events {
		fmt.Println(formatLogLine(ev))
	}
	footer := logsCountLabel(len(events))
	if isStdoutTTY() {
		footer = ansiGray + footer + ansiReset
	}
	fmt.Println(footer)
	return nil
}

func logsErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "Error:") {
		msg = "Error: " + msg
	}
	return fail(fmt.Errorf("%s", msg))
}
