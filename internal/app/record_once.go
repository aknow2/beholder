package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aknow2/beholder/internal/storage"
	"github.com/google/uuid"
)

func (a *App) RecordOnce(ctx context.Context) (*storage.Event, error) {
	captureResult, err := captureFullScreenPNG(a.Config)
	if err != nil {
		return nil, err
	}
	if captureResult.CleanupImage && captureResult.ImagePath != "" {
		defer func(path string) {
			_ = os.Remove(path)
		}(captureResult.ImagePath)
	}

	classification, err := a.Classifier.Classify(ctx, captureResult.ImagePath, a.Config.Categories)

	status := "OK"
	categoryID := ""
	categoryName := ""
	confidence := 0.0
	rationale := ""
	var detectedApps []string
	var detectedKeywords []string

	if err != nil {
		log.Printf("classification failed: %v", err)
		status = "FAILED"
	} else {
		categoryID = classification.SelectedCategoryID
		confidence = classification.Confidence
		rationale = classification.Rationale
		detectedApps = classification.DetectedApps
		detectedKeywords = classification.DetectedKeywords
	}

	if categoryID == "" && len(a.Config.Categories) > 0 {
		categoryID = a.Config.Categories[0].ID
	}

	// T020: Map category ID to Name from Config
	for _, cat := range a.Config.Categories {
		if cat.ID == categoryID {
			categoryName = cat.Name
			break
		}
	}

	hash := sha256.Sum256(captureResult.PNG)
	screenshotHash := hex.EncodeToString(hash[:])

	event := &storage.Event{
		ID:               uuid.NewString(),
		CapturedAt:       time.Now().UTC(),
		CategoryName:     categoryName,
		Confidence:       confidence,
		Status:           status,
		AgentVersion:     a.Config.Copilot.Model,
		ScreenshotHash:   screenshotHash,
		DetectedApps:     detectedApps,
		DetectedKeywords: detectedKeywords,
		Notes:            fmt.Sprintf("rationale=%s displayCount=%d resolution=%s", rationale, captureResult.DisplayCount, captureResult.Resolution),
		CreatedAt:        time.Now().UTC(),
	}

	if err := a.Storage.InsertEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}

func (a *App) ListEventsByDate(date time.Time) ([]storage.Event, error) {
	return a.Storage.ListEventsByDate(date)
}

func (a *App) DeleteEventsByDate(date time.Time) (int64, error) {
	return a.Storage.DeleteEventsByDate(date)
}
