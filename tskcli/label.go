package tskcli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runLabel(home string, args []string) error {
	setCommand(currentCtx, "label", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(labelHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runLabelAdd(home, args[1:])
	case "rm":
		return runLabelRm(home, args[1:])
	case "list":
		return runLabelList(home, args[1:])
	default:
		return fail(fmt.Errorf("tsk label: unknown subcommand %q", args[0]))
	}
}

func runLabelList(home string, args []string) error {
	setCommand(currentCtx, "label", append([]string{"list"}, args...))

	remaining, err := lessflags.
		Help("-h,--help", labelListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return fail(fmt.Errorf("tsk label list: unexpected arguments"))
	}

	names, err := collectLabelNames(home)
	if err != nil {
		return fail(err)
	}
	for _, name := range names {
		fmt.Println(name)
	}
	word := "labels"
	if len(names) == 1 {
		word = "label"
	}
	fmt.Printf("%d %s\n", len(names), word)
	return nil
}

func collectLabelNames(home string) ([]string, error) {
	seen := make(map[string]struct{})
	add := func(tokens []string) {
		for _, tok := range tokens {
			name, err := storage.LabelName(tok)
			if err != nil || name == "" {
				continue
			}
			seen[name] = struct{}{}
		}
	}

	ids, err := storage.ListTaskIDs(home)
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		task, dir, err := storage.LoadTaskByID(home, id)
		if err != nil {
			return nil, err
		}
		add(task.Labels)
		notes, err := storage.ReadTopicNotes(dir)
		if err != nil {
			return nil, err
		}
		for _, n := range notes {
			add(n.Labels)
		}
	}

	root := filepath.Join(home, "topics")
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == root {
				return nil
			}
			return err
		}
		if d.IsDir() {
			if path == root {
				return nil
			}
			if storage.IsTaskDirName(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != storage.TopicNotesFile {
			return nil
		}
		notes, err := storage.ReadTopicNotes(filepath.Dir(path))
		if err != nil {
			return err
		}
		for _, n := range notes {
			add(n.Labels)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out, nil
}

func runLabelAdd(home string, args []string) error {
	setCommand(currentCtx, "label", append([]string{"add"}, args...))

	if len(args) != 2 {
		return fail(fmt.Errorf("tsk label add: task id and label required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	label := args[1]

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	if !containsLabel(task.Labels, label) {
		task.Labels = append(task.Labels, label)
		sort.Strings(task.Labels)
		task.UpdatedAt = storage.NowTimestamp(task.ID)
		return storage.WriteTask(taskDir, task)
	}
	return nil
}

func runLabelRm(home string, args []string) error {
	setCommand(currentCtx, "label", append([]string{"rm"}, args...))

	if len(args) != 2 {
		return fail(fmt.Errorf("tsk label rm: task id and label required"))
	}
	id, err := parseID(args[0])
	if err != nil {
		return fail(err)
	}
	label := args[1]

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}
	var kept []string
	for _, l := range task.Labels {
		if l != label {
			kept = append(kept, l)
		}
	}
	task.Labels = kept
	task.UpdatedAt = storage.NowTimestamp(task.ID)
	return storage.WriteTask(taskDir, task)
}