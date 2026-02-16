package summary

import (
	"testing"
	"time"

	"github.com/aknow2/beholder/internal/storage"
)

func TestGen(t *testing.T) {
	s := Generate([]storage.Event{{ID: "1", CapturedAt: time.Now(), CategoryName: "T"}})
	if s.TotalCount != 1 {
		t.Error("fail")
	}
}
