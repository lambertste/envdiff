package audit

import (
	"strings"
	"testing"
)

func TestLog_EmptyFormat(t *testing.T) {
	l := &Log{}
	got := l.Format()
	if got != "no changes recorded" {
		t.Errorf("expected 'no changes recorded', got %q", got)
	}
}

func TestLog_RecordAdded(t *testing.T) {
	l := &Log{}
	l.Record("prod.env", "NEW_KEY", EventAdded, "", "value1")
	if len(l.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(l.Events))
	}
	e := l.Events[0]
	if e.Type != EventAdded || e.Key != "NEW_KEY" || e.NewValue != "value1" {
		t.Errorf("unexpected event: %+v", e)
	}
}

func TestLog_RecordRemoved(t *testing.T) {
	l := &Log{}
	l.Record("prod.env", "OLD_KEY", EventRemoved, "oldval", "")
	out := l.Format()
	if !strings.Contains(out, "- OLD_KEY") {
		t.Errorf("expected removed marker in output, got: %s", out)
	}
}

func TestLog_RecordModified(t *testing.T) {
	l := &Log{}
	l.Record("staging.env", "DB_HOST", EventModified, "localhost", "db.prod.internal")
	out := l.Format()
	if !strings.Contains(out, "~ DB_HOST") {
		t.Errorf("expected modified marker in output, got: %s", out)
	}
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "db.prod.internal") {
		t.Errorf("expected old and new values in output, got: %s", out)
	}
}

func TestLog_Summary(t *testing.T) {
	l := &Log{}
	l.Record("f", "A", EventAdded, "", "1")
	l.Record("f", "B", EventAdded, "", "2")
	l.Record("f", "C", EventRemoved, "3", "")
	l.Record("f", "D", EventModified, "x", "y")

	s := l.Summary()
	if s[EventAdded] != 2 {
		t.Errorf("expected 2 added, got %d", s[EventAdded])
	}
	if s[EventRemoved] != 1 {
		t.Errorf("expected 1 removed, got %d", s[EventRemoved])
	}
	if s[EventModified] != 1 {
		t.Errorf("expected 1 modified, got %d", s[EventModified])
	}
}

func TestLog_FormatContainsFile(t *testing.T) {
	l := &Log{}
	l.Record("myfile.env", "KEY", EventAdded, "", "val")
	out := l.Format()
	if !strings.Contains(out, "myfile.env") {
		t.Errorf("expected filename in output, got: %s", out)
	}
}
