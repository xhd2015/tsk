package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

func runTopicWhere(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"where"}, args...))

	remaining, err := lessflags.
		Help("-h,--help", topicWhereHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return topicErr("tsk topic where: topic required")
	}
	parts, dir, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts)})
	fmt.Println(dir)
	return nil
}

type topicInfoJSON struct {
	Path      string   `json:"path"`
	Title     string   `json:"title"`
	Aliases   []string `json:"aliases"`
	Dir       string   `json:"dir"`
	Notes     int      `json:"notes"`
	Tasks     int      `json:"tasks"`
	Subtopics []string `json:"subtopics"`
}

func runTopicInfo(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"info"}, args...))

	var asJSON bool
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Help("-h,--help", topicInfoHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return topicErr("tsk topic info: topic required")
	}
	parts, dir, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts)})
	if err := storage.MigrateLegacyTopicNotes(dir, parts); err != nil {
		return fail(err)
	}
	meta, err := storage.LoadTopicMeta(dir, parts)
	if err != nil {
		return fail(err)
	}
	journal, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	tasks, err := storage.CountTasksAtTopic(home, parts)
	if err != nil {
		return fail(err)
	}
	subs, err := storage.ListTopicChildNames(dir)
	if err != nil {
		return fail(err)
	}
	if subs == nil {
		subs = []string{}
	}
	if meta.Aliases == nil {
		meta.Aliases = []string{}
	}

	if asJSON {
		out := topicInfoJSON{
			Path:      storage.JoinTopicPath(parts),
			Title:     meta.Title,
			Aliases:   meta.Aliases,
			Dir:       dir,
			Notes:     len(journal),
			Tasks:     tasks,
			Subtopics: subs,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(out)
	}

	fmt.Printf("path: %s\n", storage.JoinTopicPath(parts))
	fmt.Printf("title: %s\n", meta.Title)
	if len(meta.Aliases) == 0 {
		fmt.Println("aliases:")
	} else {
		fmt.Printf("aliases: %s\n", strings.Join(meta.Aliases, ", "))
	}
	fmt.Printf("dir: %s\n", dir)
	fmt.Printf("notes: %d\n", len(journal))
	fmt.Printf("tasks: %d\n", tasks)
	fmt.Printf("subtopics: %d\n", len(subs))
	for _, name := range subs {
		fmt.Printf("  %s\n", name)
	}
	return nil
}

func runTopicNote(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"note"}, args...))

	var labels []string
	remaining, err := lessflags.
		StringSlice("--label", &labels).
		Help("-h,--help", topicNoteHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) < 2 {
		return topicErr("tsk topic note: topic and text required")
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return topicErr("tsk topic note: %v", err)
		}
	}
	parts, dir, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	if err := storage.MigrateLegacyTopicNotes(dir, parts); err != nil {
		return fail(err)
	}
	existing, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	text := joinArgs(remaining[1:])
	note := storage.TopicNote{
		TS:   storage.NowTimestamp(len(existing) + 1),
		Text: text,
	}
	if len(labels) > 0 {
		note.Labels = labels
	}
	if err := storage.AppendTopicNote(dir, note); err != nil {
		return fail(err)
	}
	invk.setData(storage.EventData{
		Topic:  storage.JoinTopicPath(parts),
		Text:   text,
		Labels: labels,
	})
	return nil
}

func runTopicNotes(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"notes"}, args...))

	var asJSON bool
	var limit int
	var labels []string
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Int("--limit", &limit).
		StringSlice("--label", &labels).
		Help("-h,--help", topicNotesHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 1 {
		return topicErr("tsk topic notes: topic required")
	}
	if limit < 0 {
		return topicErr("tsk topic notes: --limit must be >= 0")
	}
	for _, l := range labels {
		if err := storage.ValidateLabel(l); err != nil {
			return topicErr("tsk topic notes: %v", err)
		}
	}
	parts, dir, err := requireTopicDir(home, remaining[0])
	if err != nil {
		return topicErr("%s", err.Error())
	}
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts)})
	if err := storage.MigrateLegacyTopicNotes(dir, parts); err != nil {
		return fail(err)
	}
	notes, err := storage.ReadTopicNotes(dir)
	if err != nil {
		return fail(err)
	}
	notes = storage.ApplyNoteLimit(storage.FilterNotes(notes, labels), limit)
	return printNotes(notes, asJSON, false)
}

func runTopicAlias(invk *invocation, args []string) error {
	invk.setCommand("topic", append([]string{"alias"}, args...))

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(topicAliasHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runTopicAliasAdd(invk, args[1:])
	default:
		return topicErr("tsk topic alias: unknown subcommand %q", args[0])
	}
}

func runTopicAliasAdd(invk *invocation, args []string) error {
	home := invk.home
	invk.setCommand("topic", append([]string{"alias", "add"}, args...))

	remaining, err := lessflags.
		Help("-h,--help", topicAliasAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	if len(remaining) != 2 {
		return topicErr("tsk topic alias add: topic and alias required")
	}
	ref := remaining[0]
	alias := strings.Trim(strings.TrimSpace(remaining[1]), "/")
	if alias == "" {
		return topicErr("tsk topic alias add: alias required")
	}

	parts, dir, err := requireTopicDir(home, ref)
	if err != nil {
		return topicErr("%s", err.Error())
	}
	invk.setData(storage.EventData{Topic: storage.JoinTopicPath(parts), Alias: alias})

	if other, err := storage.FindAliasOwner(home, alias); err != nil {
		return topicErr("%s", err.Error())
	} else if other != nil && !topicPartsEqual(other, parts) {
		return topicErr("alias %s already used by %s", alias, storage.JoinTopicPath(other))
	}

	aliasParts := storage.SplitTopicPath(alias)
	if storage.TopicDirExists(home, aliasParts) && !topicPartsEqual(aliasParts, parts) {
		return topicErr("alias %s conflicts with topic path %s", alias, alias)
	}

	meta, err := storage.LoadTopicMeta(dir, parts)
	if err != nil {
		return fail(err)
	}
	meta.Aliases = append(meta.Aliases, alias)
	meta.UpdatedAt = storage.NowTimestamp(0)
	if err := storage.WriteTopicMeta(dir, meta); err != nil {
		return fail(err)
	}
	return nil
}
