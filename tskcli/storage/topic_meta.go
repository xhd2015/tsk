package storage

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// TopicMetaFile is the reserved metadata filename in a topic directory.
const TopicMetaFile = "topic.json"

// TopicMeta is topics/<path>/topic.json.
type TopicMeta struct {
	Path      []string `json:"path"`
	Title     string   `json:"title"`
	Aliases   []string `json:"aliases"`
	Notes     string   `json:"notes,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// SplitTopicPath splits a slash-separated topic ref into segments.
func SplitTopicPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// JoinTopicPath joins topic segments with /.
func JoinTopicPath(parts []string) string {
	return strings.Join(parts, "/")
}

// TopicAbs returns the absolute directory for a topic path.
func TopicAbs(home string, parts []string) string {
	elems := append([]string{home, "topics"}, parts...)
	return filepath.Join(elems...)
}

// TopicDirExists reports whether the topic directory exists.
func TopicDirExists(home string, parts []string) bool {
	if len(parts) == 0 {
		return false
	}
	st, err := os.Stat(TopicAbs(home, parts))
	return err == nil && st.IsDir()
}

// IsTaskDirName reports whether name looks like [id]-<slug>.
func IsTaskDirName(name string) bool {
	_, _, ok := ParseTaskDirName(name)
	return ok
}

// ParseTaskDirName splits [id]-<slug>. ok is false if the name is not a task dir.
func ParseTaskDirName(name string) (id int, slug string, ok bool) {
	if !strings.HasPrefix(name, "[") {
		return 0, "", false
	}
	closeIdx := strings.IndexByte(name, ']')
	if closeIdx < 2 {
		return 0, "", false
	}
	id, err := strconv.Atoi(name[1:closeIdx])
	if err != nil || id <= 0 {
		return 0, "", false
	}
	rest := name[closeIdx+1:]
	if !strings.HasPrefix(rest, "-") {
		return 0, "", false
	}
	slug = rest[1:]
	if slug == "" {
		return 0, "", false
	}
	return id, slug, true
}

// DefaultTopicMeta returns metadata for a directory that has no topic.json.
func DefaultTopicMeta(parts []string) TopicMeta {
	title := ""
	if n := len(parts); n > 0 {
		title = parts[n-1]
	}
	return TopicMeta{
		Path:    append([]string(nil), parts...),
		Title:   title,
		Aliases: []string{},
		Notes:   "",
	}
}

// ReadTopicMeta loads topic.json. os.IsNotExist if the file is missing.
func ReadTopicMeta(topicDir string) (TopicMeta, error) {
	data, err := os.ReadFile(filepath.Join(topicDir, TopicMetaFile))
	if err != nil {
		return TopicMeta{}, err
	}
	var meta TopicMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return TopicMeta{}, fmt.Errorf("parse topic.json: %w", err)
	}
	if meta.Aliases == nil {
		meta.Aliases = []string{}
	}
	return meta, nil
}

// LoadTopicMeta returns topic.json or defaults when the file is absent.
func LoadTopicMeta(topicDir string, parts []string) (TopicMeta, error) {
	meta, err := ReadTopicMeta(topicDir)
	if err == nil {
		meta.Path = append([]string(nil), parts...)
		if meta.Aliases == nil {
			meta.Aliases = []string{}
		}
		if strings.TrimSpace(meta.Title) == "" {
			def := DefaultTopicMeta(parts)
			meta.Title = def.Title
		}
		return meta, nil
	}
	if os.IsNotExist(err) {
		return DefaultTopicMeta(parts), nil
	}
	return TopicMeta{}, err
}

// WriteTopicMeta writes topic.json atomically.
func WriteTopicMeta(topicDir string, meta TopicMeta) error {
	cleaned := make([]string, 0, len(meta.Aliases))
	for _, a := range meta.Aliases {
		a = strings.TrimSpace(strings.Trim(a, "/"))
		if a != "" {
			cleaned = append(cleaned, a)
		}
	}
	sort.Strings(cleaned)
	uniq := cleaned[:0]
	for i, a := range cleaned {
		if i == 0 || a != cleaned[i-1] {
			uniq = append(uniq, a)
		}
	}
	meta.Aliases = uniq
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(topicDir, "topic-*.json.tmp")
	if err != nil {
		return fmt.Errorf("create topic temp: %w", err)
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return fmt.Errorf("write topic temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	dst := filepath.Join(topicDir, TopicMetaFile)
	if err := os.Rename(tmpName, dst); err != nil {
		_ = os.Remove(tmpName)
		return fmt.Errorf("rename topic.json: %w", err)
	}
	return nil
}

// ListTopicChildNames returns sorted sub-topic directory names (not task dirs).
func ListTopicChildNames(topicDir string) ([]string, error) {
	entries, err := os.ReadDir(topicDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, err
	}
	var names []string
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		if IsTaskDirName(ent.Name()) {
			continue
		}
		names = append(names, ent.Name())
	}
	sort.Strings(names)
	return names, nil
}

// CountTasksAtTopic counts index tasks whose topic_path equals parts exactly.
func CountTasksAtTopic(home string, parts []string) (int, error) {
	return countTasksForTopic(home, parts, false)
}

// CountTasksUnderTopic counts tasks whose topic_path equals parts or is under it.
func CountTasksUnderTopic(home string, parts []string) (int, error) {
	return countTasksForTopic(home, parts, true)
}

func countTasksForTopic(home string, parts []string, includeDescendants bool) (int, error) {
	ids, err := ListTaskIDs(home)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, id := range ids {
		task, _, err := LoadTaskByID(home, id)
		if err != nil {
			return 0, err
		}
		got, err := ParseTopicPath(task.TopicPath)
		if err != nil {
			return 0, err
		}
		if includeDescendants {
			if topicPathHasPrefix(got, parts) {
				n++
			}
		} else if topicPartsEqual(got, parts) {
			n++
		}
	}
	return n, nil
}

// topicPathHasPrefix reports whether path equals prefix or is under it.
func topicPathHasPrefix(path, prefix []string) bool {
	if len(prefix) == 0 || len(path) < len(prefix) {
		return false
	}
	for i := range prefix {
		if path[i] != prefix[i] {
			return false
		}
	}
	return true
}

func topicPartsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// LookupTopicRef resolves a user topic ref (path or alias) to canonical parts.
// A missing topic still returns the literal parts; callers check TopicDirExists.
// err is set only on an empty ref or an alias conflict.
func LookupTopicRef(home, ref string) ([]string, error) {
	raw := SplitTopicPath(ref)
	if len(raw) == 0 {
		return nil, fmt.Errorf("topic path required")
	}
	if TopicDirExists(home, raw) {
		return raw, nil
	}
	matches, err := aliasMatches(home, raw)
	if err != nil {
		return nil, err
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("alias conflict: %s", JoinTopicPath(raw))
	}
	if len(matches) == 1 {
		return matches[0], nil
	}
	return raw, nil
}

type topicAliasHit struct {
	path    []string
	aliases []string
}

func collectTopicAliases(home string) ([]topicAliasHit, error) {
	root := filepath.Join(home, "topics")
	var hits []topicAliasHit
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path == root {
			return nil
		}
		if IsTaskDirName(d.Name()) {
			return fs.SkipDir
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		parts := SplitTopicPath(filepath.ToSlash(rel))
		meta, err := ReadTopicMeta(path)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if len(meta.Aliases) == 0 {
			return nil
		}
		hits = append(hits, topicAliasHit{path: parts, aliases: meta.Aliases})
		return nil
	})
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return hits, nil
}

func aliasMatches(home string, raw []string) ([][]string, error) {
	hits, err := collectTopicAliases(home)
	if err != nil {
		return nil, err
	}
	joined := JoinTopicPath(raw)
	seen := map[string][]string{}
	for _, hit := range hits {
		for _, alias := range hit.aliases {
			alias = strings.Trim(alias, "/")
			if alias == "" {
				continue
			}
			var candidate []string
			switch {
			case joined == alias:
				candidate = append([]string(nil), hit.path...)
			case strings.HasPrefix(joined, alias+"/"):
				rest := strings.TrimPrefix(joined, alias+"/")
				candidate = append(append([]string(nil), hit.path...), SplitTopicPath(rest)...)
			default:
				continue
			}
			key := JoinTopicPath(candidate)
			if prev, ok := seen[key]; ok {
				_ = prev
				continue
			}
			seen[key] = candidate
		}
	}
	if len(seen) == 0 {
		return nil, nil
	}
	out := make([][]string, 0, len(seen))
	for _, p := range seen {
		out = append(out, p)
	}
	return out, nil
}

// FindAliasOwner returns the topic path that already claims alias, if any.
func FindAliasOwner(home, alias string) ([]string, error) {
	alias = strings.Trim(alias, "/")
	if alias == "" {
		return nil, nil
	}
	hits, err := collectTopicAliases(home)
	if err != nil {
		return nil, err
	}
	var owners [][]string
	for _, hit := range hits {
		for _, a := range hit.aliases {
			if strings.Trim(a, "/") == alias {
				owners = append(owners, hit.path)
				break
			}
		}
	}
	if len(owners) == 0 {
		return nil, nil
	}
	if len(owners) > 1 {
		return nil, fmt.Errorf("alias conflict: %s", alias)
	}
	return owners[0], nil
}
