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

func runNote(home string, args []string) error {
	setCommand(currentCtx, "note", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(noteHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runNoteAdd(home, args[1:])
	case "list":
		return runNoteList(home, args[1:])
	case "edit":
		return runNoteEdit(home, args[1:])
	default:
		return noteErr("tsk note: unknown subcommand %q", args[0])
	}
}

func noteErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "Error:") {
		msg = "Error: " + msg
	}
	return fail(fmt.Errorf("%s", msg))
}

func formatNoteLine(n storage.TopicNote) string {
	var parts []string
	if len(n.Labels) > 0 {
		parts = append(parts, fmt.Sprintf("[%s]", strings.Join(n.Labels, ", ")))
	}
	if n.Status != "" {
		parts = append(parts, fmt.Sprintf("(%s)", n.Status))
	}
	if len(parts) == 0 {
		return fmt.Sprintf("%s  %s", n.TS, n.Text)
	}
	return fmt.Sprintf("%s  %s  %s", n.TS, strings.Join(parts, "  "), n.Text)
}

func parseRequiredTaskID(raw string, cmd string) (int, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, noteErr("%s: --id required", cmd)
	}
	id, err := parseID(raw)
	if err != nil {
		return 0, noteErr("%s", err.Error())
	}
	return id, nil
}

func loadTaskDir(home string, id int) (string, error) {
	_, dir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return "", noteErr("task not found: %d", id)
	}
	return dir, nil
}

func runNoteAdd(home string, args []string) error {
	setCommand(currentCtx, "note", append([]string{"add"}, args...))

	var idStr string
	var labels []string
	remaining, err := lessflags.
		String("--id", &idStr).
		StringSlice("--label", &labels).
		Help("-h,--help", noteAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	id, err := parseRequiredTaskID(idStr, "tsk note add")
	if err != nil {
		return err
	}
	if len(remaining) == 0 {
		return noteErr("tsk note add: text required")
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
		TS:   storage.NowTimestamp(len(existing) + 1),
		Text: joinArgs(remaining),
	}
	if len(labels) > 0 {
		note.Labels = labels
	}
	if err := storage.AppendTopicNote(dir, note); err != nil {
		return fail(err)
	}
	fmt.Println("added note")
	return nil
}

func runNoteList(home string, args []string) error {
	setCommand(currentCtx, "note", append([]string{"list"}, args...))

	var idStr string
	var labels []string
	var asJSON bool
	var showIndex bool
	var limit int
	remaining, err := lessflags.
		String("--id", &idStr).
		StringSlice("--label", &labels).
		Bool("--json", &asJSON).
		Bool("--show-index", &showIndex).
		Int("--limit", &limit).
		Help("-h,--help", noteListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return noteErr("tsk note list: unexpected args")
	}
	if limit < 0 {
		return noteErr("tsk note list: --limit must be >= 0")
	}
	id, err := parseRequiredTaskID(idStr, "tsk note list")
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
	notes = storage.ApplyNoteLimit(storage.FilterNotes(notes, labels), limit)
	return printNotes(notes, asJSON, showIndex)
}

func printNotes(notes []storage.TopicNote, asJSON bool, showIndex bool) error {
	if asJSON {
		if notes == nil {
			notes = []storage.TopicNote{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(notes)
	}
	for i, n := range notes {
		if showIndex {
			fmt.Printf("%d.  %s\n", i+1, formatNoteLine(n))
		} else {
			fmt.Println(formatNoteLine(n))
		}
	}
	fmt.Printf("%d notes\n", len(notes))
	return nil
}

func runNoteEdit(home string, args []string) error {
	setCommand(currentCtx, "note", append([]string{"edit"}, args...))

	var idStr string
	var labels []string
	var status string
	var indexStr string
	var appendMode bool
	remaining, err := lessflags.
		String("--id", &idStr).
		StringSlice("--label", &labels).
		String("--status", &status).
		String("--index", &indexStr).
		Bool("--append", &appendMode).
		Help("-h,--help", noteEditHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	id, err := parseRequiredTaskID(idStr, "tsk note edit")
	if err != nil {
		return err
	}
	if indexStr == "" {
		return noteErr("tsk note edit: --index required")
	}
	index, err := parseID(indexStr)
	if err != nil || index <= 0 {
		return noteErr("tsk note edit: --index must be a positive integer")
	}
	if len(remaining) == 0 {
		return noteErr("tsk note edit: text required")
	}
	dir, err := loadTaskDir(home, id)
	if err != nil {
		return err
	}
	notes, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}

	// Map filtered index back to original position
	var filteredIdxs []int
	for i, n := range notes {
		if storage.NoteHasAllLabels(n, labels) {
			filteredIdxs = append(filteredIdxs, i)
		}
	}
	if len(filteredIdxs) == 0 {
		word := "notes"
		if len(labels) == 1 {
			word = "note"
		}
		return noteErr("tsk note edit: index %d out of range (have 0 %s)", index, word)
	}
	if index > len(filteredIdxs) {
		word := "notes"
		if len(filteredIdxs) == 1 {
			word = "note"
		}
		return noteErr("tsk note edit: index %d out of range (have %d %s)", index, len(filteredIdxs), word)
	}

	targetIdx := filteredIdxs[index-1]
	newText := joinArgs(remaining)
	if appendMode {
		notes[targetIdx].Text = notes[targetIdx].Text + newText
	} else {
		notes[targetIdx].Text = newText
	}
	if status != "" {
		notes[targetIdx].Status = status
	}

	if err := storage.RewriteNotes(dir, notes); err != nil {
		return fail(err)
	}
	fmt.Println("edited note")
	return nil
}
