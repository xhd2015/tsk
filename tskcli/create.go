package tskcli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runCreate(home string, args []string) error {
	setCommand(currentCtx, "create", args)

	var labels []string
	var topic string
	remaining, err := lessflags.
		StringSlice("--label", &labels).
		String("--topic", &topic).
		Parse(args)
	if err != nil {
		return fail(err)
	}
	if len(remaining) != 1 {
		return fail(fmt.Errorf("tsk create: title required"))
	}
	title := remaining[0]
	if title == "" {
		return fail(fmt.Errorf("tsk create: title required"))
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
	if topic != "" {
		topicParts = splitTopic(topic)
		relPath = storage.TopicRelPath(topic, id, stage, title)
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
		Stage:        stage,
		CreatedAt:    now,
		UpdatedAt:    now,
		StageHistory: []storage.StageHistoryEntry{},
	}
	if err := storage.WriteTask(taskDir, task); err != nil {
		return err
	}
	return storage.WriteIndex(home, id, relPath)
}