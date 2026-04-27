package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/envdiff/internal/watch"
)

// runWatch watches one or more .env files and prints a line to stdout
// whenever a file's content changes. It exits cleanly on SIGINT/SIGTERM.
func runWatch(paths []string, intervalMs int) error {
	if len(paths) == 0 {
		return fmt.Errorf("watch: at least one file path is required")
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			return fmt.Errorf("watch: cannot access %q: %w", p, err)
		}
	}

	interval := time.Duration(intervalMs) * time.Millisecond
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}

	done := make(chan struct{})
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		close(done)
	}()

	ch, err := watch.Watch(paths, interval, done)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "watching %d file(s) every %v — press Ctrl+C to stop\n", len(paths), interval)
	for ev := range ch {
		fmt.Fprintf(os.Stdout, "changed  %s  %s -> %s  at %s\n",
			ev.Path, ev.OldHash[:8], ev.NewHash[:8],
			ev.At.Format(time.RFC3339),
		)
	}
	fmt.Fprintln(os.Stdout, "watch stopped")
	return nil
}
