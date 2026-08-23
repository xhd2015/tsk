package storage

import (
	"os"
	"path/filepath"
	"sort"
)

// TopicTreeTask is one task node in a topic/inbox view (may nest sub-tasks).
type TopicTreeTask struct {
	ID    int            `json:"id"`
	Stage string         `json:"stage"`
	Slug  string         `json:"slug"`
	Dir   string         `json:"dir"`
	Tasks []TopicTreeTask `json:"tasks,omitempty"`
}

// TopicTree is a topic directory and its nested tasks / sub-topics.
type TopicTree struct {
	Path      string          `json:"path"`
	Aliases   []string        `json:"aliases"`
	Tasks     []TopicTreeTask `json:"tasks"`
	Subtopics []TopicTree     `json:"subtopics"`
}

// LoadTopicTree walks the topic directory (not task dirs, not json/jsonl files).
func LoadTopicTree(home string, parts []string) (TopicTree, error) {
	return loadTopicTreeAt(TopicAbs(home, parts), parts)
}

func loadTopicTreeAt(dir string, parts []string) (TopicTree, error) {
	meta, err := LoadTopicMeta(dir, parts)
	if err != nil {
		return TopicTree{}, err
	}
	if meta.Aliases == nil {
		meta.Aliases = []string{}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return TopicTree{}, err
	}
	tree := TopicTree{
		Path:      JoinTopicPath(parts),
		Aliases:   meta.Aliases,
		Tasks:     []TopicTreeTask{},
		Subtopics: []TopicTree{},
	}
	var subNames []string
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		if IsTaskDirName(name) {
			node, err := loadTaskTreeAt(filepath.Join(dir, name), name)
			if err != nil {
				return TopicTree{}, err
			}
			tree.Tasks = append(tree.Tasks, node)
			continue
		}
		subNames = append(subNames, name)
	}
	sort.Slice(tree.Tasks, func(i, j int) bool {
		return tree.Tasks[i].Dir < tree.Tasks[j].Dir
	})
	sort.Strings(subNames)
	for _, name := range subNames {
		childParts := append(append([]string(nil), parts...), name)
		child, err := loadTopicTreeAt(filepath.Join(dir, name), childParts)
		if err != nil {
			return TopicTree{}, err
		}
		child.Path = name
		tree.Subtopics = append(tree.Subtopics, child)
	}
	return tree, nil
}

func loadTaskTreeAt(dir, name string) (TopicTreeTask, error) {
	id, stage, slug, ok := ParseTaskDirName(name)
	if !ok {
		return TopicTreeTask{}, nil
	}
	node := TopicTreeTask{
		ID:    id,
		Stage: stage,
		Slug:  slug,
		Dir:   name,
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return TopicTreeTask{}, err
	}
	var kids []TopicTreeTask
	for _, ent := range entries {
		if !ent.IsDir() || !IsTaskDirName(ent.Name()) {
			continue
		}
		child, err := loadTaskTreeAt(filepath.Join(dir, ent.Name()), ent.Name())
		if err != nil {
			return TopicTreeTask{}, err
		}
		kids = append(kids, child)
	}
	sort.Slice(kids, func(i, j int) bool {
		return kids[i].Dir < kids[j].Dir
	})
	if len(kids) > 0 {
		node.Tasks = kids
	}
	return node, nil
}

// LoadInboxTasks returns task trees under inbox/ (topic_path == null).
func LoadInboxTasks(home string) ([]TopicTreeTask, error) {
	inboxDir := filepath.Join(home, "inbox")
	entries, err := os.ReadDir(inboxDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []TopicTreeTask{}, nil
		}
		return nil, err
	}
	var tasks []TopicTreeTask
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		if !IsTaskDirName(name) {
			continue
		}
		node, err := loadTaskTreeAt(filepath.Join(inboxDir, name), name)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, node)
	}
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Dir < tasks[j].Dir
	})
	return tasks, nil
}

// LoadTopicForest returns top-level topic trees under topics/.
func LoadTopicForest(home string) ([]TopicTree, error) {
	topicsDir := filepath.Join(home, "topics")
	entries, err := os.ReadDir(topicsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []TopicTree{}, nil
		}
		return nil, err
	}
	var names []string
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		if IsTaskDirName(name) {
			continue // a task dir at topics/ root, not a topic
		}
		names = append(names, name)
	}
	sort.Strings(names)
	forest := make([]TopicTree, 0, len(names))
	for _, name := range names {
		child, err := loadTopicTreeAt(filepath.Join(topicsDir, name), []string{name})
		if err != nil {
			return nil, err
		}
		child.Path = name
		forest = append(forest, child)
	}
	return forest, nil
}
