package tskcli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runAdd(home string, args []string) error {
	setCommand(currentCtx, "add", args)

	var labels []string
	var notes []string
	var topic string
	var parentStr string
	remaining, err := lessflags.
		StringSlice("--label", &labels).
		StringSlice("--note", &notes).
		String("--topic", &topic).
		String("--parent", &parentStr).
		Help("-h,--help", addHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk add: title required"))
	}
	title := remaining[0]
	if title == "" {
		return fail(fmt.Errorf("tsk add: title required"))
	}
	if topic != "" && parentStr != "" {
		return fail(fmt.Errorf("tsk add: --topic conflicts with --parent (child inherits parent location)"))
	}
	for i, text := range notes {
		notes[i] = strings.TrimSpace(text)
		if notes[i] == "" {
			return fail(fmt.Errorf("tsk add: --note text required"))
		}
	}

	if err := storage.EnsureLayout(home); err != nil {
		return err
	}

	id, err := storage.NextID(home)
	if err != nil {
		return err
	}

	slug := storage.Slugify(title)
	stage := "create"
	now := storage.NowTimestamp(id)

	sort.Strings(labels)
	unique := labels[:0]
	for i, l := range labels {
		if i == 0 || l != labels[i-1] {
			unique = append(unique, l)
		}
	}
	labels = unique

	var topicParts []string
	var relPath string
	var parentID int
	if parentStr != "" {
		parentID, err = parseID(parentStr)
		if err != nil {
			return fail(fmt.Errorf("tsk add: parent task not found: %s", parentStr))
		}
		parent, parentDir, loadErr := storage.LoadTaskByID(home, parentID)
		if loadErr != nil {
			return fail(fmt.Errorf("tsk add: parent task not found: %d", parentID))
		}
		parentRel, err := storage.RelFromHome(home, parentDir)
		if err != nil {
			return fail(err)
		}
		topicParts, err = storage.ParseTopicPath(parent.TopicPath)
		if err != nil {
			return fail(err)
		}
		relPath = storage.ChildRelPath(parentRel, id, stage, slug)
	} else if topic != "" {
		topicParts, err = resolveTopicInput(home, topic)
		if err != nil {
			return fail(err)
		}
		relPath = storage.TopicRelPath(storage.JoinTopicPath(topicParts), id, stage, title)
		if err := os.MkdirAll(filepath.Join(home, filepath.Dir(filepath.FromSlash(relPath))), 0o755); err != nil {
			return err
		}
	} else {
		relPath = storage.InboxRelPath(id, stage, title)
	}

	taskDir := filepath.Join(home, filepath.FromSlash(relPath))
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(taskDir, "context"), 0o755); err != nil {
		return err
	}

	topicJSON, err := storage.TopicPathJSON(topicParts)
	if err != nil {
		return err
	}

	task := storage.Task{
		ID:           id,
		Title:        title,
		Slug:         slug,
		Labels:       labels,
		TopicPath:    topicJSON,
		ParentID:     parentID,
		Stage:        stage,
		CreatedAt:    now,
		UpdatedAt:    now,
		StageHistory: []storage.StageHistoryEntry{},
	}
	if err := storage.WriteTask(taskDir, task); err != nil {
		return err
	}
	if err := storage.WriteIndex(home, id, relPath); err != nil {
		return err
	}
	for _, text := range notes {
		existing, err := storage.ReadTopicNotes(taskDir)
		if err != nil {
			return fail(err)
		}
		note := storage.TopicNote{
			TS:   storage.NowTimestamp(len(existing) + 1),
			Text: text,
		}
		if err := storage.AppendTopicNote(taskDir, note); err != nil {
			return fail(fmt.Errorf("tsk add: append note: %w", err))
		}
	}
	fmt.Println(id)
	return nil
}
