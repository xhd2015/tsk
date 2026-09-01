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
	fmt.Print(formatTopicView(tree, treeColorEnabled(false, false, asJSON)))
	return nil
}

type viewNode struct {
	name  string
	extra string
	mark  string // kind prefix e.g. "▣ " / "# "; drawn dim when color
	style string // ANSI SGR prefix for name+extra; empty = unstyled
	kids  []viewNode // explicit children (project groups, tasks, subtopics)
	color bool
}

func formatTopicView(tree storage.TopicTree, color bool) string {
	var b strings.Builder
	b.WriteString(tree.Path)
	if len(tree.Aliases) > 0 {
		b.WriteString("  aliases: ")
		b.WriteString(strings.Join(tree.Aliases, ", "))
	}
	b.WriteByte('\n')
	idWidth := maxTopicTreeIDWidth(&tree)
	kids := viewChildren(&tree, color, idWidth)
	if len(kids) == 0 {
		b.WriteString("(empty)\n")
		return b.String()
	}
	writeViewKids(&b, kids, "")
	return b.String()
}

func maxTopicTreeIDWidth(tree *storage.TopicTree) int {
	max := maxTaskIDWidth(tree.Tasks)
	for i := range tree.Subtopics {
		if w := maxTopicTreeIDWidth(&tree.Subtopics[i]); w > max {
			max = w
		}
	}
	return max
}

func viewChildren(tree *storage.TopicTree, color bool, idWidth int) []viewNode {
	out := make([]viewNode, 0, len(tree.Tasks)+len(tree.Subtopics))
	for _, task := range tree.Tasks {
		out = append(out, makeTaskViewNode(task, color, idWidth))
	}
	for i := range tree.Subtopics {
		st := &tree.Subtopics[i]
		out = append(out, viewNode{
			name:  st.Path,
			kids:  viewChildren(st, color, idWidth),
			color: color,
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
		if kid.style != "" {
			b.WriteString(kid.style)
		}
		if kid.mark != "" {
			if kid.color && kid.style == "" {
				b.WriteString(ansiGray)
				b.WriteString(kid.mark)
				b.WriteString(ansiReset)
			} else {
				b.WriteString(kid.mark)
			}
		}
		b.WriteString(kid.name)
		b.WriteString(kid.extra)
		if kid.style != "" {
			b.WriteString(ansiReset)
		}
		b.WriteByte('\n')
		if len(kid.kids) > 0 {
			writeViewKids(b, kid.kids, prefix+next)
		}
	}
}
