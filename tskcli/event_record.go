package tskcli

import (
	"os"
	"strconv"
	"strings"

	"github.com/xhd2015/tsk/tskcli/storage"
)

// nestedCommands maps dotted action prefixes to valid next path segments.
var nestedCommands = map[string][]string{
	"label":               {"add", "rm", "list"},
	"note":                {"add", "list", "edit"},
	"progress":            {"add", "list", "edit", "archive", "show"},
	"clarify":             {"add", "list", "confirm"},
	"topic":               {"set", "mkdir", "rm", "where", "info", "note", "notes", "view", "alias"},
	"topic.alias":         {"add"},
	"channel":             {"create", "list", "archive", "delete", "send", "messages", "participants", "participant"},
	"channel.participant": {"add", "remove"},
	"project":             {"add", "tree", "list", "which", "register", "unregister", "notes"},
	"project.notes":       {"add", "list", "edit", "delete"},
}

var mutationActions = map[string]struct{}{
	"add":                        {},
	"advance":                    {},
	"stage":                      {},
	"done":                       {},
	"archive":                    {},
	"delete":                     {},
	"update":                     {},
	"followup":                   {},
	"install":                    {},
	"label.add":                  {},
	"label.rm":                   {},
	"note.add":                   {},
	"note.edit":                  {},
	"progress.add":               {},
	"progress.edit":              {},
	"progress.archive":           {},
	"clarify.add":                {},
	"clarify.confirm":            {},
	"topic.set":                  {},
	"topic.mkdir":                {},
	"topic.rm":                   {},
	"topic.note":                 {},
	"topic.alias.add":            {},
	"channel.create":             {},
	"channel.archive":            {},
	"channel.delete":             {},
	"channel.send":               {},
	"channel.participant.add":    {},
	"channel.participant.remove": {},
	"project.add":                {},
	"project.register":           {},
	"project.unregister":         {},
	"project.notes.add":          {},
	"project.notes.edit":         {},
	"project.notes.delete":       {},
	"skill.install":              {},
}

func (invk *invocation) setCommand(command string, eventArgs []string) {
	if invk == nil {
		return
	}
	invk.command = command
	if eventArgs == nil {
		eventArgs = []string{}
	}
	invk.eventArgs = eventArgs
	invk.action = deriveAction(command, eventArgs)
	invk.mutation = isMutationAction(invk.action)
}

func (invk *invocation) setMutation(m bool) {
	if invk == nil {
		return
	}
	invk.mutation = m
}

func (invk *invocation) setData(d storage.EventData) {
	if invk == nil {
		return
	}
	invk.data = mergeEventData(invk.data, d)
}

func deriveAction(command string, args []string) string {
	if command == "" {
		return ""
	}
	if command == "skill" {
		return deriveSkillAction(args)
	}
	parts := []string{command}
	cur := command
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			break
		}
		if !containsString(nestedCommands[cur], a) {
			break
		}
		parts = append(parts, a)
		cur = strings.Join(parts, ".")
	}
	return strings.Join(parts, ".")
}

func deriveSkillAction(args []string) string {
	for _, a := range args {
		switch a {
		case "--install":
			return "skill.install"
		case "--show":
			return "skill.show"
		case "--list", "-l":
			return "skill.list"
		}
	}
	return "skill"
}

func isMutationAction(action string) bool {
	_, ok := mutationActions[action]
	return ok
}

func eventIsMutation(ev storage.Event) bool {
	if ev.Action != "" {
		return ev.Mutation
	}
	return isMutationAction(deriveAction(ev.Command, ev.Args))
}

func eventUser() string {
	if v := os.Getenv("TSK_USER"); v != "" {
		return v
	}
	return os.Getenv("USER")
}

func mergeEventData(dst *storage.EventData, src storage.EventData) *storage.EventData {
	if dst == nil {
		dst = &storage.EventData{}
	}
	if src.TaskID != 0 {
		dst.TaskID = src.TaskID
	}
	if src.ParentID != 0 {
		dst.ParentID = src.ParentID
	}
	if src.Topic != "" {
		dst.Topic = src.Topic
	}
	if src.ChannelID != "" {
		dst.ChannelID = src.ChannelID
	}
	if src.Project != "" {
		dst.Project = src.Project
	}
	if src.Title != "" {
		dst.Title = src.Title
	}
	if src.Text != "" {
		dst.Text = src.Text
	}
	if src.Labels != nil {
		dst.Labels = src.Labels
	}
	if src.Label != "" {
		dst.Label = src.Label
	}
	if src.Stage != "" {
		dst.Stage = src.Stage
	}
	if src.Status != "" {
		dst.Status = src.Status
	}
	if src.Index != 0 {
		dst.Index = src.Index
	}
	if src.Name != "" {
		dst.Name = src.Name
	}
	if src.Handle != "" {
		dst.Handle = src.Handle
	}
	if src.Alias != "" {
		dst.Alias = src.Alias
	}
	if src.Query != "" {
		dst.Query = src.Query
	}
	if src.Notes != nil {
		dst.Notes = src.Notes
	}
	if src.MessageID != 0 {
		dst.MessageID = src.MessageID
	}
	return dst
}

func compactEventData(d *storage.EventData) *storage.EventData {
	if d == nil || eventDataEmpty(*d) {
		return nil
	}
	return d
}

func eventDataEmpty(d storage.EventData) bool {
	return d.TaskID == 0 &&
		d.ParentID == 0 &&
		d.Topic == "" &&
		d.ChannelID == "" &&
		d.Project == "" &&
		d.Title == "" &&
		d.Text == "" &&
		len(d.Labels) == 0 &&
		d.Label == "" &&
		d.Stage == "" &&
		d.Status == "" &&
		d.Index == 0 &&
		d.Name == "" &&
		d.Handle == "" &&
		d.Alias == "" &&
		d.Query == "" &&
		len(d.Notes) == 0 &&
		d.MessageID == 0
}

func formatLogLine(ev storage.Event) string {
	status := "ok"
	if ev.ExitCode != 0 {
		status = "fail"
	}
	action := ev.Action
	if action == "" {
		action = ev.Command
	}
	if action == "" {
		action = "(unknown)"
	}
	parts := []string{ev.TS, status, action}
	parts = append(parts, formatEventDataCompact(ev.Data)...)
	return strings.Join(parts, "  ")
}

func formatEventDataCompact(d *storage.EventData) []string {
	if d == nil {
		return nil
	}
	var out []string
	if d.TaskID != 0 {
		out = append(out, "task="+strconv.Itoa(d.TaskID))
	}
	if d.ParentID != 0 {
		out = append(out, "parent="+strconv.Itoa(d.ParentID))
	}
	if d.Topic != "" {
		out = append(out, "topic="+d.Topic)
	}
	if d.ChannelID != "" {
		out = append(out, "channel="+d.ChannelID)
	}
	if d.Project != "" {
		out = append(out, "project="+d.Project)
	}
	if d.Name != "" {
		out = append(out, "name="+d.Name)
	}
	if d.Handle != "" {
		out = append(out, "handle="+d.Handle)
	}
	if d.Alias != "" {
		out = append(out, "alias="+d.Alias)
	}
	if d.Index != 0 {
		out = append(out, "index="+strconv.Itoa(d.Index))
	}
	if d.Label != "" {
		out = append(out, "label="+d.Label)
	}
	if d.Stage != "" {
		out = append(out, "stage="+d.Stage)
	}
	if d.Status != "" {
		out = append(out, "status="+d.Status)
	}
	if d.MessageID != 0 {
		out = append(out, "msg="+strconv.Itoa(d.MessageID))
	}
	return out
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func containsString(ss []string, want string) bool {
	for _, s := range ss {
		if s == want {
			return true
		}
	}
	return false
}

func applyEventLimit(events []storage.Event, limit int) []storage.Event {
	if limit > 0 && len(events) > limit {
		return events[len(events)-limit:]
	}
	return events
}

func logsCountLabel(n int) string {
	if n == 1 {
		return "1 log"
	}
	return strconv.Itoa(n) + " logs"
}
