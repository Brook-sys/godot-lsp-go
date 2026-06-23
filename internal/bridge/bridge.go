package bridge

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Brook-sys/godot-lsp-go/internal/config"
	"github.com/Brook-sys/godot-lsp-go/internal/connector"
	"github.com/Brook-sys/godot-lsp-go/internal/logging"
	"github.com/Brook-sys/godot-lsp-go/internal/lsp"
	"github.com/Brook-sys/godot-lsp-go/internal/rewriter"
	"github.com/Brook-sys/godot-lsp-go/internal/session"
)

type Bridge struct {
	cfg    config.Config
	log    *logging.Logger
	guard  session.Guard
	queue  *Queue
	conn   net.Conn
	connMu sync.RWMutex
	writer *lsp.Writer
}

func New(cfg config.Config, log *logging.Logger) *Bridge {
	return &Bridge{cfg: cfg, log: log, queue: NewQueue(cfg.MaxPendingMessages), writer: lsp.NewWriter(os.Stdout)}
}

func (b *Bridge) Run(ctx context.Context) error {
	if err := b.connect(ctx); err != nil {
		b.log.Warn("initial connection failed: %v", err)
	}
	go b.readStdin(ctx)
	for {
		select {
		case <-ctx.Done():
			b.closeConn()
			return nil
		default:
		}
		conn := b.getConn()
		if conn == nil {
			if !b.cfg.Reconnect {
				return fmt.Errorf("not connected")
			}
			if err := b.reconnect(ctx); err != nil {
				return err
			}
			continue
		}
		b.readTCP(ctx, conn)
		b.closeConn()
		if !b.cfg.Reconnect {
			return nil
		}
	}
}

func (b *Bridge) connect(ctx context.Context) error {
	res, err := connector.Connect(ctx, b.cfg.Host, b.cfg.Ports, b.cfg.ConnectTimeout)
	if err != nil {
		return err
	}
	b.setConn(res.Conn)
	b.log.Info("connected to Godot LSP on %s:%d", b.cfg.Host, res.Port)
	for _, msg := range b.queue.Drain() {
		_ = b.writeTCP(msg)
	}
	return nil
}

func (b *Bridge) reconnect(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(b.cfg.ReconnectDelay):
		}
		b.log.Info("trying to reconnect to Godot LSP")
		if err := b.connect(ctx); err != nil {
			b.log.Warn("reconnect failed: %v", err)
			continue
		}
		if b.cfg.WarmupDelay > 0 {
			time.Sleep(b.cfg.WarmupDelay)
		}
		_ = b.writer.WriteMessage(showMessage("Godot LSP server restarted. You may need to reopen files for diagnostics."))
		return nil
	}
}

func (b *Bridge) readStdin(ctx context.Context) {
	r := lsp.NewReader(os.Stdin)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, err := r.ReadMessage()
		if err != nil {
			if err != io.EOF {
				b.log.Warn("stdin read error: %v", err)
			}
			return
		}
		b.guard.TrackClientMessage(msg)
		msg.Body = rewriter.Rewrite(msg.Body, rewriter.Options{NormalizeURIs: b.cfg.NormalizeURIs, PatchOpenCode: b.cfg.PatchOpenCode, PathMaps: b.cfg.PathMaps, Direction: rewriter.ClientToGodot})
		if b.getConn() == nil {
			b.queue.Push(msg)
			continue
		}
		if err := b.writeTCP(msg); err != nil {
			b.log.Warn("tcp write failed: %v", err)
			b.queue.Push(msg)
			b.closeConn()
		}
	}
}

func (b *Bridge) readTCP(ctx context.Context, conn net.Conn) {
	r := lsp.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		msg, err := r.ReadMessage()
		if err != nil {
			if err != io.EOF {
				b.log.Warn("tcp read error: %v", err)
			}
			return
		}
		for _, out := range b.guard.HandleServerMessage(msg) {
			out.Body = rewriter.Rewrite(out.Body, rewriter.Options{NormalizeURIs: b.cfg.NormalizeURIs, PathMaps: b.cfg.PathMaps, Direction: rewriter.GodotToClient})
			if err := b.writer.WriteMessage(out); err != nil {
				b.log.Warn("stdout write failed: %v", err)
				return
			}
		}
	}
}

func (b *Bridge) writeTCP(msg lsp.Message) error {
	conn := b.getConn()
	if conn == nil {
		return fmt.Errorf("not connected")
	}
	return lsp.NewWriter(conn).WriteMessage(msg)
}

func (b *Bridge) getConn() net.Conn {
	b.connMu.RLock()
	defer b.connMu.RUnlock()
	return b.conn
}

func (b *Bridge) setConn(conn net.Conn) {
	b.connMu.Lock()
	defer b.connMu.Unlock()
	b.conn = conn
}

func (b *Bridge) closeConn() {
	b.connMu.Lock()
	defer b.connMu.Unlock()
	if b.conn != nil {
		_ = b.conn.Close()
		b.conn = nil
	}
	b.guard.ResetConnection()
}

func showMessage(message string) lsp.Message {
	msg, _ := lsp.NewMessage(map[string]any{
		"jsonrpc": "2.0",
		"method":  "window/showMessage",
		"params": map[string]any{
			"type":    2,
			"message": message,
		},
	})
	return msg
}
