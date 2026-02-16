package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type Store struct {
	DB *sql.DB
}

func Open(path string) (*Store, error) {
	// T010-T011: Path resolution logic
	resolvedPath := path

	// Expand tilde to home directory
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		resolvedPath = filepath.Join(home, path[1:])
	} else if !filepath.IsAbs(path) {
		// T010: Relative paths resolve to ~/.beholder/
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		resolvedPath = filepath.Join(home, ".beholder", path)
	}
	// T011: Absolute paths use as-is (no change needed)

	dir := filepath.Dir(resolvedPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", resolvedPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}
	return s.DB.Close()
}
