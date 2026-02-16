package summary

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aknow2/beholder/internal/storage"
)

type CategorySummary struct {
	CategoryName string
	Count        int
	Events       []storage.Event
}

type DailySummary struct {
	Date       time.Time
	Categories []CategorySummary
	TotalCount int
	FirstAt    time.Time
	LastAt     time.Time
}

func Generate(events []storage.Event) *DailySummary {
	if len(events) == 0 {
		return &DailySummary{
			Date:       time.Now().Local(),
			Categories: []CategorySummary{},
			TotalCount: 0,
		}
	}

	firstAt := events[0].CapturedAt
	lastAt := events[0].CapturedAt

	categoryMap := make(map[string]*CategorySummary)

	for _, event := range events {
		if event.CapturedAt.Before(firstAt) {
			firstAt = event.CapturedAt
		}
		if event.CapturedAt.After(lastAt) {
			lastAt = event.CapturedAt
		}

		catName := event.CategoryName
		if catName == "" {
			catName = "未分類"
		}

		if _, exists := categoryMap[catName]; !exists {
			categoryMap[catName] = &CategorySummary{
				CategoryName: catName,
				Count:        0,
				Events:       []storage.Event{},
			}
		}

		categoryMap[catName].Count++
		categoryMap[catName].Events = append(categoryMap[catName].Events, event)
	}

	categorySummaries := make([]CategorySummary, 0, len(categoryMap))
	for _, cat := range categoryMap {
		categorySummaries = append(categorySummaries, *cat)
	}

	sort.Slice(categorySummaries, func(i, j int) bool {
		return categorySummaries[i].Count > categorySummaries[j].Count
	})

	return &DailySummary{
		Date:       events[0].CapturedAt.In(time.Local),
		Categories: categorySummaries,
		TotalCount: len(events),
		FirstAt:    firstAt.In(time.Local),
		LastAt:     lastAt.In(time.Local),
	}
}

func (s *DailySummary) FormatMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Daily Report - %s\n\n", s.Date.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("**Total Events**: %d\n\n", s.TotalCount))
	sb.WriteString(fmt.Sprintf("**First Classified**: %s\n", s.FirstAt.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("**Last Classified**: %s\n\n", s.LastAt.Format("15:04:05")))

	if len(s.Categories) == 0 {
		sb.WriteString("No events recorded.\n")
		return sb.String()
	}

	sb.WriteString("## Summary by Category\n\n")
	for _, cat := range s.Categories {
		sb.WriteString(fmt.Sprintf("### %s\n", cat.CategoryName))
		sb.WriteString(fmt.Sprintf("- Count: %d\n", cat.Count))
		sb.WriteString(fmt.Sprintf("- Percentage: %.1f%%\n\n", float64(cat.Count)/float64(s.TotalCount)*100))
	}

	sb.WriteString("## Timeline\n\n")
	allEvents := []storage.Event{}
	for _, cat := range s.Categories {
		allEvents = append(allEvents, cat.Events...)
	}

	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].CapturedAt.Before(allEvents[j].CapturedAt)
	})

	for _, event := range allEvents {
		catName := event.CategoryName
		if catName == "" {
			catName = "未分類"
		}
		sb.WriteString(fmt.Sprintf("- %s | **%s** | confidence: %.2f | status: %s\n",
			event.CapturedAt.In(time.Local).Format("15:04:05"),
			catName,
			event.Confidence,
			event.Status))
	}

	return sb.String()
}

func (s *DailySummary) FormatText() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Daily Report - %s\n", s.Date.Format("2006-01-02")))
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")
	sb.WriteString(fmt.Sprintf("Total Events: %d\n\n", s.TotalCount))
	sb.WriteString(fmt.Sprintf("First Classified: %s\n", s.FirstAt.Format("15:04:05")))
	sb.WriteString(fmt.Sprintf("Last Classified: %s\n\n", s.LastAt.Format("15:04:05")))

	if len(s.Categories) == 0 {
		sb.WriteString("No events recorded.\n")
		return sb.String()
	}

	sb.WriteString("Summary by Category:\n")
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	for _, cat := range s.Categories {
		sb.WriteString(fmt.Sprintf("%s: %d events (%.1f%%)\n",
			cat.CategoryName,
			cat.Count,
			float64(cat.Count)/float64(s.TotalCount)*100))
	}

	return sb.String()
}
