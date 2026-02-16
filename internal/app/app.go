package app

import (
	"github.com/aknow2/beholder/internal/classify"
	"github.com/aknow2/beholder/internal/config"
	"github.com/aknow2/beholder/internal/storage"
)

type App struct {
	Config     *config.Config
	Storage    *storage.Store
	Classifier *classify.Client
}

func NewApp(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, err
	}
	if err := config.Validate(cfg); err != nil {
		return nil, err
	}

	store, err := storage.Open(cfg.Storage.Path)
	if err != nil {
		return nil, err
	}

	if err := store.Migrate(); err != nil {
		_ = store.Close()
		return nil, err
	}

	// T023: Remove UpsertCategories call - categories now only in Config

	return &App{
		Config:     cfg,
		Storage:    store,
		Classifier: classify.NewClient(cfg.Copilot.Model),
	}, nil
}

func (a *App) Close() {
	if a != nil && a.Storage != nil {
		_ = a.Storage.Close()
	}
}
