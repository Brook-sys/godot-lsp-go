package godot

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func FindProject(start string) (string, error) {
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	if info, err := os.Stat(abs); err == nil && !info.IsDir() {
		abs = filepath.Dir(abs)
	}
	current := abs
	for i := 0; i < 10; i++ {
		candidate := filepath.Join(current, "project.godot")
		if _, err := os.Stat(candidate); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("project.godot not found from %s", abs)
}

func FindExecutable(configured string) (string, error) {
	if configured != "" {
		if _, err := os.Stat(configured); err == nil {
			return configured, nil
		}
		if p, err := exec.LookPath(configured); err == nil {
			return p, nil
		}
		return "", fmt.Errorf("godot executable not found: %s", configured)
	}
	names := []string{"godot", "godot4", "godot-editor"}
	if runtime.GOOS == "windows" {
		names = append([]string{"godot.exe"}, names...)
	}
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return p, nil
		}
	}
	return "", fmt.Errorf("godot executable not found")
}

func CommandArgs(port int, project string, headless bool) []string {
	args := []string{"--editor"}
	if headless {
		args = append(args, "--headless", "--display-driver", "headless", "--audio-driver", "Dummy")
	}
	args = append(args, "--lsp-port", fmt.Sprintf("%d", port), "--path", project)
	return args
}

func Launch(ctx context.Context, executable string, args []string) (*exec.Cmd, error) {
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = nil
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
