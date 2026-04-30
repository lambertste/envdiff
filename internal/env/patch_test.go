package env

import (
	"testing"
)

func basePatchSet() *Set {
	s := NewSet()
	s.Set("APP_HOST", "localhost")
	s.Set("APP_PORT", "8080")
	s.Set("APP_ENV", "staging")
	return s
}

func TestPatch_SetNewKey(t *testing.T) {
	s := basePatchSet()
	out, err := Patch(s, []PatchInstruction{{Op: PatchSet, Key: "APP_DEBUG", Value: "true"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := out.Get("APP_DEBUG")
	if !ok || v != "true" {
		t.Errorf("expected APP_DEBUG=true, got %q (ok=%v)", v, ok)
	}
}

func TestPatch_SetOverwritesExisting(t *testing.T) {
	s := basePatchSet()
	out, err := Patch(s, []PatchInstruction{{Op: PatchSet, Key: "APP_PORT", Value: "9090"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, _ := out.Get("APP_PORT")
	if v != "9090" {
		t.Errorf("expected 9090, got %q", v)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	s := basePatchSet()
	out, err := Patch(s, []PatchInstruction{{Op: PatchDelete, Key: "APP_ENV"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := out.Get("APP_ENV")
	if ok {
		t.Error("expected APP_ENV to be deleted")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	s := basePatchSet()
	out, err := Patch(s, []PatchInstruction{{Op: PatchRename, Key: "APP_HOST", NewKey: "SERVICE_HOST"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, oldOk := out.Get("APP_HOST")
	v, newOk := out.Get("SERVICE_HOST")
	if oldOk {
		t.Error("old key APP_HOST should not exist after rename")
	}
	if !newOk || v != "localhost" {
		t.Errorf("expected SERVICE_HOST=localhost, got %q (ok=%v)", v, newOk)
	}
}

func TestPatch_DoesNotMutateOriginal(t *testing.T) {
	s := basePatchSet()
	_, err := Patch(s, []PatchInstruction{
		{Op: PatchSet, Key: "NEW_KEY", Value: "new"},
		{Op: PatchDelete, Key: "APP_PORT"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.Get("NEW_KEY"); ok {
		t.Error("original set should not have NEW_KEY")
	}
	if _, ok := s.Get("APP_PORT"); !ok {
		t.Error("original set should still have APP_PORT")
	}
}

func TestPatch_RenameNonExistentKey_ReturnsError(t *testing.T) {
	s := basePatchSet()
	_, err := Patch(s, []PatchInstruction{{Op: PatchRename, Key: "MISSING", NewKey: "OTHER"}})
	if err == nil {
		t.Error("expected error for renaming non-existent key")
	}
}

func TestPatch_UnknownOp_ReturnsError(t *testing.T) {
	s := basePatchSet()
	_, err := Patch(s, []PatchInstruction{{Op: "upsert", Key: "X"}})
	if err == nil {
		t.Error("expected error for unknown op")
	}
}
