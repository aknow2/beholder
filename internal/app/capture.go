package app

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aknow2/beholder/internal/config"
)

type CaptureResult struct {
	PNG          []byte
	DisplayCount int
	Resolution   string
	ImagePath    string
	CleanupImage bool
}

func captureFullScreenPNG(cfg *config.Config) (*CaptureResult, error) {
	cleanupImage := false

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	imgDir := filepath.Join(homeDir, ".beholder", "imgs")
	if cfg.Image.SaveImages {
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			return nil, err
		}
	} else {
		imgDir = os.TempDir()
		cleanupImage = true
	}

	tmpDir := os.TempDir()
	rawPath := filepath.Join(tmpDir, fmt.Sprintf("beholder-raw-%d.png", time.Now().UnixNano()))
	timestamp := time.Now().Format("20060102-150405")

	// T013: Select format based on config
	ext := "jpg"
	formatArg := "jpeg"
	if cfg.Image.Format == "png" {
		ext = "png"
		formatArg = "png"
	}
	resizedPath := filepath.Join(imgDir, fmt.Sprintf("screenshot-%s.%s", timestamp, ext))

	cmd := exec.Command("screencapture", "-x", "-t", "png", rawPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	defer os.Remove(rawPath)

	// T012: Use Config.Image.MaxWidth dynamically
	maxWidthStr := fmt.Sprintf("%d", cfg.Image.MaxWidth)
	resizeCmd := exec.Command("sips", "-s", "format", formatArg, "-Z", maxWidthStr, rawPath, "--out", resizedPath)
	if err := resizeCmd.Run(); err != nil {
		return nil, fmt.Errorf("resize failed: %w", err)
	}

	info, err := os.Stat(resizedPath)
	if err != nil {
		return nil, err
	}
	if info.Size() > 3*1024*1024 {
		return nil, fmt.Errorf("image too large after resize: %d bytes", info.Size())
	}

	f, err := os.Open(resizedPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	cfg2, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	// T015-T017: Cleanup old images if max_files is set
	if cfg.Image.SaveImages && cfg.Image.MaxFiles > 0 {
		if err := cleanupOldImages(imgDir, cfg.Image.MaxFiles); err != nil {
			// T017: Log error but continue gracefully
			log.Printf("Warning: failed to cleanup old images: %v", err)
		}
	}

	return &CaptureResult{
		PNG:          data,
		DisplayCount: 1,
		Resolution:   fmt.Sprintf("%dx%d", cfg2.Width, cfg2.Height),
		ImagePath:    resizedPath,
		CleanupImage: cleanupImage,
	}, nil
}

// T015-T016: Cleanup old images based on max_files setting
func cleanupOldImages(imgDir string, maxFiles int) error {
	files, err := os.ReadDir(imgDir)
	if err != nil {
		return err
	}

	// Filter image files and sort by name (timestamp-based)
	var imageFiles []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasPrefix(name, "screenshot-") && (strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpeg")) {
			imageFiles = append(imageFiles, name)
		}
	}

	// T016: Sort by timestamp (oldest first)
	sort.Strings(imageFiles)

	// Delete oldest files if count exceeds maxFiles
	if len(imageFiles) > maxFiles {
		toDelete := len(imageFiles) - maxFiles
		for i := 0; i < toDelete; i++ {
			filePath := filepath.Join(imgDir, imageFiles[i])
			if err := os.Remove(filePath); err != nil {
				// T017: Log warning but continue
				log.Printf("Warning: failed to delete old image %s: %v", imageFiles[i], err)
			}
		}
	}

	return nil
}
