package main

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func main() {
	OpenDB()
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/home", Home)
	http.HandleFunc("/credit", credit)
	http.HandleFunc("/transfer", transfer)
	http.HandleFunc("/balance", checkBalance)
	http.HandleFunc("/logout", logout)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
