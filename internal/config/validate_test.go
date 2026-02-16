package config

import "testing"

func TestValidateValid(t *testing.T) {
	cfg := &Config{
		Storage:    StorageConfig{Path: "test.db"},
		Copilot:    CopilotConfig{Model: "gpt-4.1"},
		Categories: []CategoryConfig{{ID: "test", Name: "Test"}},
	}
	if err := Validate(cfg); err != nil {
		t.Errorf("valid config should not error: %v", err)
	}
}

func TestValidateEmptyPath(t *testing.T) {
	cfg := &Config{
		Storage:    StorageConfig{Path: ""},
		Copilot:    CopilotConfig{Model: "gpt-4.1"},
		Categories: []CategoryConfig{{ID: "test", Name: "Test"}},
	}
	if err := Validate(cfg); err == nil {
		t.Error("empty path should error")
	}
}
