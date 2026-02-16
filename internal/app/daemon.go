package app

import (
	"context"
	"log"

	"github.com/aknow2/beholder/internal/scheduler"
)

func (a *App) StartScheduler(ctx context.Context) error {
	if a.Config.Scheduler.IntervalMinutes <= 0 {
		log.Println("scheduler interval not configured, using default 10 minutes")
		a.Config.Scheduler.IntervalMinutes = 10
	}

	recordFunc := func(ctx context.Context) error {
		_, err := a.RecordOnce(ctx)
		return err
	}

	s := scheduler.New(a.Config.Scheduler.IntervalMinutes, recordFunc)

	log.Printf("starting scheduler with %d minute interval", a.Config.Scheduler.IntervalMinutes)
	s.Start(ctx)

	return nil
}
