package watch

import (
	"os"
	"testing"
	"time"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "envdiff-watch-*.env")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatch_DetectsChange(t *testing.T) {
	path := writeTmp(t, "KEY=old\n")
	done := make(chan struct{})
	defer close(done)

	ch, err := Watch([]string{path}, 20*time.Millisecond, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=new\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	select {
	case ev := <-ch:
		if ev.Path != path {
			t.Errorf("expected path %q, got %q", path, ev.Path)
		}
		if ev.OldHash == ev.NewHash {
			t.Error("expected hashes to differ")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	path := writeTmp(t, "KEY=stable\n")
	done := make(chan struct{})
	defer close(done)

	ch, err := Watch([]string{path}, 20*time.Millisecond, done)
	if err != nil {
		t.Fatalf("Watch: %v", err)
	}

	select {
	case ev := <-ch:
		t.Errorf("unexpected event: %+v", ev)
	case <-time.After(100 * time.Millisecond):
		// expected — no change
	}
}

func TestWatch_InvalidPathReturnsError(t *testing.T) {
	_, err := Watch([]string{"/nonexistent/path.env"}, 20*time.Millisecond, make(chan struct{}))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestChecksum_ConsistentForSameContent(t *testing.T) {
	path := writeTmp(t, "A=1\nB=2\n")
	s1, err := checksum(path)
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}
	s2, err := checksum(path)
	if err != nil {
		t.Fatalf("checksum: %v", err)
	}
	if s1.Checksum != s2.Checksum {
		t.Errorf("expected same checksum, got %q vs %q", s1.Checksum, s2.Checksum)
	}
}
