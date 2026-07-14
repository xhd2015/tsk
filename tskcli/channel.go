package tskcli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
	"github.com/xhd2015/tsk/tskcli/storage"
)

const (
	ansiGreen = "\x1b[32m"
	ansiGray  = "\x1b[90m"
	ansiReset = "\x1b[0m"
)

func runChannel(home string, args []string) error {
	setCommand(currentCtx, "channel", args)

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(channelHelp())
		return nil
	}
	switch args[0] {
	case "create":
		return runChannelCreate(home, args[1:])
	case "list":
		return runChannelList(home, args[1:])
	case "archive":
		return runChannelArchive(home, args[1:])
	case "delete":
		return runChannelDelete(home, args[1:])
	case "send":
		return runChannelSend(home, args[1:])
	case "messages":
		return runChannelMessages(home, args[1:])
	case "participants":
		return runChannelParticipants(home, args[1:])
	case "participant":
		return runChannelParticipant(home, args[1:])
	default:
		return channelFail(fmt.Errorf("tsk channel: unknown subcommand %q", args[0]))
	}
}

func channelFail(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.HasPrefix(msg, "Error:") {
		return err
	}
	return fmt.Errorf("Error: %s", msg)
}

func channelSuccess(word string, tty bool) string {
	if tty {
		return ansiGreen + word + ansiReset
	}
	return word
}

func runChannelCreate(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"create"}, args...))

	var channelID, userHandle string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Help("-h,--help", channelCreateHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 1 {
		return channelFail(fmt.Errorf("tsk channel create: name required"))
	}
	name := remaining[0]
	if name == "" {
		return channelFail(fmt.Errorf("tsk channel create: name required"))
	}

	if err := storage.EnsureLayout(home); err != nil {
		return err
	}

	if channelID == "" {
		channelID = storage.Slugify(name)
	}
	channelID = strings.ToLower(strings.TrimSpace(channelID))
	if err := storage.ValidateChannelID(channelID); err != nil {
		return channelFail(err)
	}

	creator, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}
	if err := storage.CreateChannel(home, name, channelID, creator); err != nil {
		return channelFail(err)
	}
	fmt.Println(channelID)
	return nil
}

func runChannelList(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"list"}, args...))

	var all, asJSON bool
	remaining, err := lessflags.
		Bool("--all", &all).
		Bool("--json", &asJSON).
		Help("-h,--help", channelListHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 0 {
		return channelFail(fmt.Errorf("tsk channel list: unexpected arguments"))
	}

	channels, err := storage.ListChannels(home, all)
	if err != nil {
		return channelFail(err)
	}

	if asJSON {
		data, err := json.Marshal(channels)
		if err != nil {
			return channelFail(err)
		}
		fmt.Println(string(data))
		return nil
	}

	tty := isStdoutTTY()
	if len(channels) == 0 {
		if tty {
			fmt.Printf("%s0 channels%s\n", ansiGray, ansiReset)
		}
		return nil
	}

	fmt.Printf("%-20s %-24s %s\n", "ID", "Name", "Status")
	for _, ch := range channels {
		status := ch.Status
		if tty {
			status = channelSuccess(status, true)
		}
		fmt.Printf("%-20s %-24s %s\n", ch.ID, ch.Name, status)
	}
	footer := fmt.Sprintf("%d channel(s)", len(channels))
	if tty {
		footer = ansiGray + footer + ansiReset
	}
	fmt.Println(footer)
	return nil
}

func runChannelArchive(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"archive"}, args...))

	var channelID string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		Help("-h,--help", channelArchiveHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 0 {
		return channelFail(fmt.Errorf("tsk channel archive: unexpected arguments"))
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel archive: --channel-id required"))
	}

	if err := storage.ArchiveChannel(home, channelID); err != nil {
		return channelFail(err)
	}
	fmt.Printf("archived %s\n", channelSuccess(channelID, isStdoutTTY()))
	return nil
}

func runChannelDelete(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"delete"}, args...))

	var channelID string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		Help("-h,--help", channelDeleteHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 0 {
		return channelFail(fmt.Errorf("tsk channel delete: unexpected arguments"))
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel delete: --channel-id required"))
	}

	if err := storage.DeleteChannel(home, channelID); err != nil {
		return channelFail(err)
	}
	fmt.Printf("deleted %s\n", channelSuccess(channelID, isStdoutTTY()))
	return nil
}

func runChannelSend(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"send"}, args...))

	var channelID, userHandle string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Help("-h,--help", channelSendHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel send: --channel-id required"))
	}
	if len(remaining) == 0 {
		return channelFail(fmt.Errorf("tsk channel send: message required"))
	}
	body := strings.Join(remaining, " ")

	sender, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}
	msgID, err := storage.SendChannelMessage(home, channelID, sender, body)
	if err != nil {
		return channelFail(err)
	}
	fmt.Printf("sent message %s\n", channelSuccess(fmt.Sprintf("%d", msgID), isStdoutTTY()))
	return nil
}

func runChannelMessages(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"messages"}, args...))

	var channelID, userHandle string
	var limit int
	var asJSON bool
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Int("--limit", &limit).
		Bool("--json", &asJSON).
		Help("-h,--help", channelMessagesHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 0 {
		return channelFail(fmt.Errorf("tsk channel messages: unexpected arguments"))
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel messages: --channel-id required"))
	}

	actor, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}
	ch, dir, _, err := storage.LoadChannelByID(home, channelID)
	if err != nil {
		return channelFail(err)
	}
	if err := storage.RequireChannelParticipant(ch, actor); err != nil {
		return channelFail(err)
	}

	msgs, err := storage.ReadChannelMessages(dir)
	if err != nil {
		return channelFail(err)
	}
	if limit > 0 && len(msgs) > limit {
		msgs = msgs[len(msgs)-limit:]
	}

	if asJSON {
		data, err := json.Marshal(msgs)
		if err != nil {
			return channelFail(err)
		}
		fmt.Println(string(data))
		return nil
	}

	var b strings.Builder
	for _, m := range msgs {
		fmt.Fprintf(&b, "[%d] %s  %s\n%s\n\n", m.ID, m.Sender, m.CreatedAt, m.Body)
	}
	if b.Len() > 0 {
		fmt.Print(b.String())
	}
	return nil
}

func runChannelParticipants(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"participants"}, args...))

	var channelID, userHandle string
	var asJSON bool
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Bool("--json", &asJSON).
		Help("-h,--help", channelParticipantsHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if len(remaining) != 0 {
		return channelFail(fmt.Errorf("tsk channel participants: unexpected arguments"))
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel participants: --channel-id required"))
	}

	actor, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}
	ch, _, _, err := storage.LoadChannelByID(home, channelID)
	if err != nil {
		return channelFail(err)
	}
	if err := storage.RequireChannelParticipant(ch, actor); err != nil {
		return channelFail(err)
	}

	if asJSON {
		data, err := json.Marshal(ch.Participants)
		if err != nil {
			return channelFail(err)
		}
		fmt.Println(string(data))
		return nil
	}

	for _, p := range ch.Participants {
		fmt.Printf("%s  %s\n", p.Handle, p.JoinedAt)
	}
	return nil
}

func runChannelParticipant(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"participant"}, args...))

	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(channelParticipantHelp())
		return nil
	}
	switch args[0] {
	case "add":
		return runChannelParticipantAdd(home, args[1:])
	case "remove":
		return runChannelParticipantRemove(home, args[1:])
	default:
		return channelFail(fmt.Errorf("tsk channel participant: unknown subcommand %q", args[0]))
	}
}

func runChannelParticipantAdd(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"participant", "add"}, args...))

	var channelID, userHandle string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Help("-h,--help", channelParticipantAddHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel participant add: --channel-id required"))
	}
	if len(remaining) != 1 {
		return channelFail(fmt.Errorf("tsk channel participant add: handle required"))
	}
	handle := remaining[0]

	actor, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}
	added, err := storage.AddChannelParticipant(home, channelID, actor, handle)
	if err != nil {
		return channelFail(err)
	}
	handle, err = storage.NormalizeHandle(handle)
	if err != nil {
		return channelFail(err)
	}
	_ = added
	fmt.Printf("added %s\n", channelSuccess(handle, isStdoutTTY()))
	return nil
}

func runChannelParticipantRemove(home string, args []string) error {
	setCommand(currentCtx, "channel", append([]string{"participant", "remove"}, args...))

	var channelID, userHandle string
	remaining, err := lessflags.
		String("--channel-id", &channelID).
		String("--user", &userHandle).
		Help("-h,--help", channelParticipantRemoveHelp()).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return channelFail(err)
	}
	if channelID == "" {
		return channelFail(fmt.Errorf("tsk channel participant remove: --channel-id required"))
	}

	actor, err := storage.ResolveChannelIdentity(userHandle)
	if err != nil {
		return channelFail(err)
	}

	var target string
	selfLeave := false
	if len(remaining) == 0 {
		target = actor
		selfLeave = true
	} else if len(remaining) == 1 {
		target = remaining[0]
		selfLeave = target == actor
	} else {
		return channelFail(fmt.Errorf("tsk channel participant remove: unexpected arguments"))
	}

	if err := storage.RemoveChannelParticipant(home, channelID, actor, target); err != nil {
		return channelFail(err)
	}

	tty := isStdoutTTY()
	if selfLeave {
		fmt.Printf("left %s\n", channelSuccess(channelID, tty))
	} else {
		target, err = storage.NormalizeHandle(target)
		if err != nil {
			return channelFail(err)
		}
		fmt.Printf("removed %s\n", channelSuccess(target, tty))
	}
	return nil
}

func channelHelp() string {
	return `Usage: tsk channel <command> [arguments]

Subcommands:
  create        create a new channel
  list          list channels
  archive       archive a channel (readonly)
  delete        delete a channel
  send          send a message
  messages      show channel messages
  participants  list channel participants
  participant   add or remove participants

  -h, --help    show this help
`
}

func channelCreateHelp() string {
	return `Usage: tsk channel create [--channel-id ID] [--user HANDLE] <name>

Create a new channel. Creator and agent are auto-added as participants.

Flags:
  --channel-id ID   channel id (default: slugified name)
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}

func channelListHelp() string {
	return `Usage: tsk channel list [--json] [--all]

List channels (active only by default).

Flags:
  --json        output JSON array
  --all         include archived channels
  -h, --help    show this help
`
}

func channelArchiveHelp() string {
	return `Usage: tsk channel archive --channel-id ID

Archive a channel (readonly for mutations).

Flags:
  --channel-id ID   channel id
  -h, --help        show this help
`
}

func channelDeleteHelp() string {
	return `Usage: tsk channel delete --channel-id ID

Delete a channel and write a tombstone.

Flags:
  --channel-id ID   channel id
  -h, --help        show this help
`
}

func channelSendHelp() string {
	return `Usage: tsk channel send --channel-id ID [--user HANDLE] <message...>

Send a message to a channel.

Flags:
  --channel-id ID   channel id
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}

func channelMessagesHelp() string {
	return `Usage: tsk channel messages --channel-id ID [--json] [--limit N] [--user HANDLE]

Show channel messages (participants only).

Flags:
  --channel-id ID   channel id
  --json            output JSON array
  --limit N         show last N messages
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}

func channelParticipantsHelp() string {
	return `Usage: tsk channel participants --channel-id ID [--json] [--user HANDLE]

List channel participants (participants only).

Flags:
  --channel-id ID   channel id
  --json            output JSON array
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}

func channelParticipantHelp() string {
	return `Usage: tsk channel participant <command> [arguments]

Subcommands:
  add --channel-id ID <handle>       add participant
  remove --channel-id ID [<handle>]  remove participant or leave channel

  -h, --help                         show this help
`
}

func channelParticipantAddHelp() string {
	return `Usage: tsk channel participant add --channel-id ID [--user HANDLE] <handle>

Add a participant to a channel.

Flags:
  --channel-id ID   channel id
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}

func channelParticipantRemoveHelp() string {
	return `Usage: tsk channel participant remove --channel-id ID [--user HANDLE] [<handle>]

Remove a participant or leave the channel (omit handle).

Flags:
  --channel-id ID   channel id
  --user HANDLE     acting user (default: TSK_USER or $USER)
  -h, --help        show this help
`
}