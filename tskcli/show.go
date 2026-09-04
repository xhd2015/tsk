package tskcli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/pathfmt"
	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runShow(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("show", args)

	remaining, err := lessflags.
		Help("-h,--help", showHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk show: task id required"))
	}
	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{TaskID: id})

	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	topicParts, err := storage.ParseTopicPath(task.TopicPath)
	if err != nil {
		return err
	}
	var topicStr string
	if len(topicParts) == 0 {
		topicStr = "inbox"
	} else {
		topicStr = strings.Join(topicParts, "/")
	}

	fmt.Printf("id: %d\n", task.ID)
	fmt.Printf("title: %s\n", task.Title)
	fmt.Printf("slug: %s\n", task.Slug)
	fmt.Printf("stage: %s\n", task.Stage)
	fmt.Printf("topic: %s\n", topicStr)
	if task.ParentID != 0 {
		fmt.Printf("parent: %d\n", task.ParentID)
	}
	if task.Cwd != "" {
		fmt.Printf("cwd: %s\n", pathfmt.TildeHome(pathfmt.Expand(task.Cwd)))
	}
	if task.Project != nil {
		if line := formatShowProject(home, task.Project); line != "" {
			fmt.Printf("project: %s\n", line)
		}
	}
	if len(task.Labels) == 0 {
		fmt.Println("labels:")
	} else {
		fmt.Printf("labels: %s\n", strings.Join(task.Labels, ", "))
	}
	fmt.Printf("created_at: %s\n", task.CreatedAt)
	fmt.Printf("updated_at: %s\n", task.UpdatedAt)
	journal, err := storage.ReadTopicNotes(taskDir)
	if err != nil {
		return fail(err)
	}
	fmt.Printf("notes: %d\n", len(journal))

	var progressNotes []storage.TopicNote
	for _, n := range journal {
		if storage.NoteHasAllLabels(n, []string{"progress"}) {
			progressNotes = append(progressNotes, n)
		}
	}
	if len(progressNotes) > 0 {
		latest := progressNotes[len(progressNotes)-1]
		word := "entries"
		if len(progressNotes) == 1 {
			word = "entry"
		}
		if latest.Status != "" {
			fmt.Printf("progress: %s (%d %s)\n", latest.Status, len(progressNotes), word)
		} else {
			fmt.Printf("progress: %d %s\n", len(progressNotes), word)
		}
	}
	return nil
}

// formatShowProject picks the project: line value: ledger location (tilde),
// else task/registry name, else origin.
func formatShowProject(home string, ref *storage.ProjectRef) string {
	if ref == nil {
		return ""
	}
	if loc := lookupLedgerLocation(home, ref); loc != "" {
		return pathfmt.TildeHome(pathfmt.Expand(loc))
	}
	if ref.Name != "" {
		return ref.Name
	}
	return ref.Origin
}

func lookupLedgerLocation(home string, ref *storage.ProjectRef) string {
	if ref == nil {
		return ""
	}
	reg, err := storage.ReadProjects(home)
	if err == nil {
		if ref.Origin != "" {
			if e, ok := storage.FindProjectByOrigin(reg, ref.Origin); ok {
				if loc := e.EffectiveLocation(); loc != "" {
					return loc
				}
			}
		}
		if ref.Name != "" {
			if e, ok := storage.FindProjectByName(reg, ref.Name); ok {
				if loc := e.EffectiveLocation(); loc != "" {
					return loc
				}
			}
		}
	}
	auto, err := storage.ReadProjectsAuto(home)
	if err != nil {
		return ""
	}
	if e, ok := storage.FindProjectAuto(auto, *ref); ok {
		return e.EffectiveLocation()
	}
	return ""
}
