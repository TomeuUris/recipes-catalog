package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path"
	"regexp"
	"sort"
	"strconv"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

type Migration struct {
	Version int
	Query   string
}

type MigrationStore struct {
	migrations []Migration
}

func (s *MigrationStore) Load(fs embed.FS) error {
	migrationsDir := "migrations"
	entries, err := fs.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		version := versionFromFilename(entry.Name())
		query, err := fs.ReadFile(path.Join(migrationsDir, entry.Name()))
		if err != nil {
			return err
		}

		s.migrations = append(s.migrations, Migration{
			Version: version,
			Query:   string(query),
		})
	}

	sort.Slice(s.migrations, func(i, j int) bool {
		return s.migrations[i].Version < s.migrations[j].Version
	})

	return nil
}

func (s *MigrationStore) Apply(ctx context.Context, db *sql.DB) error {
	var version int
	if err := db.QueryRow(`PRAGMA user_version`).Scan(&version); err != nil {
		return fmt.Errorf("failed to get user_version: %w", err)
	}

	log.Printf("Migration counter: %d/%d", version, len(s.migrations))

	for _, migration := range s.migrations {
		if migration.Version <= version {
			continue
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to create migration transaction: %w", err)
		}

		if _, err = tx.Exec(migration.Query); err != nil {
			return fmt.Errorf("failed to run migration: %w", err)
		}

		if _, err = tx.Exec(fmt.Sprintf(`PRAGMA user_version=%d`, migration.Version)); err != nil {
			return fmt.Errorf("failed to update version: %w", err)
		}

		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration: %w", err)
		}

		log.Printf("Migration counter: %d/%d", migration.Version, len(s.migrations))
	}

	return nil
}

func versionFromFilename(filename string) int {
	versionStr := regexp.MustCompile(`migration_(\d+)_.*\.sql`).FindStringSubmatch(filename)
	version, _ := strconv.ParseInt(versionStr[1], 10, 32)
	return int(version)
}
