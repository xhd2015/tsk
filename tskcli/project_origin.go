package tskcli

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// NormalizeOriginURL converts a git remote.origin.url into a canonical
// project id (host/path, no scheme or .git) and short name (basename).
func NormalizeOriginURL(origin string) (id, name string, err error) {
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return "", "", fmt.Errorf("empty origin url")
	}

	var host, repoPath string

	switch {
	case strings.Contains(origin, "://"):
		u, perr := url.Parse(origin)
		if perr != nil {
			return "", "", fmt.Errorf("parse origin url: %w", perr)
		}
		host = u.Hostname()
		repoPath = strings.TrimPrefix(u.Path, "/")
	case strings.HasPrefix(origin, "git@") || scpLikeOrigin(origin):
		// user@host:path or git@host:path
		at := strings.IndexByte(origin, '@')
		rest := origin[at+1:]
		colon := strings.IndexByte(rest, ':')
		if colon < 0 {
			return "", "", fmt.Errorf("invalid scp-style origin %q", origin)
		}
		host = rest[:colon]
		repoPath = rest[colon+1:]
	default:
		// bare host/path
		u, perr := url.Parse("https://" + origin)
		if perr != nil {
			return "", "", fmt.Errorf("parse origin url: %w", perr)
		}
		host = u.Hostname()
		repoPath = strings.TrimPrefix(u.Path, "/")
	}

	host = strings.TrimSpace(host)
	repoPath = strings.Trim(strings.TrimSpace(repoPath), "/")
	repoPath = strings.TrimSuffix(repoPath, ".git")
	if host == "" || repoPath == "" {
		return "", "", fmt.Errorf("invalid origin url %q", origin)
	}

	id = host + "/" + repoPath
	name = path.Base(repoPath)
	if name == "" || name == "." || name == "/" {
		return "", "", fmt.Errorf("invalid origin url %q", origin)
	}
	return id, name, nil
}

func scpLikeOrigin(origin string) bool {
	at := strings.IndexByte(origin, '@')
	if at <= 0 {
		return false
	}
	rest := origin[at+1:]
	colon := strings.IndexByte(rest, ':')
	return colon > 0 && !strings.Contains(rest[:colon], "/")
}
