package classify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aknow2/beholder/internal/config"
	copilot "github.com/github/copilot-sdk/go"
)

type Result struct {
	SelectedCategoryID string   `json:"selectedCategoryId"`
	Confidence         float64  `json:"confidence"`
	Rationale          string   `json:"rationale"`
	DetectedApps       []string `json:"detectedApps,omitempty"`
	DetectedKeywords   []string `json:"detectedKeywords,omitempty"`
}

type Client struct {
	Model string
}

func NewClient(model string) *Client {
	return &Client{Model: model}
}

func (c *Client) Classify(ctx context.Context, imagePath string, categories []config.CategoryConfig) (*Result, error) {
	if imagePath == "" {
		return nil, fmt.Errorf("image path is empty")
	}
	if _, err := os.Stat(imagePath); err != nil {
		return nil, fmt.Errorf("image path is not accessible: %w", err)
	}

	client := copilot.NewClient(nil)
	if err := client.Start(); err != nil {
		return nil, err
	}
	defer client.Stop()

	session, err := client.CreateSession(&copilot.SessionConfig{Model: c.Model})
	if err != nil {
		return nil, err
	}
	defer session.Destroy()

	catsJSON, err := json.Marshal(categories)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`You are a screenshot classifier.
Use the attached image to classify the screenshot.
Return ONLY valid JSON with keys: selectedCategoryId, confidence, rationale, detectedApps, detectedKeywords.
Choose exactly one category id from the list.
Categories: %s
`, string(catsJSON))

	attachments := []copilot.Attachment{
		{
			DisplayName: filepath.Base(imagePath),
			Path:        imagePath,
			Type:        copilot.File,
		},
	}

	resp, err := session.SendAndWait(copilot.MessageOptions{Prompt: prompt, Attachments: attachments}, 0)
	if err != nil {
		return nil, err
	}

	if resp == nil || resp.Data.Content == nil {
		return nil, fmt.Errorf("empty response")
	}

	var result Result
	if err := json.Unmarshal([]byte(*resp.Data.Content), &result); err != nil {
		return nil, fmt.Errorf("invalid json response: %w", err)
	}

	return &result, nil
}
