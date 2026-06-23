package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/code-xhyun/godot-lsp-go/internal/bridge"
	"github.com/code-xhyun/godot-lsp-go/internal/config"
	"github.com/code-xhyun/godot-lsp-go/internal/connector"
	"github.com/code-xhyun/godot-lsp-go/internal/godot"
	"github.com/code-xhyun/godot-lsp-go/internal/logging"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}
	if cfg.ShowVersion {
		os.Exit(0)
	}

	log, err := logging.New(cfg.Debug, cfg.Quiet, cfg.LogFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "log error: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-sigs
		cancel()
	}()

	var cmd *exec.Cmd
	if cfg.Launch {
		open := false
		for _, port := range cfg.Ports {
			if connector.IsPortOpen(ctx, cfg.Host, port, cfg.ConnectTimeout) {
				open = true
				break
			}
		}
		if !open {
			project, err := godot.FindProject(cfg.ProjectPath)
			if err != nil {
				log.Error("failed to find project: %v", err)
				os.Exit(1)
			}
			exe, err := godot.FindExecutable(cfg.GodotPath)
			if err != nil {
				log.Error("failed to find godot executable: %v", err)
				os.Exit(1)
			}
			args := godot.CommandArgs(cfg.Ports[0], project, cfg.Headless)
			log.Info("launching godot: %s %v", exe, args)
			cmd, err = godot.Launch(ctx, exe, args)
			if err != nil {
				log.Error("failed to launch godot: %v", err)
				os.Exit(1)
			}
			if cfg.CleanupGodot != "never" {
				defer func() {
					if cmd.Process != nil {
						_ = cmd.Process.Kill()
					}
				}()
			}
			log.Info("waiting for Godot LSP to start")
			time.Sleep(2 * time.Second)
		}
	}

	b := bridge.New(cfg, log)
	if err := b.Run(ctx); err != nil {
		log.Error("bridge error: %v", err)
		os.Exit(1)
	}
}
