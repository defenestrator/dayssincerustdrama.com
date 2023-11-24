package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
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
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s \n", dbUrl, error)
		os.Exit(1)
	} else {

		reportDrama(db)
		updateCount(db)
		count := strconv.Itoa(getDaysCount(db))
		buildPage(count)
		fmt.Printf("Daily update completed\n%s day(s) since drama in the Rust community.\n", count)
		defer db.Close()
	}
}

func buildPage(count string) (html string) {
	var processed bytes.Buffer

	tpl := template.Must(template.New("base.html").Funcs(sprig.FuncMap()).ParseGlob("*.html"))

	vars := map[string]interface{}{"DaysSinceDrama": count}

	err := tpl.ExecuteTemplate(&processed, "base.html", vars)

	if err != nil {
		fmt.Printf("Error during template execution: %s \n", err)
	}

	publicPath := "./public/index.html"
	p, _ := os.Create(publicPath)
	x := bufio.NewWriter(p)
	x.WriteString(string(processed.Bytes()))
	x.Flush()
	return
}

func getDaysCount(db *sql.DB) (count int) {
	err := db.QueryRow("SELECT days FROM days WHERE id = (SELECT MAX(id) FROM days);").Scan(&count)
	check(err)
	fmt.Println(count)
	return count
}

func getDramaCount(db *sql.DB) (count int) {
	err := db.QueryRow("SELECT COUNT(*) FROM 'drama';").Scan(&count)
	check(err)
	return
}

func updateCount(db *sql.DB) {
	var drama string

	err := db.QueryRow("SELECT date from drama WHERE id = (SELECT MAX(id) FROM drama);").Scan(&drama)

	check(err)

	now := time.Now()

	latestDrama, err := time.Parse(time.RFC3339, drama)

	check(err)

	days := strconv.Itoa(int(now.Sub(latestDrama).Hours() / 24))

	fmt.Printf("Days since drama %s\n", days)

	updateQuery := fmt.Sprintf("UPDATE 'days' SET days = %s WHERE id = (SELECT MAX(id) FROM days);", days)

	if int(now.Sub(latestDrama).Hours()/24) > 0 {
		_, err := db.Exec(updateQuery)
		check(err)
	} else {
		_, err := db.Exec(updateQuery)
		check(err)
		fmt.Printf("The count should be 0, unless our math is fucked up. Actual results: %s\n", days)
	}
}

func reportDrama(db *sql.DB) {
	datestamp := time.Now().Format(time.RFC3339)
	_, err := db.Exec(`INSERT INTO drama ('description', 'url', 'date') VALUES ('Some kind of accusations over nothing', 'https://dayssincerustdrama.com', ?);`, datestamp)
	check(err)
	updateCount(db)
	count := strconv.Itoa(getDaysCount(db))
	buildPage(count)
	fmt.Printf("Drama reported, days set to 0 on %s", datestamp)
}
