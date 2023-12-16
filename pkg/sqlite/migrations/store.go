package migrations

import "github.com/TomeuUris/recipes-catalog/pkg/sqlite"

var Store sqlite.MigrationStore

func init() {
	if err := Store.Load(sqlite.MigrationsFS); err != nil {
		panic(err)
	}
}