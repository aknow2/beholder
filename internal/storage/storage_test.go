package storage

import (
	"path/filepath"
	"testing"
)

func TestOpen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "test.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	if s.DB == nil {
		t.Error("nil")
	}
}
