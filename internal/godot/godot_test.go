package godot

import "testing"

func TestCommandArgs(t *testing.T) {
	args := CommandArgs(6005, "/tmp/project", true)
	joined := ""
	for _, a := range args {
		joined += a + " "
	}
	if joined == "" || args[0] != "--editor" {
		t.Fatalf("unexpected args: %v", args)
	}
}
