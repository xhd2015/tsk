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

var progressStatuses = map[string]struct{}{
	"in-progress": {},
	"blocked":     {},
	"done":        {},
	"archived":    {},
}

const ansiStrikethrough = "\x1b[9m"

func runProgress(home string, args []string) error {
	setCommand(currentCtx, "progress", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(progressHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runProgressAdd(home, args[1:])
	case "list":
		return runProgressList(home, args[1:])
	case "edit":
		return runProgressEdit(home, args[1:])
	case "archive":
		return runProgressArchive(home, args[1:])
	case "show":
		return runProgressShow(home, args[1:])
	default:
		return progressErr("tsk progress: unknown subcommand %q", args[0])
	}
}

func progressErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "Error:") {
		msg = "Error: " + msg
	}
	return fail(fmt.Errorf("%s", msg))
}

func validateProgressStatus(status string, required bool) error {
	if status == "" {
		if required {
			return progressErr("--status required (allowed: in-progress, blocked, done, archived)")
		}
		return nil
	}
	if _, ok := progressStatuses[status]; !ok {
		return progressErr("invalid --status %q (allowed: in-progress, blocked, done, archived)", status)
	}
	return nil
}

func runProgressAdd(home string, args []string) error {
	setCommand(currentCtx, "progress", append([]string{"add"}, args...))

	var idStr, status string
	remaining, err := lessflags.
		String("--id", &idStr).
		String("--status", &status).
		Help("-h,--help", progressAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	id, err := parseRequiredTaskID(idStr, "tsk progress add")
	if err != nil {
		return err
	}
	if err := validateProgressStatus(status, true); err != nil {
		return err
	}
	if len(remaining) == 0 {
		return progressErr("tsk progress add: text required")
	}
	dir, err := loadTaskDir(home, id)
	if err != nil {
		return err
	}
	existing, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	note := storage.TopicNote{
		TS:     storage.NowTimestamp(len(existing) + 1),
		Text:   joinArgs(remaining),
		Labels: []string{"progress"},
		Status: status,
	}
	if err := storage.AppendTopicNote(dir, note); err != nil {
		return fail(err)
	}
	fmt.Println("added progress")
	return nil
}

func runProgressList(home string, args []string) error {
	setCommand(currentCtx, "progress", append([]string{"list"}, args...))

	var idStr, status string
	var asJSON, showIndex bool
	var limit int
	remaining, err := lessflags.
		String("--id", &idStr).
		String("--status", &status).
		Bool("--json", &asJSON).
		Bool("--show-index", &showIndex).
		Int("--limit", &limit).
		Help("-h,--help", progressListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return progressErr("tsk progress list: unexpected args")
	}
	if limit < 0 {
		return progressErr("tsk progress list: --limit must be >= 0")
	}
	id, err := parseRequiredTaskID(idStr, "tsk progress list")
	if err != nil {
		return err
	}
	dir, err := loadTaskDir(home, id)
	if err != nil {
		return err
	}
	notes, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	notes = storage.FilterNotes(notes, []string{"progress"})
	if status != "" {
		notes = filterByStatus(notes, status)
	}
	notes = storage.ApplyNoteLimit(notes, limit)
	return printProgressEntries(notes, asJSON, showIndex)
}

func runProgressEdit(home string, args []string) error {
	setCommand(currentCtx, "progress", append([]string{"edit"}, args...))

	var idStr, indexStr, status string
	remaining, err := lessflags.
		String("--id", &idStr).
		String("--index", &indexStr).
		String("--status", &status).
		Help("-h,--help", progressEditHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	id, err := parseRequiredTaskID(idStr, "tsk progress edit")
	if err != nil {
		return err
	}
	index, err := parseRequiredProgressIndex(indexStr, "tsk progress edit")
	if err != nil {
		return err
	}
	if err := validateProgressStatus(status, true); err != nil {
		return err
	}
	return updateProgress(home, id, index, status, strings.TrimSpace(joinArgs(remaining)), "updated progress")
}

func runProgressArchive(home string, args []string) error {
	setCommand(currentCtx, "progress", append([]string{"archive"}, args...))

	var idStr, indexStr string
	remaining, err := lessflags.
		String("--id", &idStr).
		String("--index", &indexStr).
		Help("-h,--help", progressArchiveHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return progressErr("tsk progress archive: unexpected arguments")
	}
	id, err := parseRequiredTaskID(idStr, "tsk progress archive")
	if err != nil {
		return err
	}
	index, err := parseRequiredProgressIndex(indexStr, "tsk progress archive")
	if err != nil {
		return err
	}
	return updateProgress(home, id, index, "archived", "", "archived progress")
}

func parseRequiredProgressIndex(raw, cmd string) (int, error) {
	if raw == "" {
		return 0, progressErr("%s: --index required", cmd)
	}
	index, err := parseID(raw)
	if err != nil {
		return 0, progressErr("%s: --index must be a positive integer", cmd)
	}
	return index, nil
}

func updateProgress(home string, id, index int, status, text, success string) error {
	dir, err := loadTaskDir(home, id)
	if err != nil {
		return err
	}
	notes, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	var positions []int
	for i, n := range notes {
		if storage.NoteHasAllLabels(n, []string{"progress"}) {
			positions = append(positions, i)
		}
	}
	if index > len(positions) {
		word := "entries"
		if len(positions) == 1 {
			word = "entry"
		}
		return progressErr("index %d out of range (have %d %s)", index, len(positions), word)
	}
	target := positions[index-1]
	notes[target].Status = status
	if text != "" {
		notes[target].Text = text
	}
	if err := storage.RewriteNotes(dir, notes); err != nil {
		return fail(err)
	}
	fmt.Println(success)
	return nil
}

func runProgressShow(home string, args []string) error {
	setCommand(currentCtx, "progress", append([]string{"show"}, args...))

	var idStr string
	remaining, err := lessflags.
		String("--id", &idStr).
		Help("-h,--help", progressShowHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return progressErr("tsk progress show: unexpected args")
	}
	id, err := parseRequiredTaskID(idStr, "tsk progress show")
	if err != nil {
		return err
	}
	dir, err := loadTaskDir(home, id)
	if err != nil {
		return err
	}
	notes, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	notes = storage.FilterNotes(notes, []string{"progress"})
	if len(notes) == 0 {
		fmt.Println("no progress")
		return nil
	}
	latest := notes[len(notes)-1]
	fmt.Println(formatNoteLine(latest))
	return nil
}

func filterByStatus(notes []storage.TopicNote, status string) []storage.TopicNote {
	var out []storage.TopicNote
	for _, n := range notes {
		if n.Status == status {
			out = append(out, n)
		}
	}
	return out
}

func isTerminalProgressStatus(status string) bool {
	return status == "done" || status == "archived"
}

func formatProgressLine(n storage.TopicNote, color bool) string {
	line := formatNoteLine(n)
	if color && isTerminalProgressStatus(n.Status) {
		return ansiGray + ansiStrikethrough + line + ansiReset
	}
	return line
}

func printProgressEntries(notes []storage.TopicNote, asJSON, showIndex bool) error {
	if asJSON {
		if notes == nil {
			notes = []storage.TopicNote{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(notes)
	}
	for i, n := range notes {
		line := formatProgressLine(n, isStdoutTTY())
		if showIndex {
			fmt.Printf("%d.  %s\n", i+1, line)
		} else {
			fmt.Println(line)
		}
	}
	word := "entries"
	if len(notes) == 1 {
		word = "entry"
	}
	fmt.Printf("%d %s\n", len(notes), word)
	return nil
}
