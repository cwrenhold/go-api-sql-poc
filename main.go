package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cwrenhold/go-api-sql-poc/initializers"
	"github.com/jmoiron/sqlx"
)

type SqlScript int

const (
	select_tasks SqlScript = iota
	select_tasks_active
)

var sqlScripts = map[string]SqlScript{
	"select_tasks":        select_tasks,
	"select_tasks_active": select_tasks_active,
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/test-file", testFileHandler)
	http.HandleFunc("/test-file-script", testFileScriptHandler)

	port := os.Getenv("WEB_PORT")
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	initializers.ConnectToDB()

	query := "SELECT * FROM tasks"

	queryParams := r.URL.Query()

	var rows *sqlx.Rows
	var err error
	if queryParams["description"] != nil {
		query = query + " WHERE Description = :description"
		rows, err = initializers.SqlxDB.NamedQuery(query, map[string]interface{}{
			"description": queryParams["description"][0],
		})
	} else {
		rows, err = initializers.SqlxDB.Queryx(query)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	jsonData := readRowsToJsonArray(rows)

	w.Write(jsonData)
}

func testFileHandler(w http.ResponseWriter, r *http.Request) {
	initializers.ConnectToDB()

	data, err := os.ReadFile("scripts/select_tasks.sql")
	if err != nil {
		log.Fatal(err)
	}

	query := string(data)

	jsonData := runQueryToJson(query)

	w.Write(jsonData)
}

func testFileScriptHandler(w http.ResponseWriter, r *http.Request) {
	initializers.ConnectToDB()

	queryParams := r.URL.Query()

	rawScriptType := queryParams["script"][0]
	_, ok := sqlScripts[rawScriptType]

	if !ok {
		log.Fatal("Invalid script type")
	}

	data, err := os.ReadFile(fmt.Sprintf("scripts/%s.sql", rawScriptType))

	if err != nil {
		log.Fatal(err)
	}

	query := string(data)

	jsonData := runQueryToJson(query)

	w.Write(jsonData)
}

func runQueryToJson(query string) []byte {
	rows, err := initializers.SqlxDB.Queryx(query)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	return readRowsToJsonArray(rows)
}

func readRowsToJsonArray(rows *sqlx.Rows) []byte {
	var results []string

	for rows.Next() {
		result := make(map[string]interface{})
		err := rows.MapScan(result)
		if err != nil {
			log.Fatalln(err)
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatalln(err)
		}

		results = append(results, string(jsonData))
	}

	return []byte("[" + strings.Join(results, ",") + "]")
}
