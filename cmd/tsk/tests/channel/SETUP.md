# Scenario

**Feature**: `tsk channel` manages Slack-like conversational channels under `TSK_HOME/channels/`

```
# identity: --user > TSK_USER > $USER; creator + agent auto-join on create
TSK_USER=alice -> tsk channel create <name> -> active/<id>/channel.json + index + messages.jsonl
tsk channel send|messages|participants|participant * -> membership-gated reads/writes
tsk channel archive|delete -> lifecycle + tombstone blocks id reuse
```

## Preconditions

- Fresh `TSK_HOME` per leaf unless Setup seeds channels.
- Default participant identity `TSK_USER=alice` via `req.ExtraEnv` unless a leaf overrides.
- Channel storage: `channels/index/<id>`, `channels/active|archive/<id>/`, `channels/tombstones/<id>`.

## Context

Helpers mirror task helpers: `runTskCmd`/`runTskOK` with `tskEnv(req)`. Channel-specific readers assert on-disk layout and JSONL transcripts.

```go
import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

type channelParticipant struct {
	Handle   string `json:"handle"`
	JoinedAt string `json:"joined_at"`
}

type channelJSON struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Status       string               `json:"status"`
	Participants []channelParticipant `json:"participants"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
}

type channelMessage struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

type channelTombstone struct {
	ID        string `json:"id"`
	DeletedAt string `json:"deleted_at"`
}

func Setup(t *testing.T, req *Request) error {
	ensureChannelHelpersUsed()
	hasUser := false
	for _, e := range req.ExtraEnv {
		if strings.HasPrefix(e, "TSK_USER=") {
			hasUser = true
			break
		}
	}
	if !hasUser {
		req.ExtraEnv = append(req.ExtraEnv, "TSK_USER=alice")
	}
	return nil
}

func channelsRoot(req *Request) string {
	return filepath.Join(req.TskHome, "channels")
}

func channelAbs(req *Request, rel string) string {
	return filepath.Join(channelsRoot(req), filepath.FromSlash(rel))
}

func activeChannelDir(req *Request, id string) string {
	return channelAbs(req, filepath.Join("active", id))
}

func archiveChannelDir(req *Request, id string) string {
	return channelAbs(req, filepath.Join("archive", id))
}

func readChannelIndex(t *testing.T, req *Request, id string) string {
	t.Helper()
	data, err := os.ReadFile(channelAbs(req, filepath.Join("index", id)))
	if err != nil {
		t.Fatalf("read channels/index/%s: %v", id, err)
	}
	return strings.TrimSpace(string(data))
}

func readChannelJSON(t *testing.T, channelDir string) channelJSON {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(channelDir, "channel.json"))
	if err != nil {
		t.Fatalf("read %s/channel.json: %v", channelDir, err)
	}
	var ch channelJSON
	if err := json.Unmarshal(data, &ch); err != nil {
		t.Fatalf("parse channel.json: %v", err)
	}
	return ch
}

func readMessagesJSONL(t *testing.T, channelDir string) []channelMessage {
	t.Helper()
	path := filepath.Join(channelDir, "messages.jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	var msgs []channelMessage
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m channelMessage
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			t.Fatalf("parse messages.jsonl line: %v", err)
		}
		msgs = append(msgs, m)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan messages.jsonl: %v", err)
	}
	return msgs
}

func readMsgCounter(t *testing.T, channelDir string) int {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(channelDir, "msg-counter"))
	if err != nil {
		t.Fatalf("read msg-counter: %v", err)
	}
	var n int
	if _, err := fmt.Sscanf(strings.TrimSpace(string(data)), "%d", &n); err != nil {
		t.Fatalf("parse msg-counter: %v", err)
	}
	return n
}

func readTombstone(t *testing.T, req *Request, id string) channelTombstone {
	t.Helper()
	data, err := os.ReadFile(channelAbs(req, filepath.Join("tombstones", id)))
	if err != nil {
		t.Fatalf("read tombstone %s: %v", id, err)
	}
	var ts channelTombstone
	if err := json.Unmarshal(data, &ts); err != nil {
		t.Fatalf("parse tombstone: %v", err)
	}
	return ts
}

func createChannelArgs(name string, channelID string) []string {
	args := []string{"channel", "create", name}
	if channelID != "" {
		args = append(args, "--channel-id", channelID)
	}
	return args
}

func createChannel(t *testing.T, req *Request, name string, channelID string) string {
	t.Helper()
	resp := runTskOK(t, req, createChannelArgs(name, channelID)...)
	id := strings.TrimSpace(resp.Stdout)
	want := channelID
	if want == "" {
		want = slugify(name)
	}
	if id != want {
		t.Fatalf("create channel id: got %q want %q", id, want)
	}
	req.ChannelID = id
	req.ChannelName = name
	return id
}

func addParticipant(t *testing.T, req *Request, channelID, handle string) {
	t.Helper()
	runTskOK(t, req, "channel", "participant", "add", "--channel-id", channelID, handle)
}

func sendChannelMessage(t *testing.T, req *Request, channelID string, message string, extraArgs ...string) *Response {
	t.Helper()
	args := []string{"channel", "send", "--channel-id", channelID}
	args = append(args, extraArgs...)
	args = append(args, message)
	return runTskOK(t, req, args...)
}

func archiveChannel(t *testing.T, req *Request, channelID string) {
	t.Helper()
	runTskOK(t, req, "channel", "archive", "--channel-id", channelID)
}

func deleteChannel(t *testing.T, req *Request, channelID string) {
	t.Helper()
	runTskOK(t, req, "channel", "delete", "--channel-id", channelID)
}

func participantHandles(ch channelJSON) []string {
	out := make([]string, len(ch.Participants))
	for i, p := range ch.Participants {
		out[i] = p.Handle
	}
	return out
}

func assertChannelParticipantsSorted(t *testing.T, ch channelJSON, want []string) {
	t.Helper()
	got := participantHandles(ch)
	if len(got) != len(want) {
		t.Fatalf("participants: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("participants: got %v want %v", got, want)
		}
	}
}

func assertChannelIndexEquals(t *testing.T, req *Request, id, want string) {
	t.Helper()
	got := readChannelIndex(t, req, id)
	if got != want {
		t.Fatalf("channels/index/%s: got %q want %q", id, got, want)
	}
}

func assertStderrErrorPrefix(t *testing.T, stderr string) {
	t.Helper()
	assertStderrContainsCount(t, stderr, "Error:", 1)
}

func assertNoANSI(t *testing.T, s string) {
	t.Helper()
	if strings.Contains(s, "\x1b[") {
		t.Fatalf("expected no ANSI in output, got %q", s)
	}
}

func assertChannelEventCommand(t *testing.T, req *Request) {
	t.Helper()
	assertLastEventCommand(t, req, "channel")
}

func withTSKUser(t *testing.T, req *Request, user string) {
	t.Helper()
	out := make([]string, 0, len(req.ExtraEnv)+1)
	replaced := false
	for _, e := range req.ExtraEnv {
		if strings.HasPrefix(e, "TSK_USER=") {
			out = append(out, "TSK_USER="+user)
			replaced = true
			continue
		}
		out = append(out, e)
	}
	if !replaced {
		out = append(out, "TSK_USER="+user)
	}
	req.ExtraEnv = out
}

func ensureChannelHelpersUsed() {
	_ = channelsRoot
	_ = channelAbs
	_ = activeChannelDir
	_ = archiveChannelDir
	_ = readChannelIndex
	_ = readChannelJSON
	_ = readMessagesJSONL
	_ = readMsgCounter
	_ = readTombstone
	_ = createChannelArgs
	_ = createChannel
	_ = addParticipant
	_ = sendChannelMessage
	_ = archiveChannel
	_ = deleteChannel
	_ = participantHandles
	_ = assertChannelParticipantsSorted
	_ = assertChannelIndexEquals
	_ = assertStderrErrorPrefix
	_ = assertNoANSI
	_ = assertChannelEventCommand
	_ = withTSKUser
	_ = sort.Strings
}
```