package store

import (
	"context"
	"database/sql"
	"path/filepath"
	"time"

	"deal.digital/internal/content"

	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(dataDir string) (*SQLite, error) {
	dbPath := filepath.Join(dataDir, "site.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := configureSQLite(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := initSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &SQLite{db: db}, nil
}

func (s *SQLite) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func configureSQLite(db *sql.DB) error {
	pragmas := []string{
		`PRAGMA journal_mode=WAL;`,
		`PRAGMA synchronous=NORMAL;`,
		`PRAGMA foreign_keys=ON;`,
		`PRAGMA temp_store=MEMORY;`,
	}
	for _, statement := range pragmas {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}

func initSchema(db *sql.DB) error {
	schema := `
CREATE TABLE IF NOT EXISTS documents (
  route TEXT PRIMARY KEY,
  kind TEXT NOT NULL,
  source_path TEXT NOT NULL,
  checksum TEXT NOT NULL,
  title TEXT NOT NULL,
  updated_at TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS taxon_entries (
  taxonomy TEXT NOT NULL,
  term_slug TEXT NOT NULL,
  page_route TEXT NOT NULL,
  PRIMARY KEY (taxonomy, term_slug, page_route)
);
CREATE TABLE IF NOT EXISTS counters (
  name TEXT PRIMARY KEY,
  value INTEGER NOT NULL DEFAULT 0,
  updated_at TEXT NOT NULL
);`
	if _, err := db.Exec(schema); err != nil {
		return err
	}
	_, err := db.Exec(`INSERT INTO counters(name, value, updated_at) VALUES ('sample-page', 0, ?) ON CONFLICT(name) DO NOTHING`, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (s *SQLite) PersistSite(ctx context.Context, site *content.Site) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM documents`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM taxon_entries`); err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO documents(route, kind, source_path, checksum, title, updated_at) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, section := range site.Sections {
		if _, err := stmt.ExecContext(ctx, section.Route, "section", section.SourcePath, section.Checksum, section.Title, now); err != nil {
			return err
		}
	}
	for _, page := range site.Pages {
		if _, err := stmt.ExecContext(ctx, page.Route, "page", page.SourcePath, page.Checksum, page.Title, now); err != nil {
			return err
		}
	}
	for name, terms := range site.Taxonomies {
		for _, term := range terms {
			for _, page := range term.Pages {
				if _, err := tx.ExecContext(ctx, `INSERT INTO taxon_entries(taxonomy, term_slug, page_route) VALUES (?, ?, ?)`, name, term.Slug, page.Route); err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (s *SQLite) CounterValue(ctx context.Context, name string) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := s.db.ExecContext(ctx, `INSERT INTO counters(name, value, updated_at) VALUES (?, 0, ?) ON CONFLICT(name) DO NOTHING`, name, now); err != nil {
		return 0, err
	}

	var value int64
	if err := s.db.QueryRowContext(ctx, `SELECT value FROM counters WHERE name = ?`, name).Scan(&value); err != nil {
		return 0, err
	}
	return value, nil
}

func (s *SQLite) IncrementCounter(ctx context.Context, name string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := tx.ExecContext(ctx, `INSERT INTO counters(name, value, updated_at) VALUES (?, 0, ?) ON CONFLICT(name) DO NOTHING`, name, now); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE counters SET value = value + 1, updated_at = ? WHERE name = ?`, now, name); err != nil {
		return 0, err
	}

	var value int64
	if err := tx.QueryRowContext(ctx, `SELECT value FROM counters WHERE name = ?`, name).Scan(&value); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return value, nil
}
