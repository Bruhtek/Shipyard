package database

import (
	"database/sql"
	"embed"

	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog/log"
)

var DB *sql.DB

//go:embed all:migrations/*.sql
var migrations embed.FS

func InitializeDatabase() *sql.DB {
	log.Info().
		Msg("Initializing database")

	var err error
	DB, err = sql.Open("sqlite", "shipyard.db")
	if err != nil {
		panic(err)
	}

	// run migrations - each file from database/migrations, one by one
	files, err := migrations.ReadDir("migrations")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := "migrations/" + file.Name()
		migration, err := migrations.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		log.Info().
			Str("file", filePath).
			Msg("Running database migration: ")
		_, err = DB.Exec(string(migration))
		if err != nil {
			panic(err)
		}
	}

	return DB
}
