package tskcli

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runUpdate(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("update", args)

	var setProject, setTopic string
	var clearProject, clearTopic bool
	remaining, err := lessflags.
		String("--set-project", &setProject).
		Bool("--clear-project", &clearProject).
		String("--set-topic", &setTopic).
		Bool("--clear-topic", &clearTopic).
		Help("-h,--help", updateHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk update: task id required"))
	}
	if setProject == "" && !clearProject && setTopic == "" && !clearTopic {
		return fail(fmt.Errorf("tsk update: at least one of --set-project, --clear-project, --set-topic, --clear-topic required"))
	}
	if setProject != "" && clearProject {
		return fail(fmt.Errorf("tsk update: --set-project conflicts with --clear-project"))
	}
	if setTopic != "" && clearTopic {
		return fail(fmt.Errorf("tsk update: --set-topic conflicts with --clear-topic"))
	}

	id, err := parseID(remaining[0])
	if err != nil {
		return fail(err)
	}
	ev := storage.EventData{TaskID: id}
	if setTopic != "" {
		ev.Topic = setTopic
	}
	if setProject != "" {
		ev.Project = setProject
	}
	invk.setData(ev)
	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(err)
	}

	projectTouched := clearProject || setProject != ""
	if clearProject {
		task.Project = nil
	} else if setProject != "" {
		ref, err := resolveProjectRef(home, setProject)
		if err != nil {
			return fail(fmt.Errorf("tsk update: %w", err))
		}
		task.Project = ref
	}

	topicTouched := clearTopic || setTopic != ""
	if topicTouched {
		if task.ParentID != 0 {
			return fail(fmt.Errorf("tsk update: nested task %d: reparent before changing topic", id))
		}
		var topicParts []string
		if !clearTopic {
			topicParts, err = resolveTopicInput(home, setTopic)
			if err != nil {
				return fail(err)
			}
		}
		// MoveTaskDir writes topic_path and preserves in-memory project fields.
		if _, err = storage.MoveTaskDir(home, &task, taskDir, topicParts); err != nil {
			return fail(err)
		}
	} else if projectTouched {
		task.UpdatedAt = storage.NowTimestamp(task.ID)
		if err := storage.WriteTask(taskDir, task); err != nil {
			return fail(err)
		}
	}

	fmt.Printf("updated %d\n", id)
	return nil
}

// resolveProjectRef maps REF to a ProjectRef:
// registered name → exact origin → unique origin basename.
func resolveProjectRef(home, ref string) (*storage.ProjectRef, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("project ref required")
	}
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return nil, err
	}
	auto, err := storage.ReadProjectsAuto(home)
	if err != nil {
		return nil, err
	}

	if e, ok := storage.FindProjectByName(reg, ref); ok {
		if e.Origin != "" {
			return &storage.ProjectRef{Origin: e.Origin}, nil
		}
		if e.Name != "" {
			return &storage.ProjectRef{Name: e.Name}, nil
		}
	}
	if e, ok := storage.FindProjectByOrigin(reg, ref); ok {
		return &storage.ProjectRef{Origin: e.Origin}, nil
	}
	for _, e := range auto.Projects {
		if e.Origin == ref && e.Origin != "" {
			return &storage.ProjectRef{Origin: e.Origin}, nil
		}
		if e.Name == ref && e.Origin == "" && e.Name != "" {
			return &storage.ProjectRef{Name: e.Name}, nil
		}
	}
	if strings.Contains(ref, "/") {
		return &storage.ProjectRef{Origin: ref}, nil
	}

	var matches []string
	seen := map[string]struct{}{}
	consider := func(origin string) {
		if origin == "" {
			return
		}
		if _, ok := seen[origin]; ok {
			return
		}
		base := filepath.Base(origin)
		if _, b, err := NormalizeOriginURL("https://" + origin); err == nil && b != "" {
			base = b
		}
		if base == ref || filepath.Base(origin) == ref {
			seen[origin] = struct{}{}
			matches = append(matches, origin)
		}
	}
	for _, e := range reg.Projects {
		consider(e.Origin)
	}
	for _, e := range auto.Projects {
		consider(e.Origin)
	}
	if len(matches) == 1 {
		return &storage.ProjectRef{Origin: matches[0]}, nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("ambiguous project %q matches %d origins; pass a full origin", ref, len(matches))
	}
	return nil, fmt.Errorf("unknown project %q", ref)
}
