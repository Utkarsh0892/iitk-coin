package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func db(rn int, coin int) {
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, coins INTEGER)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO user (rollno, coins) VALUES (?, ?)")
	statement.Exec(rn, coin)
}

func credit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	rollno := r.FormValue("rollno")
	coins := r.FormValue("coins")
	rn, _ := strconv.Atoi(rollno)
	c, _ := strconv.Atoi(coins)
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, coins INTEGER)")
	statement.Exec()
	rows, _ :=
		database.Query("SELECT rollno, coins FROM user")
	var rolln int
	var coin int
	flag := 0
	for rows.Next() {
		rows.Scan(&rolln, &coin)
		if rn == rolln {
			_, err := database.Exec("UPDATE user SET coins = coins + ? WHERE rollno = ?", c, rn)
			flag = 1
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if flag == 0 {
		db(rn, c)
	}
}

func transfer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	rollnof := r.FormValue("rollnofrom")
	rollnot := r.FormValue("rollnoto")
	amt := r.FormValue("coins")
	rnf, _ := strconv.Atoi(rollnof)
	rnt, _ := strconv.Atoi(rollnot)
	c, _ := strconv.Atoi(amt)
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	tx,err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.Exec("UPDATE user SET coins = coins - ? WHERE rollno = ? AND coins - ? >= 0", c, rnf, c)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error")
		return
	}
	_, err = database.Exec("UPDATE user SET coins = coins + ? WHERE rollno = ?", c, rnt)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error")
		return
	}
	tx.Commit()
	fmt.Fprintf(w, "Coin Transfer successful\n")
}

func balance(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("this route can only handle GET request")
		return
	} 
    rolln := r.FormValue("rollno")
	rn, _ := strconv.Atoi(rolln)
	fmt.Fprintf(w, "GET request successful\n")
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	rows, err := database.Query("SELECT rollno ,coins FROM user")
	if err != nil {
		panic(err)
	}
	var rollno int
	var coins int
	for rows.Next() {
		rows.Scan(&rollno, &coins)
		if rollno == rn {
			fmt.Fprintf(w, "you have %d coins", coins)
			return
		}
	}
	fmt.Fprintf(w, "Roll no not in database")
}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/credit", credit)
	http.HandleFunc("/transfer", transfer)
	http.HandleFunc("/balance", balance)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
