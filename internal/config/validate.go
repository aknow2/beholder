package config

import "fmt"

func Validate(cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if cfg.Storage.Path == "" {
		return fmt.Errorf("storage.path is required")
	}
	if cfg.Copilot.Model == "" {
		return fmt.Errorf("copilot.model is required")
	}
	if len(cfg.Categories) == 0 {
		return fmt.Errorf("categories must contain at least one entry")
	}

	// T005: Image settings validation
	if cfg.Image.MaxWidth < 100 || cfg.Image.MaxWidth > 4096 {
		return fmt.Errorf("image.max_width must be between 100 and 4096, got: %d", cfg.Image.MaxWidth)
	}
	if cfg.Image.MaxFiles < 0 {
		return fmt.Errorf("image.max_files must be >= 0, got: %d", cfg.Image.MaxFiles)
	}
	if cfg.Image.Format != "jpeg" && cfg.Image.Format != "png" {
		return fmt.Errorf("image.format must be 'jpeg' or 'png', got: %s", cfg.Image.Format)
	}

	ids := map[string]struct{}{}
	for _, c := range cfg.Categories {
		if c.ID == "" || c.Name == "" {
			return fmt.Errorf("category id and name are required")
		}
		if _, ok := ids[c.ID]; ok {
			return fmt.Errorf("duplicate category id: %s", c.ID)
		}
		ids[c.ID] = struct{}{}
	}

	return nil
}
