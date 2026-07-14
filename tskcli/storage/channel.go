package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
)

var channelIDRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// ChannelParticipant is one roster entry in channel.json.
type ChannelParticipant struct {
	Handle   string `json:"handle"`
	JoinedAt string `json:"joined_at"`
}

// Channel is the on-disk channel.json schema.
type Channel struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Status       string               `json:"status"`
	Participants []ChannelParticipant `json:"participants"`
	CreatedAt    string               `json:"created_at"`
	UpdatedAt    string               `json:"updated_at"`
}

// ChannelMessage is one line in messages.jsonl.
type ChannelMessage struct {
	ID        int    `json:"id"`
	Sender    string `json:"sender"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// ChannelTombstone blocks channel id reuse after delete.
type ChannelTombstone struct {
	ID        string `json:"id"`
	DeletedAt string `json:"deleted_at"`
}

// ChannelListEntry is a channel summary for list output.
type ChannelListEntry struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

func channelsRoot(home string) string {
	return filepath.Join(home, "channels")
}

func channelIndexPath(home, id string) string {
	return filepath.Join(channelsRoot(home), "index", id)
}

func channelTombstonePath(home, id string) string {
	return filepath.Join(channelsRoot(home), "tombstones", id)
}

func channelActiveDir(home, id string) string {
	return filepath.Join(channelsRoot(home), "active", id)
}

func channelArchiveDir(home, id string) string {
	return filepath.Join(channelsRoot(home), "archive", id)
}

// ChannelSeq returns a deterministic sequence number from a channel id.
func ChannelSeq(id string) int {
	h := 0
	for _, c := range id {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}

// ValidateChannelID checks channel id format.
func ValidateChannelID(id string) error {
	id = strings.ToLower(strings.TrimSpace(id))
	if !channelIDRe.MatchString(id) {
		return fmt.Errorf("invalid channel id %q", id)
	}
	return nil
}

// NormalizeHandle lowercases and validates a participant handle.
func NormalizeHandle(handle string) (string, error) {
	handle = strings.ToLower(strings.TrimSpace(handle))
	if !channelIDRe.MatchString(handle) {
		return "", fmt.Errorf("invalid handle %q", handle)
	}
	return handle, nil
}

// ResolveChannelIdentity returns the current user handle.
// Precedence: --user flag > TSK_USER env > $USER.
func ResolveChannelIdentity(userFlag string) (string, error) {
	if userFlag != "" {
		return NormalizeHandle(userFlag)
	}
	if v := os.Getenv("TSK_USER"); v != "" {
		return NormalizeHandle(v)
	}
	user := os.Getenv("USER")
	if user == "" {
		return "", fmt.Errorf("cannot resolve user identity")
	}
	return NormalizeHandle(user)
}

// ChannelIndexStatus returns active, archive, or empty if missing.
func ChannelIndexStatus(home, id string) (string, error) {
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

// ChannelTombstoned reports whether a channel id was deleted.
func ChannelTombstoned(home, id string) (bool, error) {
	_, err := os.Stat(channelTombstonePath(home, id))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ChannelIDAvailable reports whether a channel id can be created.
func ChannelIDAvailable(home, id string) (bool, error) {
	status, err := ChannelIndexStatus(home, id)
	if err != nil {
		return false, err
	}
	if status != "" {
		return false, nil
	}
	tomb, err := ChannelTombstoned(home, id)
	if err != nil {
		return false, err
	}
	return !tomb, nil
}

func writeChannelIndex(home, id, rel string) error {
	dir := filepath.Join(channelsRoot(home), "index")
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

func sortParticipants(parts []ChannelParticipant) {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Handle < parts[j].Handle
	})
}

// WriteChannel writes channel.json atomically with sorted participants.
func WriteChannel(channelDir string, ch Channel) error {
	sortParticipants(ch.Participants)
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

// ReadChannel loads channel.json from a channel directory.
func ReadChannel(channelDir string) (Channel, error) {
	data, err := os.ReadFile(filepath.Join(channelDir, "channel.json"))
	if err != nil {
		return Channel{}, fmt.Errorf("read channel.json: %w", err)
	}
	var ch Channel
	if err := json.Unmarshal(data, &ch); err != nil {
		return Channel{}, fmt.Errorf("parse channel.json: %w", err)
	}
	return ch, nil
}

// LoadChannelByID loads a channel by id from active or archive storage.
func LoadChannelByID(home, id string) (Channel, string, string, error) {
	status, err := ChannelIndexStatus(home, id)
	if err != nil {
		return Channel{}, "", "", err
	}
	if status == "" {
		return Channel{}, "", "", fmt.Errorf("channel %q not found", id)
	}
	var dir string
	switch status {
	case "active":
		dir = channelActiveDir(home, id)
	case "archive":
		dir = channelArchiveDir(home, id)
	default:
		return Channel{}, "", "", fmt.Errorf("channel %q not found", id)
	}
	ch, err := ReadChannel(dir)
	if err != nil {
		return Channel{}, "", "", err
	}
	return ch, dir, status, nil
}

func channelIsParticipant(ch Channel, handle string) bool {
	for _, p := range ch.Participants {
		if p.Handle == handle {
			return true
		}
	}
	return false
}

// RequireChannelParticipant verifies membership.
func RequireChannelParticipant(ch Channel, handle string) error {
	if !channelIsParticipant(ch, handle) {
		return fmt.Errorf("not a participant in channel %q", ch.ID)
	}
	return nil
}

// CreateChannel creates a new active channel with creator and agent participants.
func CreateChannel(home, name, id, creator string) error {
	if err := ValidateChannelID(id); err != nil {
		return err
	}
	avail, err := ChannelIDAvailable(home, id)
	if err != nil {
		return err
	}
	if !avail {
		return fmt.Errorf("channel %q already exists", id)
	}

	dir := channelActiveDir(home, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create channel dir: %w", err)
	}

	now := NowTimestamp(ChannelSeq(id))
	ch := Channel{
		ID:     id,
		Name:   name,
		Status: "active",
		Participants: []ChannelParticipant{
			{Handle: "agent", JoinedAt: now},
			{Handle: creator, JoinedAt: now},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := WriteChannel(dir, ch); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "messages.jsonl"), nil, 0o644); err != nil {
		return fmt.Errorf("create messages.jsonl: %w", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "msg-counter"), []byte("0"), 0o644); err != nil {
		return fmt.Errorf("create msg-counter: %w", err)
	}
	return writeChannelIndex(home, id, "active/"+id)
}

// ListChannels returns channel summaries; active only unless all is true.
func ListChannels(home string, all bool) ([]ChannelListEntry, error) {
	indexDir := filepath.Join(channelsRoot(home), "index")
	entries, err := os.ReadDir(indexDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []ChannelListEntry
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		id := ent.Name()
		status, err := ChannelIndexStatus(home, id)
		if err != nil {
			continue
		}
		if status == "" {
			continue
		}
		if !all && status == "archive" {
			continue
		}
		ch, _, _, err := LoadChannelByID(home, id)
		if err != nil {
			continue
		}
		out = append(out, ChannelListEntry{
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

// ArchiveChannel moves an active channel to archive storage.
func ArchiveChannel(home, id string) error {
	status, err := ChannelIndexStatus(home, id)
	if err != nil {
		return err
	}
	if status == "" {
		return fmt.Errorf("channel %q not found", id)
	}
	if status == "archive" {
		return fmt.Errorf("channel %q is already archived", id)
	}
	activeDir := channelActiveDir(home, id)
	archiveDir := channelArchiveDir(home, id)
	if err := os.MkdirAll(filepath.Dir(archiveDir), 0o755); err != nil {
		return err
	}
	if err := os.Rename(activeDir, archiveDir); err != nil {
		return fmt.Errorf("archive channel: %w", err)
	}
	ch, err := ReadChannel(archiveDir)
	if err != nil {
		return err
	}
	ch.Status = "archived"
	ch.UpdatedAt = NowTimestamp(ChannelSeq(id))
	if err := WriteChannel(archiveDir, ch); err != nil {
		return err
	}
	return writeChannelIndex(home, id, "archive/"+id)
}

// DeleteChannel removes a channel and writes a tombstone.
func DeleteChannel(home, id string) error {
	status, err := ChannelIndexStatus(home, id)
	if err != nil {
		return err
	}
	if status == "" {
		return fmt.Errorf("channel %q not found", id)
	}
	var dir string
	switch status {
	case "active":
		dir = channelActiveDir(home, id)
	case "archive":
		dir = channelArchiveDir(home, id)
	}
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("delete channel dir: %w", err)
	}
	if err := os.Remove(channelIndexPath(home, id)); err != nil && !os.IsNotExist(err) {
		return err
	}
	ts := ChannelTombstone{
		ID:        id,
		DeletedAt: NowTimestamp(ChannelSeq(id)),
	}
	data, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tombDir := filepath.Join(channelsRoot(home), "tombstones")
	if err := os.MkdirAll(tombDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(channelTombstonePath(home, id), data, 0o644)
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

// SendChannelMessage appends a message from a participant to an active channel.
func SendChannelMessage(home, id, sender, body string) (int, error) {
	ch, dir, status, err := LoadChannelByID(home, id)
	if err != nil {
		return 0, err
	}
	if status == "archive" {
		return 0, fmt.Errorf("channel %q is archived", id)
	}
	if err := RequireChannelParticipant(ch, sender); err != nil {
		return 0, err
	}
	msgID, err := nextMessageID(dir)
	if err != nil {
		return 0, err
	}
	msg := ChannelMessage{
		ID:        msgID,
		Sender:    sender,
		Body:      body,
		CreatedAt: NowTimestamp(ChannelSeq(id) + msgID),
	}
	line, err := json.Marshal(msg)
	if err != nil {
		return 0, err
	}
	f, err := os.OpenFile(filepath.Join(dir, "messages.jsonl"), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, fmt.Errorf("open messages.jsonl: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(append(line, '\n')); err != nil {
		return 0, fmt.Errorf("append message: %w", err)
	}
	return msgID, nil
}

// ReadChannelMessages loads messages in chronological order.
func ReadChannelMessages(channelDir string) ([]ChannelMessage, error) {
	path := filepath.Join(channelDir, "messages.jsonl")
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()
	var msgs []ChannelMessage
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m ChannelMessage
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

// AddChannelParticipant adds a handle to an active channel (idempotent).
func AddChannelParticipant(home, id, actor, handle string) (bool, error) {
	ch, dir, status, err := LoadChannelByID(home, id)
	if err != nil {
		return false, err
	}
	if status == "archive" {
		return false, fmt.Errorf("channel %q is archived", id)
	}
	if err := RequireChannelParticipant(ch, actor); err != nil {
		return false, err
	}
	handle, err = NormalizeHandle(handle)
	if err != nil {
		return false, err
	}
	if channelIsParticipant(ch, handle) {
		return false, nil
	}
	now := NowTimestamp(ChannelSeq(id))
	ch.Participants = append(ch.Participants, ChannelParticipant{
		Handle:   handle,
		JoinedAt: now,
	})
	ch.UpdatedAt = now
	if err := WriteChannel(dir, ch); err != nil {
		return false, err
	}
	return true, nil
}

// RemoveChannelParticipant removes a handle from an active channel.
func RemoveChannelParticipant(home, id, actor, handle string) error {
	ch, dir, status, err := LoadChannelByID(home, id)
	if err != nil {
		return err
	}
	if status == "archive" {
		return fmt.Errorf("channel %q is archived", id)
	}
	if err := RequireChannelParticipant(ch, actor); err != nil {
		return err
	}
	handle, err = NormalizeHandle(handle)
	if err != nil {
		return err
	}
	if !channelIsParticipant(ch, handle) {
		return fmt.Errorf("participant %q not in channel %q", handle, id)
	}
	if len(ch.Participants) <= 1 {
		return fmt.Errorf("cannot remove the last participant from channel %q", id)
	}
	var kept []ChannelParticipant
	for _, p := range ch.Participants {
		if p.Handle != handle {
			kept = append(kept, p)
		}
	}
	ch.Participants = kept
	ch.UpdatedAt = NowTimestamp(ChannelSeq(id))
	return WriteChannel(dir, ch)
}