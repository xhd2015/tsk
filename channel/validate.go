package channel

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

var idRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// Slugify converts a name into a channel id slug.
func Slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
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
	runes := []rune(out)
	if len(runes) > 64 {
		out = string(runes[:64])
		out = strings.Trim(out, "-")
	}
	return out
}

// ValidateID checks channel id format.
func ValidateID(id string) error {
	id = strings.ToLower(strings.TrimSpace(id))
	if !idRe.MatchString(id) {
		return fmt.Errorf("invalid channel id %q", id)
	}
	return nil
}

// NormalizeHandle lowercases and validates a participant handle.
func NormalizeHandle(handle string) (string, error) {
	handle = strings.ToLower(strings.TrimSpace(handle))
	if !idRe.MatchString(handle) {
		return "", fmt.Errorf("invalid handle %q", handle)
	}
	return handle, nil
}

// ResolveIdentity returns the current user handle.
// Precedence: userFlag > TSK_USER env > $USER.
func ResolveIdentity(userFlag string) (string, error) {
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

// Seq returns a deterministic sequence number from a channel id.
func Seq(id string) int {
	h := 0
	for _, c := range id {
		h = h*31 + int(c)
	}
	if h < 0 {
		h = -h
	}
	return h
}