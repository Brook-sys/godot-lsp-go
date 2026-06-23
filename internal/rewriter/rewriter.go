package rewriter

import (
	"encoding/json"
	"strings"
)

type Options struct {
	NormalizeURIs bool
	PatchOpenCode bool
}

func Rewrite(body []byte, opts Options) []byte {
	var msg any
	if err := json.Unmarshal(body, &msg); err != nil {
		return body
	}
	changed := false
	if opts.NormalizeURIs {
		msg, changed = normalizeValue(msg)
	}
	if opts.PatchOpenCode {
		if patchLanguageID(msg) {
			changed = true
		}
	}
	if !changed {
		return body
	}
	out, err := json.Marshal(msg)
	if err != nil {
		return body
	}
	return out
}

func normalizeValue(v any) (any, bool) {
	switch t := v.(type) {
	case string:
		n := NormalizeFileURI(t)
		return n, n != t
	case []any:
		changed := false
		for i, item := range t {
			var c bool
			t[i], c = normalizeValue(item)
			changed = changed || c
		}
		return t, changed
	case map[string]any:
		changed := false
		for k, item := range t {
			var c bool
			t[k], c = normalizeValue(item)
			changed = changed || c
		}
		return t, changed
	default:
		return v, false
	}
}

func NormalizeFileURI(uri string) string {
	if !strings.HasPrefix(uri, "file://") {
		return uri
	}
	path := strings.TrimPrefix(uri, "file://")
	path = strings.ReplaceAll(path, "\\", "/")
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return "file://" + path
}
