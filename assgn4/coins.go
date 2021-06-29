package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	/* jwt "github.com/dgrijalva/jwt-go" */
	_ "github.com/mattn/go-sqlite3"
)

/* var jwt_key = []byte("secret_key")

type Claims struct {
	Rollno string `json:"rollno"`
	jwt.StandardClaims
} */

func dbu(rn int, coin int) {
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, coins INTEGER)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO user (rollno, coins) VALUES (?, ?)")
	statement.Exec(rn, coin)
}

func dbt(aw int, frn int, trn int, coin int, tax int, ts string) {
	database, _ :=
		sql.Open("sqlite3", "./transactions.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS transactions (award INTEGER, fromrn INTEGER, torn INTEGER, coins INTEGER, tax INTEGER, timestamp TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO user (award, fromrn ,torn , coins, tax, timestamp) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(aw, frn, trn, coin, tax, ts)
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
		dbu(rn, c)
	}
	dbt(1, 0, rn, c, 0, time.Now().String())
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
	b1 := rnf
	b2 := rnt
	tax := 0
	for b1 >= 100 {
		b1 /= 10
	}
	for b2 >= 100 {
		b2 /= 10
	}
	if b1 == b2 {
		tax = (int)(2 * c / 100)
	} else {
		tax = (int)(33 * c / 100)
	}
	database, _ :=
		sql.Open("sqlite3", "./transactions.db")
	tx, err := database.Begin()
	if err != nil {
		log.Fatal(err)
	}
	count := 0
	var aw int
	var trn int
	rows, _ :=
		database.Query("SELECT award, ttorn FROM transactions")
	for rows.Next() {
		rows.Scan(&aw, &trn)
		if rnt == trn && aw == 1 {
			count = count + 1
		}
	}
	if count == 0 {
		fmt.Fprintf(w, "You can't recieve money because you haven't participated in any event")
		return
	}
	_, err = database.Exec("UPDATE user SET coins = coins - ?  WHERE rollno = ? AND coins - ? >= 0", c, rnf, c)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error")
		return
	}
	_, err = database.Exec("UPDATE user SET coins = coins + ? - ? WHERE rollno = ?", c, tax, rnt)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error")
		return
	}
	tx.Commit()
	dbt(0, rnf, rnt, c, tax, time.Now().String())
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
