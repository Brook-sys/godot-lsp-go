package rewriter

import (
	"strings"
	"testing"
)

func TestFileURIToPath(t *testing.T) {
	cases := map[string]string{
		"file:///home/alpine/game/a.gd":     "/home/alpine/game/a.gd",
		"file:///C:/Users/Murilo/game/a.gd": "C:/Users/Murilo/game/a.gd",
		"file:///C%3A/Users/Murilo/game.gd": "C:/Users/Murilo/game.gd",
		`file://C:\Users\Murilo\game\a.gd`:  "C:/Users/Murilo/game/a.gd",
	}
	for input, want := range cases {
		got, ok := FileURIToPath(input)
		if !ok {
			t.Fatalf("expected %s to parse", input)
		}
		if got != want {
			t.Fatalf("%s: want %s got %s", input, want, got)
		}
	}
}

func TestMapFileURIClientToGodot(t *testing.T) {
	maps := []PathMap{{ClientRoot: "/home/alpine/game", GodotRoot: "C:/Users/Murilo/Godot/game"}}
	got, changed := MapFileURI("file:///home/alpine/game/scripts/player.gd", maps, ClientToGodot)
	if !changed {
		t.Fatal("expected change")
	}
	want := "file:///C:/Users/Murilo/Godot/game/scripts/player.gd"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestMapFileURIGodotToClient(t *testing.T) {
	maps := []PathMap{{ClientRoot: "/home/alpine/game", GodotRoot: "C:/Users/Murilo/Godot/game"}}
	got, changed := MapFileURI("file:///C:/Users/Murilo/Godot/game/scripts/player.gd", maps, GodotToClient)
	if !changed {
		t.Fatal("expected change")
	}
	want := "file:///home/alpine/game/scripts/player.gd"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestPathMapDoesNotMatchSiblingPrefix(t *testing.T) {
	maps := []PathMap{{ClientRoot: "/home/alpine/game", GodotRoot: "/remote/game"}}
	_, changed := MapPlainPath("/home/alpine/gameplay/a.gd", maps, ClientToGodot)
	if changed {
		t.Fatal("must not map sibling prefix")
	}
}

func TestMostSpecificPathMapWins(t *testing.T) {
	maps := []PathMap{
		{ClientRoot: "/home/alpine/game", GodotRoot: "/remote/game"},
		{ClientRoot: "/home/alpine/game/addons", GodotRoot: "/remote/addons"},
	}
	got, changed := MapPlainPath("/home/alpine/game/addons/plugin.gd", maps, ClientToGodot)
	if !changed {
		t.Fatal("expected change")
	}
	want := "/remote/addons/plugin.gd"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestRewriteInitializeRootPath(t *testing.T) {
	maps := []PathMap{{ClientRoot: "/home/alpine/game", GodotRoot: "C:/Users/Murilo/Godot/game"}}
	in := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"rootPath":"/home/alpine/game","rootUri":"file:///home/alpine/game"}}`)
	out := Rewrite(in, Options{NormalizeURIs: true, PathMaps: maps, Direction: ClientToGodot})
	if !strings.Contains(string(out), `"rootPath":"C:/Users/Murilo/Godot/game"`) {
		t.Fatalf("rootPath not mapped: %s", out)
	}
	if !strings.Contains(string(out), `"rootUri":"file:///C:/Users/Murilo/Godot/game"`) {
		t.Fatalf("rootUri not mapped: %s", out)
	}
}

func TestRewriteGodotChangeWorkspacePath(t *testing.T) {
	maps := []PathMap{{ClientRoot: "/home/alpine/game", GodotRoot: "C:/Users/Murilo/Godot/game"}}
	in := []byte(`{"jsonrpc":"2.0","method":"gdscript_client/changeWorkspace","params":{"path":"C:/Users/Murilo/Godot/game"}}`)
	out := Rewrite(in, Options{PathMaps: maps, Direction: GodotToClient})
	if !strings.Contains(string(out), `"path":"/home/alpine/game"`) {
		t.Fatalf("path not mapped: %s", out)
	}
}
