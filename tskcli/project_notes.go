package tskcli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runProjectNotes(invk *invocation, args []string) error {
	invk.setCommand("project", append([]string{"notes"}, args...))

	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help":
			fmt.Print(projectNotesHelp())
			return nil
		case "add":
			return runProjectNotesAdd(invk, args[1:])
		case "list":
			return runProjectNotesList(invk, args[1:])
		case "edit":
			return runProjectNotesEdit(invk, args[1:])
		case "delete":
			return runProjectNotesDelete(invk, args[1:])
		}
	}
	// Bare `tsk project notes` lists.
	return runProjectNotesList(invk, args)
}

func resolveNotesProject(home, dirFlag, nameFlag, projectFlag string) (cwd string, ref storage.ProjectRef, err error) {
	nameFlag = strings.TrimSpace(nameFlag)
	projectFlag = strings.TrimSpace(projectFlag)
	if nameFlag != "" && projectFlag != "" {
		return "", storage.ProjectRef{}, fmt.Errorf("tsk project notes: --name and --project are mutually exclusive")
	}
	if strings.TrimSpace(dirFlag) != "" && (nameFlag != "" || projectFlag != "") {
		return "", storage.ProjectRef{}, fmt.Errorf("tsk project notes: --dir conflicts with --name/--project")
	}
	projectName := projectFlag
	if projectName == "" {
		projectName = nameFlag
	}
	cwd, ref, err = resolveAddProject(home, dirFlag, projectName)
	if err != nil {
		// Rewrite add-specific prefix for notes.
		msg := err.Error()
		msg = strings.Replace(msg, "tsk project add:", "tsk project notes:", 1)
		return "", storage.ProjectRef{}, fmt.Errorf("%s", msg)
	}
	return cwd, ref, nil
}

func ensureNotesDir(home, dirFlag, nameFlag, projectFlag string) (notesDir string, err error) {
	cwd, ref, err := resolveNotesProject(home, dirFlag, nameFlag, projectFlag)
	if err != nil {
		return "", err
	}
	loc := pathfmt.TildeHome(mainRepoDirForAuto(cwd))
	id, err := storage.EnsureProjectID(home, ref, loc)
	if err != nil {
		return "", err
	}
	return storage.ProjectNotesDir(home, id), nil
}

func runProjectNotesAdd(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("project", append([]string{"notes", "add"}, args...))

	var dirFlag, nameFlag, projectFlag string
	var labels []string
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		String("--name", &nameFlag).
		String("--project", &projectFlag).
		StringSlice("--label", &labels).
		Help("-h,--help", projectNotesAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return projectFail(fmt.Errorf("tsk project notes add: %v", err))
		}
	}
	if len(remaining) == 0 {
		return projectFail(fmt.Errorf("tsk project notes add: text required"))
	}

	notesDir, err := ensureNotesDir(home, dirFlag, nameFlag, projectFlag)
	if err != nil {
		return projectFail(err)
	}
	existing, err := storage.ReadTopicNotes(notesDir)
	if err != nil {
		return projectFail(err)
	}
	text := joinArgs(remaining)
	note := storage.TopicNote{
		TS:   storage.NowTimestamp(len(existing) + 1),
		Text: text,
	}
	if len(labels) > 0 {
		note.Labels = labels
	}
	if err := storage.AppendTopicNote(notesDir, note); err != nil {
		return projectFail(err)
	}
	invk.setData(storage.EventData{
		Name:    firstNonEmpty(nameFlag, projectFlag),
		Project: firstNonEmpty(projectFlag, nameFlag),
		Text:    text,
		Labels:  labels,
	})
	fmt.Println("added note")
	return nil
}

func runProjectNotesList(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("project", append([]string{"notes", "list"}, args...))

	var dirFlag, nameFlag, projectFlag string
	var labels []string
	var asJSON, showIndex bool
	var limit int
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		String("--name", &nameFlag).
		String("--project", &projectFlag).
		StringSlice("--label", &labels).
		Bool("--json", &asJSON).
		Bool("--show-index", &showIndex).
		Int("--limit", &limit).
		Help("-h,--help", projectNotesListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project notes list: unexpected args"))
	}
	if limit < 0 {
		return projectFail(fmt.Errorf("tsk project notes list: --limit must be >= 0"))
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return projectFail(fmt.Errorf("tsk project notes list: %v", err))
		}
	}

	notesDir, err := ensureNotesDir(home, dirFlag, nameFlag, projectFlag)
	if err != nil {
		return projectFail(err)
	}
	invk.setData(storage.EventData{
		Name:    firstNonEmpty(nameFlag, projectFlag),
		Project: firstNonEmpty(projectFlag, nameFlag),
	})
	notes, err := storage.ReadTopicNotes(notesDir)
	if err != nil {
		return projectFail(err)
	}
	notes = storage.ApplyNoteLimit(storage.FilterNotes(notes, labels), limit)
	return printNotes(notes, asJSON, showIndex)
}

func runProjectNotesEdit(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("project", append([]string{"notes", "edit"}, args...))

	var dirFlag, nameFlag, projectFlag, indexStr string
	var labels []string
	var appendMode bool
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		String("--name", &nameFlag).
		String("--project", &projectFlag).
		StringSlice("--label", &labels).
		String("--index", &indexStr).
		Bool("--append", &appendMode).
		Help("-h,--help", projectNotesEditHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return projectFail(fmt.Errorf("tsk project notes edit: %v", err))
		}
	}
	if indexStr == "" {
		return projectFail(fmt.Errorf("tsk project notes edit: --index required"))
	}
	index, err := parseID(indexStr)
	if err != nil || index <= 0 {
		return projectFail(fmt.Errorf("tsk project notes edit: --index must be a positive integer"))
	}
	if len(remaining) == 0 {
		return projectFail(fmt.Errorf("tsk project notes edit: text required"))
	}

	notesDir, err := ensureNotesDir(home, dirFlag, nameFlag, projectFlag)
	if err != nil {
		return projectFail(err)
	}
	notes, err := storage.ReadTopicNotes(notesDir)
	if err != nil {
		return projectFail(err)
	}

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
		return projectFail(fmt.Errorf("tsk project notes edit: index %d out of range (have 0 %s)", index, word))
	}
	if index > len(filteredIdxs) {
		word := "notes"
		if len(filteredIdxs) == 1 {
			word = "note"
		}
		return projectFail(fmt.Errorf("tsk project notes edit: index %d out of range (have %d %s)", index, len(filteredIdxs), word))
	}

	targetIdx := filteredIdxs[index-1]
	newText := joinArgs(remaining)
	invk.setData(storage.EventData{
		Name:    firstNonEmpty(nameFlag, projectFlag),
		Project: firstNonEmpty(projectFlag, nameFlag),
		Index:   index,
		Text:    newText,
		Labels:  labels,
	})
	if appendMode {
		notes[targetIdx].Text = notes[targetIdx].Text + newText
	} else {
		notes[targetIdx].Text = newText
	}
	if err := storage.RewriteNotes(notesDir, notes); err != nil {
		return projectFail(err)
	}
	fmt.Println("edited note")
	return nil
}

func runProjectNotesDelete(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("project", append([]string{"notes", "delete"}, args...))

	var dirFlag, nameFlag, projectFlag, indexStr string
	var labels []string
	remaining, err := lessflags.
		String("--dir", &dirFlag).
		String("--name", &nameFlag).
		String("--project", &projectFlag).
		StringSlice("--label", &labels).
		String("--index", &indexStr).
		Help("-h,--help", projectNotesDeleteHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return projectFail(err)
	}
	if len(remaining) != 0 {
		return projectFail(fmt.Errorf("tsk project notes delete: unexpected args"))
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return projectFail(fmt.Errorf("tsk project notes delete: %v", err))
		}
	}
	if indexStr == "" {
		return projectFail(fmt.Errorf("tsk project notes delete: --index required"))
	}
	index, err := parseID(indexStr)
	if err != nil || index <= 0 {
		return projectFail(fmt.Errorf("tsk project notes delete: --index must be a positive integer"))
	}
	invk.setData(storage.EventData{
		Name:    firstNonEmpty(nameFlag, projectFlag),
		Project: firstNonEmpty(projectFlag, nameFlag),
		Index:   index,
		Labels:  labels,
	})

	notesDir, err := ensureNotesDir(home, dirFlag, nameFlag, projectFlag)
	if err != nil {
		return projectFail(err)
	}
	notes, err := storage.ReadTopicNotes(notesDir)
	if err != nil {
		return projectFail(err)
	}

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
		return projectFail(fmt.Errorf("tsk project notes delete: index %d out of range (have 0 %s)", index, word))
	}
	if index > len(filteredIdxs) {
		word := "notes"
		if len(filteredIdxs) == 1 {
			word = "note"
		}
		return projectFail(fmt.Errorf("tsk project notes delete: index %d out of range (have %d %s)", index, len(filteredIdxs), word))
	}

	targetIdx := filteredIdxs[index-1]
	notes = append(notes[:targetIdx], notes[targetIdx+1:]...)
	if err := storage.RewriteNotes(notesDir, notes); err != nil {
		return projectFail(err)
	}
	fmt.Println("deleted note")
	return nil
}
