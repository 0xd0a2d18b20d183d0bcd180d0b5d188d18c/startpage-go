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

// Check is db file there
// If yes, read and prepare
// If not, create and put examples
// CRUD data
// Authentication by uuid at url's parameters (for first iteration)

type Item struct {
	Id       int    `db:"id"`
	URL      string `db:"url"`
	Shortcut string `db:"shortcut"`
	Desc     string `db:"desc"`
}

type Wrapper struct {
	Value []Item
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	dat, err := os.ReadFile("startpage.html")
	check(err)
	fmt.Printf("That's root")
	io.WriteString(w, string(dat))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		db := getDb()

		item := []Item{}
		db.Select(&item, "SELECT * FROM items")
		// fmt.Printf("%#v\n", item)
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
		tx.MustExec("INSERT INTO items (url, shortcut, desc) VALUES ($1, $2, $3)", item.URL, item.Shortcut, item.Desc)
		tx.Commit()

		db.Close()
	case "PUT":
	case "DELETE":
	default:
		fmt.Fprint(w, "Sorry, wrong http request type")
	}
}

func main() {
	fmt.Printf("It's working")

	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/items", getItems)

	err := http.ListenAndServe(":3333", mux)
	check(err)
}

// func createDb() {
// 	_, err := os.Create("database.db")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	fmt.Printf("File created")
// }

// func isDbExists() {
// 	fmt.Printf("Here")
// 	db, err := sql.Open("sqlite3", "database.db")
// 	if err != nil {
// 		fmt.Printf("Miss file")
// 		panic(err)
// 	}

// 	db.Close()
// }

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
