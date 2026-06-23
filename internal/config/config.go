package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/code-xhyun/godot-lsp-go/internal/rewriter"
	"github.com/code-xhyun/godot-lsp-go/internal/version"
)

type Config struct {
	Host               string
	Ports              []int
	Launch             bool
	GodotPath          string
	ProjectPath        string
	Headless           bool
	XVFB               string
	CleanupGodot       string
	Reconnect          bool
	ConnectTimeout     time.Duration
	StartupTimeout     time.Duration
	ReconnectDelay     time.Duration
	WarmupDelay        time.Duration
	MaxBufferSize      int
	MaxPendingMessages int
	NormalizeURIs      bool
	PatchOpenCode      bool
	PathMaps           []rewriter.PathMap
	Debug              bool
	Quiet              bool
	LogFile            string
	ShowVersion        bool
}

func Parse(args []string) (Config, error) {
	cfg := Config{
		Host:               envString("GODOT_LSP_HOST", "127.0.0.1"),
		Ports:              envPorts("GODOT_LSP_PORTS", envPort("GODOT_LSP_PORT", []int{6005, 6007, 6008})),
		Headless:           true,
		XVFB:               "auto",
		CleanupGodot:       "auto",
		Reconnect:          true,
		ConnectTimeout:     2 * time.Second,
		StartupTimeout:     30 * time.Second,
		ReconnectDelay:     5 * time.Second,
		WarmupDelay:        time.Second,
		MaxBufferSize:      10 * 1024 * 1024,
		MaxPendingMessages: 1000,
		NormalizeURIs:      true,
		PatchOpenCode:      true,
		Debug:              os.Getenv("GODOT_LSP_BRIDGE_DEBUG") == "true",
		LogFile:            os.Getenv("GODOT_LSP_BRIDGE_LOG"),
		GodotPath:          os.Getenv("GODOT_PATH"),
		ProjectPath:        firstNonEmpty(os.Getenv("GODOT_PROJECT"), os.Getenv("OPENCODE_PROJECT_ROOT")),
	}

	fs := flag.NewFlagSet("godot-lsp-go", flag.ContinueOnError)
	ports := intsToCSV(cfg.Ports)
	var pathMaps repeatedString
	if envMap := os.Getenv("GODOT_LSP_PATH_MAP"); envMap != "" {
		_ = pathMaps.Set(envMap)
	}
	fs.StringVar(&cfg.Host, "host", cfg.Host, "Godot LSP host")
	fs.String("port", "", "fixed Godot LSP port")
	fs.StringVar(&ports, "ports", ports, "comma-separated Godot LSP ports")
	fs.BoolVar(&cfg.Launch, "launch", cfg.Launch, "launch Godot if LSP is unavailable")
	fs.StringVar(&cfg.GodotPath, "godot", cfg.GodotPath, "path to Godot executable")
	fs.StringVar(&cfg.ProjectPath, "project", cfg.ProjectPath, "path to Godot project")
	fs.BoolVar(&cfg.Headless, "headless", cfg.Headless, "launch Godot headless")
	noHeadless := fs.Bool("no-headless", false, "launch Godot with visible window")
	fs.StringVar(&cfg.XVFB, "xvfb", cfg.XVFB, "xvfb mode: auto, always, never")
	fs.StringVar(&cfg.CleanupGodot, "cleanup-godot", cfg.CleanupGodot, "cleanup launched Godot: auto, always, never")
	fs.BoolVar(&cfg.Reconnect, "reconnect", cfg.Reconnect, "reconnect when TCP connection closes")
	noReconnect := fs.Bool("no-reconnect", false, "disable reconnect")
	fs.DurationVar(&cfg.ConnectTimeout, "connect-timeout", cfg.ConnectTimeout, "TCP connection timeout")
	fs.DurationVar(&cfg.StartupTimeout, "startup-timeout", cfg.StartupTimeout, "Godot startup timeout")
	fs.DurationVar(&cfg.ReconnectDelay, "reconnect-delay", cfg.ReconnectDelay, "delay between reconnect attempts")
	fs.DurationVar(&cfg.WarmupDelay, "warmup-delay", cfg.WarmupDelay, "delay after reconnect")
	fs.IntVar(&cfg.MaxBufferSize, "max-buffer-size", cfg.MaxBufferSize, "maximum message size in bytes")
	fs.IntVar(&cfg.MaxPendingMessages, "max-pending-messages", cfg.MaxPendingMessages, "maximum pending messages")
	fs.BoolVar(&cfg.NormalizeURIs, "normalize-uris", cfg.NormalizeURIs, "normalize file URIs")
	noNormalize := fs.Bool("no-normalize-uris", false, "disable URI normalization")
	fs.BoolVar(&cfg.PatchOpenCode, "patch-opencode", cfg.PatchOpenCode, "patch OpenCode plaintext languageId")
	noPatch := fs.Bool("no-patch-opencode", false, "disable OpenCode patch")
	fs.Var(&pathMaps, "path-map", "path mapping in client-root=godot-root format; repeatable")
	fs.BoolVar(&cfg.Debug, "debug", cfg.Debug, "enable debug logs")
	fs.BoolVar(&cfg.Quiet, "quiet", cfg.Quiet, "reduce logs")
	fs.StringVar(&cfg.LogFile, "log-file", cfg.LogFile, "log file path")
	fs.BoolVar(&cfg.ShowVersion, "version", false, "show version")

	if err := fs.Parse(args); err != nil {
		return cfg, err
	}
	if p := fs.Lookup("port").Value.String(); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return cfg, err
		}
		cfg.Ports = []int{n}
	} else {
		parsed, err := parsePorts(ports)
		if err != nil {
			return cfg, err
		}
		cfg.Ports = parsed
	}
	if *noHeadless {
		cfg.Headless = false
	}
	if *noReconnect {
		cfg.Reconnect = false
	}
	if *noNormalize {
		cfg.NormalizeURIs = false
	}
	if *noPatch {
		cfg.PatchOpenCode = false
	}
	maps, err := parsePathMaps(pathMaps)
	if err != nil {
		return cfg, err
	}
	cfg.PathMaps = maps
	if cfg.ShowVersion {
		fmt.Printf("godot-lsp-go %s (%s, %s)\n", version.Version, version.Commit, version.Date)
	}
	return cfg, nil
}

func envString(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envPort(key string, fallback []int) []int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return []int{n}
		}
	}
	return fallback
}

func envPorts(key string, fallback []int) []int {
	if v := os.Getenv(key); v != "" {
		if p, err := parsePorts(v); err == nil {
			return p
		}
	}
	return fallback
}

func parsePorts(s string) ([]int, error) {
	parts := strings.Split(s, ",")
	ports := make([]int, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ports = append(ports, n)
	}
	if len(ports) == 0 {
		return nil, fmt.Errorf("no ports configured")
	}
	return ports, nil
}

func intsToCSV(values []int) string {
	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, strconv.Itoa(v))
	}
	return strings.Join(parts, ",")
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

type repeatedString []string

func (r *repeatedString) String() string {
	return strings.Join(*r, ";")
}

func (r *repeatedString) Set(value string) error {
	*r = append(*r, value)
	return nil
}

func parsePathMaps(values []string) ([]rewriter.PathMap, error) {
	var maps []rewriter.PathMap
	for _, value := range values {
		for _, part := range strings.Split(value, ";") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			m, err := rewriter.ParsePathMap(part)
			if err != nil {
				return nil, err
			}
			maps = append(maps, m)
		}
	}
	return rewriter.NormalizePathMaps(maps), nil
}
