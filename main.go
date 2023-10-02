package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"time"

	"github.com/Masterminds/sprig"
	"github.com/joho/godotenv"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

func check(e error) {
	if e != nil {
		log.Fatal(e.Error())
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var dbUrl = "libsql://dayssincerustdrama-defenestrator.turso.io?auth_token=" + os.Getenv("TURSO_API_TOKEN")
	db, error := sql.Open("libsql", dbUrl)
	if error != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, error)
		os.Exit(1)
	} else {
		// updateCount(db)
		tpl := template.Must(template.New("base").Funcs(sprig.FuncMap()).ParseGlob("*.html"))
		t, err := tpl.ParseFiles("index.html")
		check(err)
		fmt.Println(t)
		// To do :: figure out naked returns in Go. What a a newbarino.
		defer db.Close()
	}
}

func createDramaTable(db *sql.DB) {

	db.Exec("DROP TABLE drama;")

	createTable := `CREATE TABLE IF NOT EXISTS drama (
		"id" 			INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"description" 	TEXT,
		"url" 			TEXT,
		"date" 			TEXT NOT NULL
	  );`
	log.Println("Creating drama table...")
	dramaReportsTable, err := db.Prepare(createTable)
	check(err)
	dramaReportsTable.Exec()
	log.Println("drama reports table created")
}

func createDaysTable(db *sql.DB) {
	log.Println("Creating days table...")
	createTable := `CREATE TABLE IF NOT EXISTS days (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"days" BIGINT UNSIGNED NOT NULL
	  );`
	daysSinceDrama, err := db.Prepare(createTable)
	check(err)
	daysSinceDrama.Exec()
	log.Println("days table created")
}

func initializeDays(db *sql.DB) {
	zero := `INSERT INTO days ("days") VALUES (0);`
	initialization, err := db.Prepare(zero)
	check(err)
	initialization.Exec()
	log.Println("days table initialized")
}

func initializeDrama(db *sql.DB) {
	zero := `INSERT INTO drama ('drama', 'description', 'date', 'url') VALUES (TRUE, 'Some kind of bigotry accusations over nothing, probably, real dumb stuff', '2023-10-01', 'https://dayssincerustdrama.com');`
	initialization, err := db.Prepare(zero)
	check(err)
	initialization.Exec()
	log.Println("drama table initialized")
}

func getDaysCount(db *sql.DB) {
	var id uint64
	var days uint64
	err := db.QueryRow("SELECT * FROM days WHERE id = ?;", 1).Scan(&id, &days)
	check(err)
}

func getDramaCount(db *sql.DB) {
	var count uint64
	err := db.QueryRow("SELECT COUNT(*) FROM 'drama';").Scan(&count)
	check(err)
	fmt.Println(count)
}

func updateCount(db *sql.DB) {
	_, err := db.Exec("UPDATE 'days' SET days = days +1 WHERE id = 1;")
	check(err)
}

func reportDrama(db *sql.DB, post any) {
	datestamp := time.Now().Format(time.DateOnly)
	_, err := db.Exec(`INSERT INTO drama ('description', 'url', 'date') VALUES (TRUE, 'Some kind of bigotry accusations over nothing, probably, real dumb stuff', 'https://dayssincerustdrama.com', ?");`, datestamp)
	check(err)
}
