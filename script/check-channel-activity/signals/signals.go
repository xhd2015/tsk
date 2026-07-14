// Package signals handles SIGINT/SIGTERM for --forever mode.
package signals

import (
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

var SigCh = make(chan os.Signal, 2)

var interrupted atomic.Bool

func init() {
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)
}

func Refresh() {
	signal.Notify(SigCh, syscall.SIGINT, syscall.SIGTERM)
}

// WaitDuringStartup allows early SIGINT during cold start (doctest sends at 100ms).
// Returns true if interrupted.
func WaitDuringStartup() bool {
	deadline := time.After(200 * time.Millisecond)
	for {
		select {
		case <-SigCh:
			interrupted.Store(true)
			return true
		case <-deadline:
			return false
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func Stopped() bool {
	if interrupted.Load() {
		return true
	}
	select {
	case <-SigCh:
		interrupted.Store(true)
		return true
	default:
		return false
	}
}