package rewriter

import "encoding/json"

type Options struct {
	NormalizeURIs bool
	PatchOpenCode bool
	PathMaps      []PathMap
	Direction     Direction
}

func Rewrite(body []byte, opts Options) []byte {
	var msg any
	if err := json.Unmarshal(body, &msg); err != nil {
		return body
	}
	changed := false
	if opts.PatchOpenCode && opts.Direction == ClientToGodot {
		if patchLanguageID(msg) {
			changed = true
		}
	}
	if opts.NormalizeURIs || len(opts.PathMaps) > 0 {
		var c bool
		msg, c = rewriteValue(msg, opts)
		changed = changed || c
	}
	if opts.Direction == ClientToGodot {
		if patchInitializeRootPath(msg, opts.PathMaps) {
			changed = true
		}
	} else {
		if patchGodotPlainPaths(msg, opts.PathMaps) {
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

func rewriteValue(v any, opts Options) (any, bool) {
	switch t := v.(type) {
	case string:
		if len(opts.PathMaps) > 0 {
			mapped, changed := MapFileURI(t, opts.PathMaps, opts.Direction)
			if changed {
				return mapped, true
			}
		}
		if opts.NormalizeURIs {
			n := NormalizeFileURI(t)
			return n, n != t
		}
		return t, false
	case []any:
		changed := false
		for i, item := range t {
			var c bool
			t[i], c = rewriteValue(item, opts)
			changed = changed || c
		}
		return t, changed
	case map[string]any:
		changed := false
		for k, item := range t {
			var c bool
			t[k], c = rewriteValue(item, opts)
			changed = changed || c
		}
		return t, changed
	default:
		return v, false
	}
}

func patchInitializeRootPath(v any, maps []PathMap) bool {
	msg, ok := v.(map[string]any)
	if !ok {
		return false
	}
	method, _ := msg["method"].(string)
	if method != "initialize" {
		return false
	}
	params, ok := msg["params"].(map[string]any)
	if !ok {
		return false
	}
	rootPath, ok := params["rootPath"].(string)
	if !ok || rootPath == "" {
		return false
	}
	mapped, changed := MapPlainPath(rootPath, maps, ClientToGodot)
	if changed {
		params["rootPath"] = mapped
	}
	return changed
}

func patchGodotPlainPaths(v any, maps []PathMap) bool {
	msg, ok := v.(map[string]any)
	if !ok {
		return false
	}
	method, _ := msg["method"].(string)
	if method != "gdscript_client/changeWorkspace" {
		return false
	}
	params, ok := msg["params"].(map[string]any)
	if !ok {
		return false
	}
	path, ok := params["path"].(string)
	if !ok || path == "" {
		return false
	}
	mapped, changed := MapPlainPath(path, maps, GodotToClient)
	if changed {
		params["path"] = mapped
	}
	return changed
}
