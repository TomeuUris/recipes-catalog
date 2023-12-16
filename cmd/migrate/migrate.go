package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TomeuUris/recipes-catalog/pkg/sqlite"
	"github.com/TomeuUris/recipes-catalog/pkg/sqlite/migrations"
)

func main() {
	db := sqlite.MustOpen("database.sqlite")

	if err := migrations.Store.Apply(context.Background(), db); err != nil {
		log.Panicf("error apply database migrations: %v", err.Error())
	}

	_ = db.Close()
	fmt.Println("Migrations successfully applied")
}
