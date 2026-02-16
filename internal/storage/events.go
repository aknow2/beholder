package storage

import (
	"database/sql"
	"encoding/json"
	"time"
)

func dateRangeUTC(date time.Time) (time.Time, time.Time) {
	loc := date.Location()
	if loc == nil {
		loc = time.Local
	}
	startLocal := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	endLocal := startLocal.Add(24 * time.Hour)
	return startLocal.UTC(), endLocal.UTC()
}

func (s *Store) InsertEvent(event *Event) error {
	appsJSON, _ := json.Marshal(event.DetectedApps)
	keywordsJSON, _ := json.Marshal(event.DetectedKeywords)

	_, err := s.DB.Exec(`INSERT INTO events (
		id, captured_at, category_name, confidence, status, agent_version, screenshot_hash, detected_apps, detected_keywords, notes, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		event.ID,
		event.CapturedAt.UTC().Format(time.RFC3339),
		event.CategoryName,
		event.Confidence,
		event.Status,
		event.AgentVersion,
		event.ScreenshotHash,
		string(appsJSON),
		string(keywordsJSON),
		event.Notes,
		event.CreatedAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (s *Store) ListEventsByDate(date time.Time) ([]Event, error) {
	start, end := dateRangeUTC(date)

	rows, err := s.DB.Query(`SELECT id, captured_at, category_name, confidence, status, agent_version, screenshot_hash, detected_apps, detected_keywords, notes, created_at
		FROM events WHERE captured_at >= ? AND captured_at < ? ORDER BY captured_at ASC`,
		start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Event
	for rows.Next() {
		var e Event
		var capturedAt string
		var createdAt string
		var detectedApps string
		var detectedKeywords string
		if err := rows.Scan(&e.ID, &capturedAt, &e.CategoryName, &e.Confidence, &e.Status, &e.AgentVersion, &e.ScreenshotHash, &detectedApps, &detectedKeywords, &e.Notes, &createdAt); err != nil {
			return nil, err
		}
		e.CapturedAt, _ = time.Parse(time.RFC3339, capturedAt)
		e.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		_ = json.Unmarshal([]byte(detectedApps), &e.DetectedApps)
		_ = json.Unmarshal([]byte(detectedKeywords), &e.DetectedKeywords)
		results = append(results, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (s *Store) DeleteEventsByDate(date time.Time) (int64, error) {
	start, end := dateRangeUTC(date)

	res, err := s.DB.Exec(`DELETE FROM events WHERE captured_at >= ? AND captured_at < ?`,
		start.Format(time.RFC3339), end.Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
}

var _ = sql.ErrNoRows
