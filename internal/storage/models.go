package storage

import "time"

type Category struct {
	ID          string
	Name        string
	Description string
	Examples    []string
	Color       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Event struct {
	ID               string
	CapturedAt       time.Time
	CategoryName     string
	Confidence       float64
	Status           string
	AgentVersion     string
	ScreenshotHash   string
	DetectedApps     []string
	DetectedKeywords []string
	Notes            string
	CreatedAt        time.Time
}
