package main

import (
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
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
	http.HandleFunc("/redeem", redeem)
	http.HandleFunc("/makeAdmin", makeAdmin)
	http.HandleFunc("/manageRedeem", manageRedeemRequests)
	http.HandleFunc("/viewRedeem", viewRedeemRequests)
	http.HandleFunc("/updateInfo", updateInfo)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
