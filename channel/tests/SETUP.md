# Scenario

**Feature**: isolated TSK_HOME per test; direct FileStore calls

```
# temp WorkRoot + {WorkRoot}/.tsk per leaf
TSK_DATE=2026-07-09 -> file.FileStore -> channel.Store methods -> on-disk layout
```

## Preconditions

- Module root is `github.com/xhd2015/tsk` (parent of `channel/`).
- Go toolchain on PATH.
- Per-leaf `TSK_HOME` at `{WorkRoot}/.tsk`; never shared across leaves.
- Default creator handle `alice` unless a leaf overrides `req.Creator`.

## Context

Helpers read normalized on-disk layout: metadata-only `channel.json`, separate
`participants.jsonl`, `messages.jsonl`, `msg-counter`, `index/<id>`,
`tombstones/<id>.json`.

```go
import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/tsk/channel"
)

const tskDate = "2026-07-09"

type channelParticipant struct {
	Handle   string `json:"handle"`
	JoinedAt string `json:"joined_at"`
}

type channelMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
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
	workRoot, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		return fmt.Errorf("resolve work root: %w", err)
	}
	req.WorkRoot = workRoot
	req.TskHome = filepath.Join(req.WorkRoot, ".tsk")
	if req.Creator == "" {
		req.Creator = "alice"
	}
	ensureStoreHelpersUsed()
	return os.MkdirAll(req.TskHome, 0o755)
}

func storeEnv(req *Request) []string {
	env := make([]string, 0, len(os.Environ())+2)
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "TSK_HOME=") || strings.HasPrefix(e, "TSK_DATE=") {
			continue
		}
		env = append(env, e)
	}
	env = append(env, "TSK_HOME="+req.TskHome, "TSK_DATE="+tskDate)
	return env
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

func slugify(name string) string {
	s := strings.ToLower(name)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('-')
		}
	}
	out := b.String()
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	out = strings.Trim(out, "-")
	if len(out) > 64 {
		out = out[:64]
		out = strings.Trim(out, "-")
	}
	return out
}

func readChannelIndex(t *testing.T, req *Request, id string) string {
	t.Helper()
	data, err := os.ReadFile(channelAbs(req, filepath.Join("index", id)))
	if err != nil {
		t.Fatalf("read channels/index/%s: %v", id, err)
	}
	return strings.TrimSpace(string(data))
}

func readChannelMetadata(t *testing.T, channelDir string) channelMetadata {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(channelDir, "channel.json"))
	if err != nil {
		t.Fatalf("read %s/channel.json: %v", channelDir, err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("parse channel.json: %v", err)
	}
	if _, ok := raw["participants"]; ok {
		t.Fatalf("channel.json must not contain participants field; got keys %v", raw)
	}
	var ch channelMetadata
	if err := json.Unmarshal(data, &ch); err != nil {
		t.Fatalf("decode channel.json metadata: %v", err)
	}
	return ch
}

func readParticipantsJSONL(t *testing.T, channelDir string) []channelParticipant {
	t.Helper()
	path := filepath.Join(channelDir, "participants.jsonl")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	var out []channelParticipant
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var p channelParticipant
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			t.Fatalf("parse participants.jsonl line: %v", err)
		}
		out = append(out, p)
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan participants.jsonl: %v", err)
	}
	return out
}

func participantHandles(parts []channelParticipant) []string {
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = p.Handle
	}
	return out
}

func assertParticipantHandlesSorted(t *testing.T, channelDir string, want []string) {
	t.Helper()
	got := participantHandles(readParticipantsJSONL(t, channelDir))
	if len(got) != len(want) {
		t.Fatalf("participants: got %v want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("participants: got %v want %v", got, want)
		}
	}
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
	data, err := os.ReadFile(channelAbs(req, filepath.Join("tombstones", id+".json")))
	if err != nil {
		t.Fatalf("read tombstone %s.json: %v", id, err)
	}
	var ts channelTombstone
	if err := json.Unmarshal(data, &ts); err != nil {
		t.Fatalf("parse tombstone: %v", err)
	}
	return ts
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	if st, err := os.Stat(path); err != nil || !st.IsDir() {
		t.Fatalf("expected dir %s: %v", path, err)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s: %v", path, err)
	}
}

func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected missing %s, err=%v", path, err)
	}
}

func assertStoreErr(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("expected store error, got nil")
	}
}

func assertStoreOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected store error: %v", err)
	}
}

func seedChannel(t *testing.T, req *Request, name, id string) {
	t.Helper()
	prev := req.Op
	req.ChannelName = name
	if id != "" {
		req.ChannelID = id
	} else {
		req.ChannelID = slugify(name)
	}
	restore := pushStoreEnv(req)
	defer restore()
	store := newFileStore(t, req)
	_, err := store.Create(context.Background(), channel.CreateRequest{
		Name:    req.ChannelName,
		ID:      req.ChannelID,
		Creator: req.Creator,
	})
	if err != nil {
		t.Fatalf("seed create: %v", err)
	}
	req.Op = prev
}

func ensureStoreHelpersUsed() {
	_ = storeEnv
	_ = channelsRoot
	_ = channelAbs
	_ = activeChannelDir
	_ = archiveChannelDir
	_ = slugify
	_ = readChannelIndex
	_ = readChannelMetadata
	_ = readParticipantsJSONL
	_ = participantHandles
	_ = assertParticipantHandlesSorted
	_ = readMessagesJSONL
	_ = readMsgCounter
	_ = readTombstone
	_ = assertDirExists
	_ = assertFileExists
	_ = assertFileNotExists
	_ = assertStoreErr
	_ = assertStoreOK
	_ = seedChannel
}
```