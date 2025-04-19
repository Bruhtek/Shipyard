package database

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"os"
)

var DB *sql.DB = initializeDatabase()

func initializeDatabase() *sql.DB {
	DB, err := sql.Open("sqlite", "shipyard.db")
	if err != nil {
		panic(err)
	}

	// run migrations - each file from database/migrations, one by one
	files, err := os.ReadDir("database/migrations")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := "database/migrations/" + file.Name()
		migration, err := os.ReadFile(filePath)
		if err != nil {
			panic(err)
		}
		log.Println("Running migration: ", filePath)
		_, err = DB.Exec(string(migration))
		if err != nil {
			panic(err)
		}
	}

	return DB
}
