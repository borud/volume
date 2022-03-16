package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	httpListenAddr string
	dbFilename     string
)

func init() {
	flag.StringVar(&httpListenAddr, "w", ":8080", "HTTP listen address")
	flag.StringVar(&dbFilename, "d", "app.db", "database file path")
	flag.Parse()
}

func main() {
	db := openDB(dbFilename)
	defer db.Close()

	// insert one value every 5 seconds
	go func() {
		for {
			db.Exec("INSERT INTO data (num) VALUES(?)", time.Now().UnixMilli())
			time.Sleep(5 * time.Second)
		}
	}()

	// quick and dirty webserver
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var nums []int
		db.Select(&nums, "SELECT num FROM data LIMIT 100")

		err := json.NewEncoder(w).Encode(nums)
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshalling JSON: %v", err), http.StatusInternalServerError)
		}
	})

	log.Printf("listening on [%s]", httpListenAddr)
	log.Printf("%v", http.ListenAndServe(httpListenAddr, nil))
}

func openDB(filename string) *sqlx.DB {
	// Make a note of whether file existed
	_, err := os.Stat(dbFilename)
	dbNeedsCreation := os.IsNotExist(err)

	db, err := sqlx.Open("sqlite3", filename)
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to open ping: %v", err)
	}

	if dbNeedsCreation {
		log.Printf("database [%s] didn't exist, creating", filename)
		_, err := db.Exec("CREATE TABLE data (num INT);")
		if err != nil {
			log.Fatalf("Unable to create table [data]: %v", err)
		}
	}

	return db
}
