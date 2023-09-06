package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func main() {
	var dbUrl = "http://127.0.0.1:8080"
	db, error := sql.Open("libsql", dbUrl)
	if error != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, error)
		os.Exit(1)
	} else {
		defer db.Close() // Defer Closing the database
		createTables(db)

	}
}

func createTables(db *sql.DB) {
	createDramaTableSQL := `CREATE TABLE IF NOT EXISTS drama (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"drama" BOOLEAN,
		"description" TEXT,
		"date" 
		"url" TEXT		
	  );`
	log.Println("Creating drama table...")
	statement, err := db.Prepare(createDramaTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("drama table created")
}

func getCount() {
	// @TODO
}

func updateCount() {
	/// @TODO
}
