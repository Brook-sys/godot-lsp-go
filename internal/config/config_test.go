package config

import "testing"

func TestParsePathMapFlag(t *testing.T) {
	cfg, err := Parse([]string{"--path-map", "/client=/godot", "--path-map", "/client/addons=/godot/addons"})
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.PathMaps) != 2 {
		t.Fatalf("expected 2 path maps, got %d", len(cfg.PathMaps))
	}
	if cfg.PathMaps[0].ClientRoot != "/client/addons" {
		t.Fatalf("expected most specific map first, got %#v", cfg.PathMaps)
	}
}
