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
	dat, err := os.ReadFile("startpage.html")
	check(err)
	io.WriteString(w, string(dat))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		db := getDb()
		query := "SELECT * FROM items i WHERE i.uuid='" + r.Header.Get("X-User-UUID") + "'"

		item := []Item{}
		db.Select(&item, query)
		wrapper := &Wrapper{item}

		b, err := json.Marshal(wrapper)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(b))
		io.WriteString(w, string(b))
		db.Close()
	case "POST":
		db := getDb()
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		check(err)
		tx := db.MustBegin()
		tx.MustExec("INSERT INTO items (url, shortcut, desc, uuid) VALUES ($1, $2, $3, $4)", item.URL, item.Shortcut, item.Desc, r.Header.Get("X-User-UUID"))
		tx.Commit()

		db.Close()
	case "PUT":
		db := getDb()
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		check(err)
		tx := db.MustBegin()
		tx.MustExec("UPDATE items SET url = $1, desc = $2, shortcut = $3 WHERE id = $4 and uuid = $5", item.URL, item.Desc, item.Shortcut, item.Id, r.Header.Get("X-User-UUID"))
		tx.Commit()
		db.Close()
	case "DELETE":
		db := getDb()
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		check(err)
		tx := db.MustBegin()
		tx.MustExec("DELETE FROM items WHERE id = $1 and uuid = $2", item.Id, r.Header.Get("X-User-UUID"))
		tx.Commit()
		db.Close()
	default:
		fmt.Fprint(w, "Sorry, wrong http request type")
	}
}

func main() {
	fmt.Println("It's work on http://localhost:3333")
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/items", getItems)

	err := http.ListenAndServe(":3333", mux)
	check(err)
}

func getDb() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "database.db")
	check(err)
	return db
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
