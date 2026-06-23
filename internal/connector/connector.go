package connector

import (
	"context"
	"fmt"
	"net"
	"time"
)

type Result struct {
	Conn net.Conn
	Port int
}

func Connect(ctx context.Context, host string, ports []int, timeout time.Duration) (Result, error) {
	var lastErr error
	for _, port := range ports {
		d := net.Dialer{Timeout: timeout}
		conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", host, port))
		if err == nil {
			return Result{Conn: conn, Port: port}, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no ports configured")
	}
	return Result{}, lastErr
}

func IsPortOpen(ctx context.Context, host string, port int, timeout time.Duration) bool {
	res, err := Connect(ctx, host, []int{port}, timeout)
	if err != nil {
		return false
	}
	_ = res.Conn.Close()
	return true
}
