package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cwrenhold/go-api-sql-poc/initializers"
	"github.com/jmoiron/sqlx"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/test", testHandler)

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

	var results []string

	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.MapScan(result)
		if err != nil {
			log.Fatalln(err)
		}

		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatalln(err)
		}

		results = append(results, string(jsonData))
	}

	jsonData := []byte("[" + strings.Join(results, ",") + "]")

	w.Write(jsonData)
}
