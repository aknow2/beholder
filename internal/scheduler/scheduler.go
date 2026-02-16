package scheduler

import (
	"context"
	"log"
	"time"
)

type RecordFunc func(ctx context.Context) error

type Scheduler struct {
	interval   time.Duration
	recordFunc RecordFunc
	stopCh     chan struct{}
	doneCh     chan struct{}
}

func New(intervalMinutes int, recordFunc RecordFunc) *Scheduler {
	return &Scheduler{
		interval:   time.Duration(intervalMinutes) * time.Minute,
		recordFunc: recordFunc,
		stopCh:     make(chan struct{}),
		doneCh:     make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	defer close(s.doneCh)

	log.Printf("scheduler started with interval: %v", s.interval)

	for {
		select {
		case <-ticker.C:
			if err := s.recordFunc(ctx); err != nil {
				log.Printf("scheduled record failed: %v", err)
			}
		case <-s.stopCh:
			log.Println("scheduler stopped")
			return
		case <-ctx.Done():
			log.Println("scheduler context cancelled")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
	<-s.doneCh
}
