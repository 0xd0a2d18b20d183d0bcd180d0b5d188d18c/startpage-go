package main

import (
	// "database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	dat, err := os.ReadFile("startpage.html")
	check(err)
	io.WriteString(w, string(dat))
}
func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	// fmt.Println("FUCK U")
	// isDbExists()

	// mux := http.NewServeMux()
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)

	err := http.ListenAndServe(":3333", nil)
	fmt.Println(err)

	// if errors.Is(err, http.ErrServerClosed) {
	// 	fmt.Printf("server closed\n")
	// } else if err != nil {
	// 	fmt.Printf("error starting server: %s\n", err)
	// 	os.Exit(1)
	// }
}

// func createDb() {
// 	file, err := os.Create("database.db")
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	file.Close()
// }

// func isDbExists() {
// 	db, err := sql.Open("sqlite3", "database.db")
// 	if err != nil {
// 		createDb()
// 	}

// 	// _, err = db.Query("SELECT * FROM items LIMIT 1")

// 	// _, err = db.Exec(".tables")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// _, err = db.Exec("CREATE TABLE `items` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `shortcut` VARCHAR(64) NOT NULL, `url` VARCHAR(255) NOT NULL, `desc` VARCHAR(255) NULL)")

// }

// Check is db file there
// If yes, read and prepare
// If not, create and put examples
// CRUD data
// Authentication by uuid at url's parameters (for first iteration)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
