package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runTopicView(home string, args []string) error {
	setCommand(currentCtx, "topic", append([]string{"view"}, args...))

	var asJSON bool
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Help("-h,--help", topicViewHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return topicErr("tsk topic view: topic required")
	}
	parts, _, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	tree, err := storage.LoadTopicTree(home, parts)
	if err != nil {
		return fail(err)
	}
	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(tree)
	}
	fmt.Print(formatTopicView(tree))
	return nil
}

type viewNode struct {
	name   string
	extra  string
	nested *storage.TopicTree
	tasks  []storage.TopicTreeTask
	color  bool
	styled bool
}

func formatTopicView(tree storage.TopicTree) string {
	var b strings.Builder
	b.WriteString(tree.Path)
	if len(tree.Aliases) > 0 {
		b.WriteString("  aliases: ")
		b.WriteString(strings.Join(tree.Aliases, ", "))
	}
	b.WriteByte('\n')
	kids := viewChildren(&tree, false)
	if len(kids) == 0 {
		b.WriteString("(empty)\n")
		return b.String()
	}
	writeViewKids(&b, kids, "")
	return b.String()
}

func viewChildren(tree *storage.TopicTree, color bool) []viewNode {
	out := make([]viewNode, 0, len(tree.Tasks)+len(tree.Subtopics))
	for _, task := range tree.Tasks {
		out = append(out, viewNode{
			name:   task.Dir,
			extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
			styled: color && isTerminalTaskStage(task.Stage),
			tasks:  task.Tasks,
			color:  color,
		})
	}
	for i := range tree.Subtopics {
		st := &tree.Subtopics[i]
		out = append(out, viewNode{name: st.Path, nested: st, color: color})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func taskViewChildren(tasks []storage.TopicTreeTask, color bool) []viewNode {
	out := make([]viewNode, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, viewNode{
			name:   task.Dir,
			extra:  fmt.Sprintf("  task %d  %s", task.ID, task.Stage),
			styled: color && isTerminalTaskStage(task.Stage),
			tasks:  task.Tasks,
			color:  color,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

func writeViewKids(b *strings.Builder, kids []viewNode, prefix string) {
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
		var nested []viewNode
		if kid.nested != nil {
			nested = viewChildren(kid.nested, kid.color)
		} else if len(kid.tasks) > 0 {
			nested = taskViewChildren(kid.tasks, kid.color)
		}
		if len(nested) > 0 {
			writeViewKids(b, nested, prefix+next)
		}
	}
}
