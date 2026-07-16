package file

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/xhd2015/tsk/channel"
)

// Options configures a FileStore.
type Options struct {
	Home string
}

// FileStore implements channel.Store on disk under TSK_HOME/channels/.
type FileStore struct {
	home string
}

// NewStore opens a file-backed channel store.
func NewStore(opts Options) (channel.Store, error) {
	home := opts.Home
	if home == "" {
		var err error
		home, err = ResolveHome()
		if err != nil {
			return nil, err
		}
	}
	home, err := filepath.Abs(home)
	if err != nil {
		return nil, fmt.Errorf("resolve home: %w", err)
	}
	if err := ensureLayout(home); err != nil {
		return nil, err
	}
	return &FileStore{home: home}, nil
}

// ResolveHome returns TSK_HOME or ~/.tsk.
func ResolveHome() (string, error) {
	if v := os.Getenv("TSK_HOME"); v != "" {
		return filepath.Abs(v)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".tsk"), nil
}

func ensureLayout(home string) error {
	for _, dir := range []string{
		"channels/index", "channels/active", "channels/archive", "channels/tombstones",
	} {
		if err := os.MkdirAll(filepath.Join(home, dir), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", dir, err)
		}
	}
	return nil
}

func (s *FileStore) channelsRoot() string {
	return filepath.Join(s.home, "channels")
}

func channelIndexPath(home, id string) string {
	return filepath.Join(home, "channels", "index", id)
}

func channelTombstonePath(home, id string) string {
	return filepath.Join(home, "channels", "tombstones", id+".json")
}

func channelActiveDir(home, id string) string {
	return filepath.Join(home, "channels", "active", id)
}

func channelArchiveDir(home, id string) string {
	return filepath.Join(home, "channels", "archive", id)
}

func nowTimestamp(seq int) string {
	if date := os.Getenv("TSK_DATE"); date != "" {
		return fmt.Sprintf("%sT%02d:00:00Z", date, seq%24)
	}
	return time.Now().UTC().Format(time.RFC3339)
}

func indexStatus(home, id string) (string, error) {
	data, err := os.ReadFile(channelIndexPath(home, id))
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	line := strings.TrimSpace(string(data))
	switch {
	case strings.HasPrefix(line, "active/"):
		return "active", nil
	case strings.HasPrefix(line, "archive/"):
		return "archive", nil
	default:
		return "", fmt.Errorf("invalid index entry for channel %q: %q", id, line)
	}
}

func tombstoned(home, id string) (bool, error) {
	_, err := os.Stat(channelTombstonePath(home, id))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func idAvailable(home, id string) (bool, error) {
	status, err := indexStatus(home, id)
	if err != nil {
		return false, err
	}
	if status != "" {
		return false, nil
	}
	tomb, err := tombstoned(home, id)
	if err != nil {
		return false, err
	}
	return !tomb, nil
}

func writeIndex(home, id, rel string) error {
	dir := filepath.Join(home, "channels", "index")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create channel index dir: %w", err)
	}
	tmp, err := os.CreateTemp(dir, id+"-*.tmp")
	if err != nil {
		return fmt.Errorf("create channel index temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.WriteString(rel); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write channel index temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, channelIndexPath(home, id)); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename channel index: %w", err)
	}
	return nil
}

type channelLocation struct {
	dir    string
	status string
}

func (s *FileStore) locate(id string) (channelLocation, error) {
	status, err := indexStatus(s.home, id)
	if err != nil {
		return channelLocation{}, err
	}
	if status == "" {
		return channelLocation{}, fmt.Errorf("channel %q not found", id)
	}
	var dir string
	switch status {
	case "active":
		dir = channelActiveDir(s.home, id)
	case "archive":
		dir = channelArchiveDir(s.home, id)
	default:
		return channelLocation{}, fmt.Errorf("channel %q not found", id)
	}
	return channelLocation{dir: dir, status: status}, nil
}

func readChannelMetadata(channelDir string) (channel.Channel, error) {
	data, err := os.ReadFile(filepath.Join(channelDir, "channel.json"))
	if err != nil {
		return channel.Channel{}, fmt.Errorf("read channel.json: %w", err)
	}
	var ch channel.Channel
	if err := json.Unmarshal(data, &ch); err != nil {
		return channel.Channel{}, fmt.Errorf("parse channel.json: %w", err)
	}
	return ch, nil
}

func writeChannelMetadata(channelDir string, ch channel.Channel) error {
	data, err := json.MarshalIndent(ch, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(channelDir, "channel-*.json.tmp")
	if err != nil {
		return fmt.Errorf("create channel temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write channel temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	dst := filepath.Join(channelDir, "channel.json")
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename channel.json: %w", err)
	}
	return nil
}

func readParticipants(channelDir string) ([]channel.Participant, error) {
	path := filepath.Join(channelDir, "participants.jsonl")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var out []channel.Participant
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var p channel.Participant
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			return nil, fmt.Errorf("parse participants.jsonl: %w", err)
		}
		out = append(out, p)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func sortParticipants(parts []channel.Participant) {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Handle < parts[j].Handle
	})
}

func writeParticipants(channelDir string, parts []channel.Participant) error {
	sortParticipants(parts)
	var b strings.Builder
	for _, p := range parts {
		line, err := json.Marshal(p)
		if err != nil {
			return err
		}
		b.Write(line)
		b.WriteByte('\n')
	}
	tmp, err := os.CreateTemp(channelDir, "participants-*.jsonl.tmp")
	if err != nil {
		return fmt.Errorf("create participants temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.WriteString(b.String()); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write participants temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	dst := filepath.Join(channelDir, "participants.jsonl")
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename participants.jsonl: %w", err)
	}
	return nil
}

func isParticipant(parts []channel.Participant, handle string) bool {
	for _, p := range parts {
		if p.Handle == handle {
			return true
		}
	}
	return false
}

func requireParticipant(parts []channel.Participant, chID, handle string) error {
	if !isParticipant(parts, handle) {
		return fmt.Errorf("%q is not a participant in channel %q", handle, chID)
	}
	return nil
}

func (s *FileStore) Create(_ context.Context, req channel.CreateRequest) (*channel.Channel, error) {
	id := strings.ToLower(strings.TrimSpace(req.ID))
	if id == "" {
		id = channel.Slugify(req.Name)
	}
	if err := channel.ValidateID(id); err != nil {
		return nil, err
	}
	creator, err := channel.NormalizeHandle(req.Creator)
	if err != nil {
		return nil, err
	}
	avail, err := idAvailable(s.home, id)
	if err != nil {
		return nil, err
	}
	if !avail {
		return nil, fmt.Errorf("channel %q already exists", id)
	}

	dir := channelActiveDir(s.home, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create channel dir: %w", err)
	}

	now := nowTimestamp(channel.Seq(id))
	ch := channel.Channel{
		ID:        id,
		Name:      req.Name,
		Status:    "active",
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := writeChannelMetadata(dir, ch); err != nil {
		return nil, err
	}
	if err := writeParticipants(dir, []channel.Participant{
		{Handle: creator, JoinedAt: now},
	}); err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "messages.jsonl"), nil, 0o644); err != nil {
		return nil, fmt.Errorf("create messages.jsonl: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "msg-counter"), []byte("0"), 0o644); err != nil {
		return nil, fmt.Errorf("create msg-counter: %w", err)
	}
	if err := writeIndex(s.home, id, "active/"+id); err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *FileStore) List(_ context.Context, opts channel.ListOptions) ([]channel.ListEntry, error) {
	indexDir := filepath.Join(s.channelsRoot(), "index")
	entries, err := os.ReadDir(indexDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []channel.ListEntry
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		id := ent.Name()
		status, err := indexStatus(s.home, id)
		if err != nil || status == "" {
			continue
		}
		if !opts.All && status == "archive" {
			continue
		}
		loc, err := s.locate(id)
		if err != nil {
			continue
		}
		ch, err := readChannelMetadata(loc.dir)
		if err != nil {
			continue
		}
		out = append(out, channel.ListEntry{
			ID:     ch.ID,
			Name:   ch.Name,
			Status: ch.Status,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out, nil
}

func (s *FileStore) Get(_ context.Context, id string) (*channel.Channel, error) {
	loc, err := s.locate(id)
	if err != nil {
		return nil, err
	}
	ch, err := readChannelMetadata(loc.dir)
	if err != nil {
		return nil, err
	}
	return &ch, nil
}

func (s *FileStore) Archive(_ context.Context, id string) error {
	status, err := indexStatus(s.home, id)
	if err != nil {
		return err
	}
	if status == "" {
		return fmt.Errorf("channel %q not found", id)
	}
	if status == "archive" {
		return fmt.Errorf("channel %q is already archived", id)
	}
	activeDir := channelActiveDir(s.home, id)
	archiveDir := channelArchiveDir(s.home, id)
	if err := os.MkdirAll(filepath.Dir(archiveDir), 0o755); err != nil {
		return err
	}
	if err := os.Rename(activeDir, archiveDir); err != nil {
		return fmt.Errorf("archive channel: %w", err)
	}
	ch, err := readChannelMetadata(archiveDir)
	if err != nil {
		return err
	}
	ch.Status = "archived"
	ch.UpdatedAt = nowTimestamp(channel.Seq(id))
	if err := writeChannelMetadata(archiveDir, ch); err != nil {
		return err
	}
	return writeIndex(s.home, id, "archive/"+id)
}

func (s *FileStore) Delete(_ context.Context, id string) error {
	status, err := indexStatus(s.home, id)
	if err != nil {
		return err
	}
	if status == "" {
		return fmt.Errorf("channel %q not found", id)
	}
	var dir string
	switch status {
	case "active":
		dir = channelActiveDir(s.home, id)
	case "archive":
		dir = channelArchiveDir(s.home, id)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("delete channel dir: %w", err)
	}
	if err := os.Remove(channelIndexPath(s.home, id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	ts := struct {
		ID        string `json:"id"`
		DeletedAt string `json:"deleted_at"`
	}{
		ID:        id,
		DeletedAt: nowTimestamp(channel.Seq(id)),
	}
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tombDir := filepath.Join(s.channelsRoot(), "tombstones")
	if err := os.MkdirAll(tombDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(channelTombstonePath(s.home, id), data, 0o644)
}

func nextMessageID(channelDir string) (int, error) {
	path := filepath.Join(channelDir, "msg-counter")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return 0, fmt.Errorf("open msg-counter: %w", err)
	}
	defer f.Close()

	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return 0, fmt.Errorf("flock msg-counter: %w", err)
	}
	defer func() { _ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN) }()

	data, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("read msg-counter: %w", err)
	}
	cur := 0
	if len(data) > 0 {
		cur, err = strconv.Atoi(strings.TrimSpace(string(data)))
		if err != nil {
			return 0, fmt.Errorf("parse msg-counter: %w", err)
		}
	}
	next := cur + 1
	if err := os.WriteFile(path, []byte(strconv.Itoa(next)), 0o644); err != nil {
		return 0, fmt.Errorf("write msg-counter: %w", err)
	}
	return next, nil
}

func (s *FileStore) SendMessage(_ context.Context, req channel.SendMessageRequest) (*channel.Message, error) {
	loc, err := s.locate(req.ChannelID)
	if err != nil {
		return nil, err
	}
	if loc.status == "archive" {
		return nil, fmt.Errorf("channel %q is archived", req.ChannelID)
	}
	parts, err := readParticipants(loc.dir)
	if err != nil {
		return nil, err
	}
	if err := requireParticipant(parts, req.ChannelID, req.Sender); err != nil {
		return nil, err
	}
	msgID, err := nextMessageID(loc.dir)
	if err != nil {
		return nil, err
	}
	msg := channel.Message{
		ID:        msgID,
		Sender:    req.Sender,
		Body:      req.Body,
		CreatedAt: nowTimestamp(channel.Seq(req.ChannelID) + msgID),
	}
	line, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(filepath.Join(loc.dir, "messages.jsonl"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open messages.jsonl: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return nil, fmt.Errorf("append message: %w", err)
	}
	return &msg, nil
}

func (s *FileStore) ListMessages(_ context.Context, channelID string) ([]channel.Message, error) {
	loc, err := s.locate(channelID)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(loc.dir, "messages.jsonl")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var msgs []channel.Message
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m channel.Message
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			return nil, fmt.Errorf("parse messages.jsonl: %w", err)
		}
		msgs = append(msgs, m)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *FileStore) ListParticipants(_ context.Context, channelID string) ([]channel.Participant, error) {
	loc, err := s.locate(channelID)
	if err != nil {
		return nil, err
	}
	return readParticipants(loc.dir)
}

func (s *FileStore) AddParticipant(_ context.Context, req channel.ParticipantChangeRequest) (bool, error) {
	loc, err := s.locate(req.ChannelID)
	if err != nil {
		return false, err
	}
	if loc.status == "archive" {
		return false, fmt.Errorf("channel %q is archived", req.ChannelID)
	}
	parts, err := readParticipants(loc.dir)
	if err != nil {
		return false, err
	}
	if err := requireParticipant(parts, req.ChannelID, req.Actor); err != nil {
		return false, err
	}
	handle, err := channel.NormalizeHandle(req.Handle)
	if err != nil {
		return false, err
	}
	if isParticipant(parts, handle) {
		return false, nil
	}
	now := nowTimestamp(channel.Seq(req.ChannelID))
	parts = append(parts, channel.Participant{Handle: handle, JoinedAt: now})
	if err := writeParticipants(loc.dir, parts); err != nil {
		return false, err
	}
	ch, err := readChannelMetadata(loc.dir)
	if err != nil {
		return false, err
	}
	ch.UpdatedAt = now
	if err := writeChannelMetadata(loc.dir, ch); err != nil {
		return false, err
	}
	return true, nil
}

func (s *FileStore) RemoveParticipant(_ context.Context, req channel.ParticipantChangeRequest) error {
	loc, err := s.locate(req.ChannelID)
	if err != nil {
		return err
	}
	if loc.status == "archive" {
		return fmt.Errorf("channel %q is archived", req.ChannelID)
	}
	parts, err := readParticipants(loc.dir)
	if err != nil {
		return err
	}
	if err := requireParticipant(parts, req.ChannelID, req.Actor); err != nil {
		return err
	}
	handle, err := channel.NormalizeHandle(req.Handle)
	if err != nil {
		return err
	}
	if !isParticipant(parts, handle) {
		return fmt.Errorf("participant %q not in channel %q", handle, req.ChannelID)
	}
	if len(parts) <= 1 {
		return fmt.Errorf("cannot remove the last participant from channel %q", req.ChannelID)
	}
	var kept []channel.Participant
	for _, p := range parts {
		if p.Handle != handle {
			kept = append(kept, p)
		}
	}
	if err := writeParticipants(loc.dir, kept); err != nil {
		return err
	}
	ch, err := readChannelMetadata(loc.dir)
	if err != nil {
		return err
	}
	ch.UpdatedAt = nowTimestamp(channel.Seq(req.ChannelID))
	return writeChannelMetadata(loc.dir, ch)
}