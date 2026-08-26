package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/terminal/color"
	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

type searchKindSet struct {
	task     bool
	note     bool
	progress bool
	topic    bool
}

func resolveSearchKinds(flags lessflags.Flags) searchKindSet {
	var set searchKindSet
	var sawAll bool
	for _, f := range flags {
		switch f.Flag {
		case "--all":
			sawAll = true
		case "--task":
			set.task = true
		case "--note":
			set.note = true
		case "--progress":
			set.progress = true
		case "--topic":
			set.topic = true
		}
	}
	if sawAll || (!set.task && !set.note && !set.progress && !set.topic) {
		return searchKindSet{task: true, note: true, progress: true, topic: true}
	}
	return set
}

// searchHit is one match under TSK_HOME.
type searchHit struct {
	Kind   string   `json:"kind"` // task | note | progress | topic
	TaskID int      `json:"task_id,omitempty"`
	Topic  string   `json:"topic,omitempty"`
	Text   string   `json:"text"`
	TS     string   `json:"ts,omitempty"`
	Labels []string `json:"labels,omitempty"`
	Status string   `json:"status,omitempty"`
}

func runSearch(home string, args []string) error {
	setCommand(currentCtx, "search", args)

	var kinds lessflags.Flags
	var asJSON bool
	var colorFlag, noColorFlag bool
	remaining, err := lessflags.
		Bool("--json", &asJSON).
		Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
		Group(
			lessflags.CollectParsedFlags(&kinds).
				Bool("--task").
				Bool("--note").
				Bool("--progress").
				Bool("--topic").
				Bool("--all"),
		).
		Help("-h,--help", searchHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return fail(err)
	}
	mode, err := color.ModeFromFlags(colorFlag, noColorFlag)
	if err != nil {
		return searchErr("%s", err.Error())
	}
	query := strings.TrimSpace(joinArgs(remaining))
	if query == "" {
		return searchErr("tsk search: query required")
	}

	want := resolveSearchKinds(kinds)
	hits, err := collectSearchHits(home, query, want)
	if err != nil {
		return fail(err)
	}
	style := color.Style{Enabled: !asJSON && color.EnabledFor(mode, os.Stdout)}
	return printSearchHits(hits, query, asJSON, style)
}

func searchErr(format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !strings.HasPrefix(msg, "Error:") {
		msg = "Error: " + msg
	}
	return fail(fmt.Errorf("%s", msg))
}

func collectSearchHits(home, query string, want searchKindSet) ([]searchHit, error) {
	var hits []searchHit
	q := strings.ToLower(query)

	if want.task || want.note || want.progress {
		ids, err := storage.ListTaskIDs(home)
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			task, dir, err := storage.LoadTaskByID(home, id)
			if err != nil {
				return nil, err
			}
			topicStr, err := topicDisplay(task)
			if err != nil {
				return nil, err
			}
			if want.task && strings.Contains(strings.ToLower(task.Title), q) {
				hits = append(hits, searchHit{
					Kind:   "task",
					TaskID: id,
					Topic:  topicStr,
					Text:   task.Title,
				})
			}
			if want.note || want.progress {
				notes, err := storage.ReadTopicNotes(dir)
				if err != nil {
					return nil, err
				}
				for _, n := range notes {
					isProgress := storage.NoteHasAllLabels(n, []string{"progress"})
					if isProgress && !want.progress {
						continue
					}
					if !isProgress && !want.note {
						continue
					}
					if !strings.Contains(strings.ToLower(n.Text), q) {
						continue
					}
					kind := "note"
					if isProgress {
						kind = "progress"
					}
					hit := searchHit{
						Kind:   kind,
						TaskID: id,
						Topic:  topicStr,
						Text:   n.Text,
						TS:     n.TS,
						Status: n.Status,
					}
					if len(n.Labels) > 0 {
						hit.Labels = n.Labels
					}
					hits = append(hits, hit)
				}
			}
		}
	}

	if want.topic {
		topicHits, err := searchTopicNotes(home, q)
		if err != nil {
			return nil, err
		}
		hits = append(hits, topicHits...)
	}
	return hits, nil
}

func topicDisplay(task storage.Task) (string, error) {
	parts, err := storage.ParseTopicPath(task.TopicPath)
	if err != nil {
		return "", err
	}
	if len(parts) == 0 {
		return "inbox", nil
	}
	return storage.JoinTopicPath(parts), nil
}

func searchTopicNotes(home, qLower string) ([]searchHit, error) {
	root := filepath.Join(home, "topics")
	var hits []searchHit
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) && path == root {
				return nil
			}
			return err
		}
		if d.IsDir() {
			if path == root {
				return nil
			}
			if storage.IsTaskDirName(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Name() != storage.TopicNotesFile {
			return nil
		}
		topicDir := filepath.Dir(path)
		rel, err := filepath.Rel(root, topicDir)
		if err != nil {
			return err
		}
		topicPath := filepath.ToSlash(rel)
		notes, err := storage.ReadTopicNotes(topicDir)
		if err != nil {
			return err
		}
		for _, n := range notes {
			if !strings.Contains(strings.ToLower(n.Text), qLower) {
				continue
			}
			hit := searchHit{
				Kind:  "topic",
				Topic: topicPath,
				Text:  n.Text,
				TS:    n.TS,
			}
			if len(n.Labels) > 0 {
				hit.Labels = n.Labels
			}
			hits = append(hits, hit)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return hits, nil
}

func printSearchHits(hits []searchHit, query string, asJSON bool, style color.Style) error {
	if asJSON {
		if hits == nil {
			hits = []searchHit{}
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetEscapeHTML(false)
		return enc.Encode(hits)
	}
	for _, h := range hits {
		fmt.Println(formatSearchHitHeader(h, style))
		fmt.Printf("  %s\n", highlightQuery(style, h.Text, query))
	}
	n := len(hits)
	var footer string
	if n == 1 {
		footer = "1 match"
	} else {
		footer = fmt.Sprintf("%d matches", n)
	}
	fmt.Println(style.Gray(footer))
	return nil
}

func formatSearchHitHeader(h searchHit, style color.Style) string {
	switch h.Kind {
	case "task", "note", "progress":
		loc := h.Topic
		if loc == "" {
			loc = "inbox"
		}
		anchor := style.Green(fmt.Sprintf("task %d", h.TaskID))
		kind := style.Blue(h.Kind)
		meta := style.Gray(loc)
		extra := ""
		if len(h.Labels) > 0 {
			extra += "  " + style.Gray("labels="+strings.Join(h.Labels, ","))
		}
		if h.Status != "" {
			extra += "  " + style.Gray("status="+h.Status)
		}
		return fmt.Sprintf("%s  %s  %s%s", anchor, kind, meta, extra)
	case "topic":
		anchor := style.Green("topic")
		kind := style.Blue("note")
		meta := style.Gray(h.Topic)
		extra := ""
		if len(h.Labels) > 0 {
			extra += "  " + style.Gray("labels="+strings.Join(h.Labels, ","))
		}
		return fmt.Sprintf("%s  %s  %s%s", anchor, kind, meta, extra)
	default:
		return h.Kind
	}
}

// highlightQuery bolds the first case-insensitive query match when style is on.
func highlightQuery(style color.Style, text, query string) string {
	if !style.Enabled || query == "" {
		return text
	}
	runes := []rune(text)
	qRunes := []rune(query)
	n := len(qRunes)
	if n == 0 || n > len(runes) {
		return text
	}
	for i := 0; i+n <= len(runes); i++ {
		if strings.EqualFold(string(runes[i:i+n]), query) {
			before := string(runes[:i])
			mid := string(runes[i : i+n])
			after := string(runes[i+n:])
			return before + style.Bold(mid) + after
		}
	}
	return text
}
