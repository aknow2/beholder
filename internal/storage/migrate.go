package storage

import "database/sql"

func (s *Store) Migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS events (
			id TEXT PRIMARY KEY,
			captured_at TEXT NOT NULL,
			category_name TEXT,
			confidence REAL,
			status TEXT NOT NULL,
			agent_version TEXT,
			screenshot_hash TEXT,
			detected_apps TEXT,
			detected_keywords TEXT,
			notes TEXT,
			created_at TEXT NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := s.DB.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func withTx(db *sql.DB, fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
