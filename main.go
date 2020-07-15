package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type CallClass struct {
	Name        string
	Phone       string
	Description string
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conent-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "GET called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, err := ioutil.ReadAll(r.Body)
	fmt.Printf("%s", b)
	if err != nil {
		log.Fatal(err)
	} else {
		var message CallClass
		e := json.Unmarshal(b, &message)
		if e != nil {
			log.Fatal(e)
		} else {
			fmt.Println("---------------")
			fmt.Println(message)
			fmt.Println("---------------")

			db := connectToDB()
			var id = 0
			err = db.QueryRow(
				`INSERT INTO glosite.call_sub(name, phone, description) VALUES ($1, $2, $3)`,
				message.Name, message.Phone, message.Description).Scan(&id)
			if err != nil {
				fmt.Println(err)
			}

			w.Write([]byte(http.StatusText(http.StatusCreated)))
		}

	}

}

func connectToDB() *sql.DB {
	connStr := "user=glosite password=glosite dbname=glosite sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		return db
	}
}

func main() {
	db := connectToDB()
	rows, err := db.Query("SELECT * FROM glosite.call_sub")
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			var row CallClass
			var temp string
			if err := rows.Scan(&temp, &temp, &row.Name, &row.Phone, &row.Description); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				log.Fatal(err)
			}
			fmt.Println(row)
		}
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", get).Methods(http.MethodGet)
	api.HandleFunc("/", post).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
