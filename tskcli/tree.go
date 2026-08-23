package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runTree(home string, args []string) error {
	setCommand(currentCtx, "tree", args)

	var asJSON bool
	var colorFlag bool
	var plain bool
	var idStr string
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Bool("--color", &colorFlag).
		Bool("--plain", &plain).
		String("--id", &idStr).
		Help("-h,--help", treeHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 0 {
		return fail(fmt.Errorf("tsk tree: unexpected arguments"))
	}
	color := !asJSON && !plain && (colorFlag || isStdoutTTY())

	if idStr != "" {
		id, err := parseID(idStr)
		if err != nil {
			return fail(err)
		}
		return runTreeID(home, id, asJSON, color)
	}

	inbox, err := storage.LoadInboxTasks(home)
	if err != nil {
		return fail(err)
	}
	forest, err := storage.LoadTopicForest(home)
	if err != nil {
		return fail(err)
	}

	if asJSON {
		root := treeJSONRoot{Inbox: inbox, Topics: forest}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(root)
	}
	fmt.Print(formatTree(inbox, forest, color))
	return nil
}

// runTreeID renders a pruned branch from root to one task, with its notes
// and progress entries nested under the task leaf.
func runTreeID(home string, id int, asJSON, color bool) error {
	task, taskDir, err := storage.LoadTaskByID(home, id)
	if err != nil {
		return fail(fmt.Errorf("Error: tsk tree: task not found: %d", id))
	}

	topicParts, err := storage.ParseTopicPath(task.TopicPath)
	if err != nil {
		return fail(err)
	}

	notes, err := storage.ReadTopicNotes(taskDir)
	if err != nil {
		return fail(err)
	}

	var regularNotes, progressNotes []storage.TopicNote
	for _, n := range notes {
		if storage.NoteHasAllLabels(n, []string{"progress"}) {
			progressNotes = append(progressNotes, n)
		} else {
			regularNotes = append(regularNotes, n)
		}
	}

	if asJSON {
		topicPath := topicParts
		if topicPath == nil {
			topicPath = []string{}
		}
		if regularNotes == nil {
			regularNotes = []storage.TopicNote{}
		}
		if progressNotes == nil {
			progressNotes = []storage.TopicNote{}
		}
		root := treeIDJSON{
			Task: treeIDTask{
				ID:        task.ID,
				Stage:     task.Stage,
				Slug:      task.Slug,
				Dir:       filepath.Base(taskDir),
				TopicPath: topicPath,
			},
			Notes:    regularNotes,
			Progress: progressNotes,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(root)
	}

	fmt.Print(formatTreeID(home, task, taskDir, topicParts, regularNotes, progressNotes, color))
	return nil
}

type treeIDJSON struct {
	Task     treeIDTask          `json:"task"`
	Notes    []storage.TopicNote `json:"notes"`
	Progress []storage.TopicNote `json:"progress"`
}

type treeIDTask struct {
	ID        int      `json:"id"`
	Stage     string   `json:"stage"`
	Slug      string   `json:"slug"`
	Dir       string   `json:"dir"`
	TopicPath []string `json:"topic_path"`
}

// renderNode is a generic tree node for rendering arbitrary nesting.
type renderNode struct {
	name     string
	extra    string
	styled   bool
	children []*renderNode
}

func formatTreeID(home string, task storage.Task, taskDir string, topicParts []string, regularNotes, progressNotes []storage.TopicNote, color bool) string {
	var b strings.Builder
	b.WriteString(".\n")

	taskNode := &renderNode{
		name:   filepath.Base(taskDir),
		extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
		styled: color && isTerminalTaskStage(task.Stage),
	}

	// build notes/progress sections under the task
	var sections []*renderNode
	if len(regularNotes) > 0 {
		notesNode := &renderNode{name: "notes"}
		for _, n := range regularNotes {
			notesNode.children = append(notesNode.children, &renderNode{name: formatNoteLine(n)})
		}
		sections = append(sections, notesNode)
	}
	if len(progressNotes) > 0 {
		progNode := &renderNode{name: "progress"}
		for _, n := range progressNotes {
			progNode.children = append(progNode.children, &renderNode{name: formatProgressLine(n, color)})
		}
		sections = append(sections, progNode)
	}
	taskNode.children = sections

	var rootKids []*renderNode
	if len(topicParts) == 0 {
		// inbox task: direct root child
		rootKids = []*renderNode{taskNode}
	} else {
		// build nested topic chain
		chain := buildTopicChain(home, topicParts)
		// attach task to the deepest node
		deepest := chain
		for len(deepest.children) > 0 {
			deepest = deepest.children[0]
		}
		deepest.children = []*renderNode{taskNode}
		rootKids = []*renderNode{chain}
	}

	writeRenderKids(&b, rootKids, "")

	topicCount := len(topicParts)
	topicWord := "topics"
	if topicCount == 1 {
		topicWord = "topic"
	}
	b.WriteString(fmt.Sprintf("\n1 task, %d %s\n", topicCount, topicWord))
	return b.String()
}

// buildTopicChain creates a nested renderNode chain for a topic path,
// loading aliases from topic.json at each level.
func buildTopicChain(home string, parts []string) *renderNode {
	var root *renderNode
	var current *renderNode
	for i, part := range parts {
		subParts := parts[:i+1]
		topicDir := storage.TopicAbs(home, subParts)
		meta, _ := storage.LoadTopicMeta(topicDir, subParts)
		node := &renderNode{name: part}
		if len(meta.Aliases) > 0 {
			node.extra = "  aliases: " + strings.Join(meta.Aliases, ", ")
		}
		if root == nil {
			root = node
		} else {
			current.children = []*renderNode{node}
		}
		current = node
	}
	return root
}

func writeRenderKids(b *strings.Builder, kids []*renderNode, prefix string) {
	for i, kid := range kids {
		last := i == len(kids)-1
		branch, next := "├── ", "│   "
		if last {
			branch, next = "└── ", "    "
		}
		b.WriteString(prefix)
		b.WriteString(branch)
		if kid.styled {
			b.WriteString(ansiGray + ansiStrikethrough)
		}
		b.WriteString(kid.name)
		b.WriteString(kid.extra)
		if kid.styled {
			b.WriteString(ansiReset)
		}
		b.WriteByte('\n')
		if len(kid.children) > 0 {
			writeRenderKids(b, kid.children, prefix+next)
		}
	}
}

type treeJSONRoot struct {
	Inbox  []storage.TopicTreeTask `json:"inbox"`
	Topics []storage.TopicTree     `json:"topics"`
}

func isTerminalTaskStage(stage string) bool {
	return stage == "done"
}

func formatTree(inbox []storage.TopicTreeTask, forest []storage.TopicTree, color bool) string {
	var b strings.Builder
	b.WriteString(".\n")

	rootKids := treeRootChildren(inbox, forest, color)
	if len(rootKids) == 0 {
		b.WriteString("(empty)\n")
		b.WriteString("0 tasks, 0 topics\n")
		return b.String()
	}
	writeViewKids(&b, rootKids, "")

	taskCount := 0
	for i := range inbox {
		taskCount += countTaskNode(&inbox[i])
	}
	for i := range forest {
		taskCount += countTreeTasks(&forest[i])
	}
	taskWord := "tasks"
	if taskCount == 1 {
		taskWord = "task"
	}
	topicWord := "topics"
	if len(forest) == 1 {
		topicWord = "topic"
	}
	b.WriteString(fmt.Sprintf("\n%d %s, %d %s\n", taskCount, taskWord, len(forest), topicWord))
	return b.String()
}

func treeRootChildren(inbox []storage.TopicTreeTask, forest []storage.TopicTree, color bool) []viewNode {
	out := make([]viewNode, 0, len(inbox)+len(forest))
	for _, task := range inbox {
		out = append(out, viewNode{
			name:   task.Dir,
			extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
			styled: color && isTerminalTaskStage(task.Stage),
			tasks:  task.Tasks,
			color:  color,
		})
	}
	for i := range forest {
		st := &forest[i]
		extra := ""
		if len(st.Aliases) > 0 {
			extra = "  aliases: " + strings.Join(st.Aliases, ", ")
		}
		out = append(out, viewNode{name: st.Path, extra: extra, nested: st, color: color})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func countTreeTasks(t *storage.TopicTree) int {
	n := 0
	for i := range t.Tasks {
		n += countTaskNode(&t.Tasks[i])
	}
	for i := range t.Subtopics {
		n += countTreeTasks(&t.Subtopics[i])
	}
	return n
}

func countTaskNode(t *storage.TopicTreeTask) int {
	n := 1
	for i := range t.Tasks {
		n += countTaskNode(&t.Tasks[i])
	}
	return n
}
