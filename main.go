package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type CallApplicationClass struct {
	Id          int
	Time        time.Time
	Name        string
	Phone       string
	Description string
}

func connectToDB() *sql.DB {
	database, err := sql.Open("sqlite3", "./call.db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		return database
	}
}

func createTable(database *sql.DB) {
	statement, err :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, time TIMESTAMP, name TEXT, phone TEXT, description TEXT)")
	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		statement.Exec()
	}
}

func insertIntoDb(database *sql.DB, name string, phone string, description string) {
	statement, err :=
		database.Prepare("INSERT INTO people (time, name, phone, description) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	} else {
		statement.Exec(time.Now(), name, phone, description)
	}

}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Conent-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "GET called"}`))
}

func post(w http.ResponseWriter, r *http.Request) {
	writeCorsHeaders(&w)
	w.WriteHeader(http.StatusCreated)

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	} else {
		var message CallApplicationClass
		e := json.Unmarshal(b, &message)
		if e != nil {
			log.Fatal(e)
		} else {
			db := connectToDB()
			insertIntoDb(db, message.Name, message.Phone, message.Description)
			w.Write([]byte(http.StatusText(http.StatusCreated) + "-" + strconv.Itoa(http.StatusCreated)))
		}

	}

}

func writeCorsHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func main() {
	port := ":8080"
	if len(os.Args) > 1 && os.Args[1] != "" {
		port = os.Args[1]
	}
	fmt.Println("Server is listening on", port)
	database := connectToDB()
	createTable(database)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/", get).Methods(http.MethodGet)
	api.HandleFunc("/call", post).Methods(http.MethodPost)

	/*
		Only for developer use only
	*/
	corsHandled := cors.AllowAll().Handler(r)

	log.Fatal(http.ListenAndServe(port, corsHandled))
}
