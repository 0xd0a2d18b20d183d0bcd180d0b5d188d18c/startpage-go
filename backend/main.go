package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type Item struct {
	Id       int    `db:"id"`
	URL      string `db:"url"`
	Shortcut string `db:"shortcut"`
	Desc     string `db:"desc"`
	UUID     string `db:"uuid"`
}

type Wrapper struct {
	Value []Item
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	page, err := os.ReadFile("frontend/index.html")
	check(err)
	io.WriteString(w, string(page))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	db := getDb()
	switch r.Method {
	case "GET":
		items := []Item{}
		db.Select(&items, "SELECT * FROM items i WHERE i.uuid='"+r.Header.Get("X-User-UUID")+"'")
		wrapper := &Wrapper{items}
		b, err := json.Marshal(wrapper)
		check(err)
		io.WriteString(w, string(b))
	case "POST":
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		check(err)
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO items (url, shortcut, desc, uuid) VALUES ($1, $2, $3, $4)", item.URL, item.Shortcut, item.Desc, r.Header.Get("X-User-UUID"))
		tx.Commit()
	case "DELETE":
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		check(err)
		tx := db.MustBegin()
		tx.MustExec("DELETE FROM items WHERE shortcut = $1 and uuid = $2", item.Shortcut, r.Header.Get("X-User-UUID"))
		tx.Commit()
	default:
		fmt.Fprint(w, "Sorry, wrong http request type")
	}
	db.Close()
}

func main() {
	fmt.Println("It works on http://localhost:3333")
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/items", getItems)
	err := http.ListenAndServe(":3333", mux)
	check(err)
}

func getDb() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "database/database.db")
	check(err)
	return db
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
