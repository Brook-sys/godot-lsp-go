package rewriter

import (
	"net/url"
	"sort"
	"strings"
)

type Direction int

const (
	ClientToGodot Direction = iota
	GodotToClient
)

type PathMap struct {
	ClientRoot string
	GodotRoot  string
}

func ParsePathMap(value string) (PathMap, error) {
	left, right, ok := strings.Cut(value, "=")
	if !ok {
		return PathMap{}, errInvalidPathMap(value)
	}
	left = canonicalPathRoot(left)
	right = canonicalPathRoot(right)
	if left == "" || right == "" {
		return PathMap{}, errInvalidPathMap(value)
	}
	return PathMap{ClientRoot: left, GodotRoot: right}, nil
}

func NormalizePathMaps(maps []PathMap) []PathMap {
	out := make([]PathMap, 0, len(maps))
	for _, m := range maps {
		m.ClientRoot = canonicalPathRoot(m.ClientRoot)
		m.GodotRoot = canonicalPathRoot(m.GodotRoot)
		if m.ClientRoot != "" && m.GodotRoot != "" {
			out = append(out, m)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		return len(out[i].ClientRoot)+len(out[i].GodotRoot) > len(out[j].ClientRoot)+len(out[j].GodotRoot)
	})
	return out
}

func MapFileURI(uri string, maps []PathMap, direction Direction) (string, bool) {
	path, ok := FileURIToPath(uri)
	if !ok {
		return uri, false
	}
	mapped, changed := MapPlainPath(path, maps, direction)
	if !changed {
		return NormalizeFileURI(uri), NormalizeFileURI(uri) != uri
	}
	return PathToFileURI(mapped), true
}

func MapPlainPath(path string, maps []PathMap, direction Direction) (string, bool) {
	path = canonicalPathRoot(path)
	for _, m := range NormalizePathMaps(maps) {
		from := m.ClientRoot
		to := m.GodotRoot
		if direction == GodotToClient {
			from = m.GodotRoot
			to = m.ClientRoot
		}
		if hasPathPrefix(path, from) {
			rest := strings.TrimPrefix(path, from)
			rest = strings.TrimPrefix(rest, "/")
			if rest == "" {
				return to, true
			}
			return strings.TrimRight(to, "/") + "/" + rest, true
		}
	}
	return path, false
}

func FileURIToPath(uri string) (string, bool) {
	if !strings.HasPrefix(uri, "file://") {
		return "", false
	}
	rest := strings.TrimPrefix(uri, "file://")
	rest = strings.ReplaceAll(rest, "\\", "/")
	if strings.HasPrefix(rest, "localhost/") {
		rest = strings.TrimPrefix(rest, "localhost")
	}
	if strings.HasPrefix(rest, "/") && len(rest) >= 4 && isDrive(rest[1:3]) && (rest[3] == '/' || rest[3] == '\\') {
		rest = rest[1:]
	}
	decoded, err := url.PathUnescape(rest)
	if err == nil {
		rest = decoded
	}
	if strings.HasPrefix(rest, "/") && len(rest) >= 4 && isDrive(rest[1:3]) && (rest[3] == '/' || rest[3] == '\\') {
		rest = rest[1:]
	}
	return canonicalPathRoot(rest), true
}

func PathToFileURI(path string) string {
	path = canonicalPathRoot(path)
	if isWindowsPath(path) {
		path = "/" + path
	}
	parts := strings.Split(path, "/")
	for i, part := range parts {
		parts[i] = strings.ReplaceAll(url.PathEscape(part), "%3A", ":")
	}
	return "file://" + strings.Join(parts, "/")
}

func NormalizeFileURI(uri string) string {
	path, ok := FileURIToPath(uri)
	if !ok {
		return uri
	}
	return PathToFileURI(path)
}

func canonicalPathRoot(path string) string {
	path = strings.TrimSpace(path)
	if strings.HasPrefix(path, "file://") {
		if parsed, ok := FileURIToPath(path); ok {
			path = parsed
		}
	}
	path = strings.ReplaceAll(path, "\\", "/")
	for strings.Contains(path, "//") {
		path = strings.ReplaceAll(path, "//", "/")
	}
	if strings.HasSuffix(path, "/") && path != "/" && !isDrive(path[:len(path)-1]) {
		path = strings.TrimRight(path, "/")
	}
	if len(path) >= 2 && isDrive(path[:2]) {
		path = strings.ToUpper(path[:1]) + path[1:]
	}
	return path
}

func hasPathPrefix(path, root string) bool {
	path = canonicalPathRoot(path)
	root = canonicalPathRoot(root)
	return path == root || strings.HasPrefix(path, strings.TrimRight(root, "/")+"/")
}

func isWindowsPath(path string) bool {
	return len(path) >= 2 && isDrive(path[:2])
}

func isDrive(s string) bool {
	if len(s) != 2 || s[1] != ':' {
		return false
	}
	c := s[0]
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

type pathMapError string

func (e pathMapError) Error() string { return "invalid path map: " + string(e) }

func errInvalidPathMap(value string) error { return pathMapError(value) }
