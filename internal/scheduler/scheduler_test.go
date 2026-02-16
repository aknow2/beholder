package scheduler

import (
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := New(10, func(ctx context.Context) error { return nil })
	if s == nil {
		t.Error("nil")
	}
}
func TestRun(t *testing.T) {
	c := 0
	s := New(1, func(ctx context.Context) error { c++; return nil })
	s.interval = 50 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()
	go s.Start(ctx)
	<-ctx.Done()
	if c < 2 {
		t.Error("not called")
	}
}
