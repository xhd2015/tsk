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
	fancy := !asJSON && !plain && isStdoutTTY()

	if idStr != "" {
		id, err := parseID(idStr)
		if err != nil {
			return fail(err)
		}
		return runTreeID(home, id, asJSON, color, fancy)
	}

	inbox, err := storage.LoadInboxTasks(home)
	if err != nil {
		return fail(err)
	}
	forest, err := storage.LoadTopicForest(home)
	if err != nil {
		return fail(err)
	}
	reg, err := storage.ReadProjects(home)
	if err != nil {
		return fail(err)
	}

	if asJSON {
		root := buildTreeJSON(inbox, forest, reg)
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(root)
	}
	fmt.Print(formatTree(inbox, forest, reg, color, fancy))
	return nil
}

// runTreeID renders a pruned branch from root to one task, with its notes
// and progress entries nested under the task leaf.
func runTreeID(home string, id int, asJSON, color, fancy bool) error {
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
				Project:   task.Project,
			},
			Notes:    regularNotes,
			Progress: progressNotes,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(root)
	}

	reg, _ := storage.ReadProjects(home)
	fmt.Print(formatTreeID(home, task, taskDir, topicParts, regularNotes, progressNotes, reg, color, fancy))
	return nil
}

type treeIDJSON struct {
	Task     treeIDTask          `json:"task"`
	Notes    []storage.TopicNote `json:"notes"`
	Progress []storage.TopicNote `json:"progress"`
}

type treeIDTask struct {
	ID        int                 `json:"id"`
	Stage     string              `json:"stage"`
	Slug      string              `json:"slug"`
	Dir       string              `json:"dir"`
	TopicPath []string            `json:"topic_path"`
	Project   *storage.ProjectRef `json:"project,omitempty"`
}

// renderNode is a generic tree node for rendering arbitrary nesting.
type renderNode struct {
	name     string
	extra    string
	styled   bool
	children []*renderNode
}

func formatTreeID(home string, task storage.Task, taskDir string, topicParts []string, regularNotes, progressNotes []storage.TopicNote, reg storage.ProjectsFile, color, fancy bool) string {
	var b strings.Builder
	b.WriteString(".\n")

	taskNode := &renderNode{
		name:   filepath.Base(taskDir),
		extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
		styled: color && isTerminalTaskStage(task.Stage),
	}

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

	leaf := taskNode
	projCount := 0
	if task.Project != nil && (task.Project.Origin != "" || task.Project.Name != "") {
		projCount = 1
		leaf = &renderNode{
			name:     projectMark(fancy) + projectDisplayLabel(task.Project, reg),
			children: []*renderNode{taskNode},
		}
	}

	var rootKids []*renderNode
	if len(topicParts) == 0 {
		rootKids = []*renderNode{leaf}
	} else {
		chain := buildTopicChain(home, topicParts, fancy)
		deepest := chain
		for len(deepest.children) > 0 {
			deepest = deepest.children[0]
		}
		deepest.children = []*renderNode{leaf}
		rootKids = []*renderNode{chain}
	}

	writeRenderKids(&b, rootKids, "")

	topicCount := len(topicParts)
	b.WriteString(formatTreeFooter(1, topicCount, projCount))
	return b.String()
}

// buildTopicChain creates a nested renderNode chain for a topic path,
// loading aliases from topic.json at each level.
func buildTopicChain(home string, parts []string, fancy bool) *renderNode {
	var root *renderNode
	var current *renderNode
	mark := topicMark(fancy)
	for i, part := range parts {
		subParts := parts[:i+1]
		topicDir := storage.TopicAbs(home, subParts)
		meta, _ := storage.LoadTopicMeta(topicDir, subParts)
		node := &renderNode{name: mark + part}
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
	Inbox         []storage.TopicTreeTask     `json:"inbox"`
	InboxProjects []storage.TopicProjectGroup `json:"inbox_projects"`
	Topics        []storage.TopicTree         `json:"topics"`
}

func isTerminalTaskStage(stage string) bool {
	return stage == "done"
}

func topicMark(fancy bool) string {
	if fancy {
		return "▣ "
	}
	return "# "
}

func projectMark(fancy bool) string {
	if fancy {
		return "◆ "
	}
	return "@ "
}

func formatTree(inbox []storage.TopicTreeTask, forest []storage.TopicTree, reg storage.ProjectsFile, color, fancy bool) string {
	var b strings.Builder
	b.WriteString(".\n")

	rootKids, projCount := treeRootChildren(inbox, forest, reg, color, fancy)
	if len(rootKids) == 0 {
		b.WriteString("(empty)\n")
		b.WriteString(strings.TrimPrefix(formatTreeFooter(0, 0, 0), "\n"))
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
	b.WriteString(formatTreeFooter(taskCount, len(forest), projCount))
	return b.String()
}

func formatTreeFooter(tasks, topics, projects int) string {
	taskWord := "tasks"
	if tasks == 1 {
		taskWord = "task"
	}
	topicWord := "topics"
	if topics == 1 {
		topicWord = "topic"
	}
	projectWord := "projects"
	if projects == 1 {
		projectWord = "project"
	}
	return fmt.Sprintf("\n%d %s, %d %s, %d %s\n", tasks, taskWord, topics, topicWord, projects, projectWord)
}

func treeRootChildren(inbox []storage.TopicTreeTask, forest []storage.TopicTree, reg storage.ProjectsFile, color, fancy bool) ([]viewNode, int) {
	out := make([]viewNode, 0, len(inbox)+len(forest))
	projCount := 0

	ungrouped, groups := partitionTasksByProject(inbox, reg)
	for _, g := range groups {
		projCount++
		out = append(out, makeProjectViewNode(g, color, fancy))
	}
	for _, task := range ungrouped {
		out = append(out, makeTaskViewNode(task, color))
	}
	for i := range forest {
		node, pc := makeTopicViewNode(&forest[i], reg, color, fancy)
		projCount += pc
		out = append(out, node)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, projCount
}

func makeTaskViewNode(task storage.TopicTreeTask, color bool) viewNode {
	return viewNode{
		name:   task.Dir,
		extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
		styled: color && isTerminalTaskStage(task.Stage),
		tasks:  task.Tasks,
		color:  color,
	}
}

func makeProjectViewNode(g storage.TopicProjectGroup, color, fancy bool) viewNode {
	kids := make([]viewNode, 0, len(g.Tasks))
	for _, t := range g.Tasks {
		kids = append(kids, makeTaskViewNode(t, color))
	}
	sort.Slice(kids, func(i, j int) bool { return kids[i].name < kids[j].name })
	return viewNode{
		name:  g.Label,
		mark:  projectMark(fancy),
		kids:  kids,
		color: color,
	}
}

func makeTopicViewNode(tree *storage.TopicTree, reg storage.ProjectsFile, color, fancy bool) (viewNode, int) {
	extra := ""
	if len(tree.Aliases) > 0 {
		extra = "  aliases: " + strings.Join(tree.Aliases, ", ")
	}
	kids, pc := topicLevelChildren(tree, reg, color, fancy)
	return viewNode{
		name:  tree.Path,
		mark:  topicMark(fancy),
		extra: extra,
		kids:  kids,
		color: color,
	}, pc
}

func topicLevelChildren(tree *storage.TopicTree, reg storage.ProjectsFile, color, fancy bool) ([]viewNode, int) {
	out := make([]viewNode, 0, len(tree.Tasks)+len(tree.Subtopics))
	projCount := 0
	ungrouped, groups := partitionTasksByProject(tree.Tasks, reg)
	for _, g := range groups {
		projCount++
		out = append(out, makeProjectViewNode(g, color, fancy))
	}
	for _, task := range ungrouped {
		out = append(out, makeTaskViewNode(task, color))
	}
	for i := range tree.Subtopics {
		node, pc := makeTopicViewNode(&tree.Subtopics[i], reg, color, fancy)
		projCount += pc
		out = append(out, node)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out, projCount
}

type projectBucket struct {
	origin string
	name   string
	label  string
	tasks  []storage.TopicTreeTask
}

func partitionTasksByProject(tasks []storage.TopicTreeTask, reg storage.ProjectsFile) (ungrouped []storage.TopicTreeTask, groups []storage.TopicProjectGroup) {
	buckets := map[string]*projectBucket{}
	var order []string
	for _, t := range tasks {
		if t.Project == nil || (t.Project.Origin == "" && t.Project.Name == "") {
			ungrouped = append(ungrouped, t)
			continue
		}
		key := storage.ProjectAutoKey(t.Project.Origin, t.Project.Name)
		b, ok := buckets[key]
		if !ok {
			b = &projectBucket{
				origin: t.Project.Origin,
				name:   t.Project.Name,
				label:  projectDisplayLabel(t.Project, reg),
			}
			buckets[key] = b
			order = append(order, key)
		}
		b.tasks = append(b.tasks, t)
	}
	sort.Strings(order)
	for _, key := range order {
		b := buckets[key]
		sort.Slice(b.tasks, func(i, j int) bool { return b.tasks[i].Dir < b.tasks[j].Dir })
		groups = append(groups, storage.TopicProjectGroup{
			Origin: b.origin,
			Name:   displayNameForProjectRef(b.origin, b.name, reg),
			Label:  b.label,
			Tasks:  b.tasks,
		})
	}
	sort.Slice(ungrouped, func(i, j int) bool { return ungrouped[i].Dir < ungrouped[j].Dir })
	return ungrouped, groups
}

func projectDisplayLabel(ref *storage.ProjectRef, reg storage.ProjectsFile) string {
	if ref == nil {
		return ""
	}
	if ref.Origin != "" {
		name := displayNameForProjectRef(ref.Origin, ref.Name, reg)
		if name == "" {
			name = filepath.Base(ref.Origin)
		}
		return name + "  " + ref.Origin
	}
	return ref.Name
}

func displayNameForProjectRef(origin, name string, reg storage.ProjectsFile) string {
	if name != "" {
		return name
	}
	if origin == "" {
		return ""
	}
	if e, ok := storage.FindProjectByOrigin(reg, origin); ok && e.Name != "" {
		return e.Name
	}
	_, base, err := NormalizeOriginURL("https://" + origin)
	if err == nil && base != "" {
		return base
	}
	return filepath.Base(origin)
}

func buildTreeJSON(inbox []storage.TopicTreeTask, forest []storage.TopicTree, reg storage.ProjectsFile) treeJSONRoot {
	ungrouped, groups := partitionTasksByProject(inbox, reg)
	topics := make([]storage.TopicTree, len(forest))
	for i := range forest {
		topics[i] = topicTreeForJSON(&forest[i], reg)
	}
	if ungrouped == nil {
		ungrouped = []storage.TopicTreeTask{}
	}
	if groups == nil {
		groups = []storage.TopicProjectGroup{}
	}
	return treeJSONRoot{
		Inbox:         ungrouped,
		InboxProjects: groups,
		Topics:        topics,
	}
}

func topicTreeForJSON(tree *storage.TopicTree, reg storage.ProjectsFile) storage.TopicTree {
	ungrouped, groups := partitionTasksByProject(tree.Tasks, reg)
	subs := make([]storage.TopicTree, len(tree.Subtopics))
	for i := range tree.Subtopics {
		subs[i] = topicTreeForJSON(&tree.Subtopics[i], reg)
	}
	if ungrouped == nil {
		ungrouped = []storage.TopicTreeTask{}
	}
	out := storage.TopicTree{
		Path:      tree.Path,
		Aliases:   tree.Aliases,
		Tasks:     ungrouped,
		Projects:  groups,
		Subtopics: subs,
	}
	if out.Aliases == nil {
		out.Aliases = []string{}
	}
	if out.Projects == nil {
		out.Projects = []storage.TopicProjectGroup{}
	}
	if out.Subtopics == nil {
		out.Subtopics = []storage.TopicTree{}
	}
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
