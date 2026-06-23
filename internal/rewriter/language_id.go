package rewriter

import "strings"

func patchLanguageID(v any) bool {
	msg, ok := v.(map[string]any)
	if !ok {
		return false
	}
	method, _ := msg["method"].(string)
	if method != "textDocument/didOpen" {
		return false
	}
	params, ok := msg["params"].(map[string]any)
	if !ok {
		return false
	}
	doc, ok := params["textDocument"].(map[string]any)
	if !ok {
		return false
	}
	languageID, _ := doc["languageId"].(string)
	uri, _ := doc["uri"].(string)
	if languageID != "plaintext" {
		return false
	}
	if !(strings.HasSuffix(uri, ".gd") || strings.HasSuffix(uri, ".gdshader")) {
		return false
	}
	doc["languageId"] = "gdscript"
	return true
}
