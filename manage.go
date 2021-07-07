package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	_ "github.com/mattn/go-sqlite3"
)

func credit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollno := r.FormValue("rollno")
	coins := r.FormValue("coins")
	rn, _ := strconv.Atoi(rollno)
	c, _ := strconv.Atoi(coins)
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	if isAdminLoggedIn(w, r) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value

		claims := &Claims{}

		_, _ = jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
		checkr, _ := strconv.Atoi(claims.Rollno)
		if !isAdmin(checkr) {
			fmt.Fprintf(w, "Only ADMIN can access this page")
			return
		}
		DBC.Exec("PRAGMA journal_mode=WAL;")
		tx, err := DBC.Begin()
		if err != nil {
			fmt.Println(err)
		}
		rows, err :=
			DBC.Query("SELECT rollno, coins FROM balance")
		if err != nil {
			fmt.Println(err)
			return
		}
		var rolln int
		var coin int
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&rolln, &coin)
			if rn == rolln {
				DBC.Exec("UPDATE balance SET coins = coins + ? WHERE rollno = ?", c, rn)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
		err = tx.Commit()
		if err != nil {
			fmt.Println("Error")
			return
		}
		dbt(1, 0, rn, c, 0, time.Now().String())
		fmt.Fprintf(w, "Coin Reward successful\n")
	} else {
		fmt.Fprintf(w, "Login as ADMIN to credit coins")
	}

}

func transfer(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollnof := r.FormValue("rollnofrom")
	rollnot := r.FormValue("rollnoto")
	amt := r.FormValue("coins")
	rnf, _ := strconv.Atoi(rollnof)
	rnt, _ := strconv.Atoi(rollnot)
	c, _ := strconv.Atoi(amt)
	if isLoggedIn(w, r) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value

		claims := &Claims{}

		_, _ = jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
		if claims.Rollno != rollnof {
			w.Write([]byte(fmt.Sprintf("Login as roll number %s to transfer coins from this account", rollnof)))
			return
		}
	} else {
		w.Write([]byte(fmt.Sprintf("Login as roll number %s to transfer coins from this account", rollnof)))
		return
	}
	if !userExists(rnt) {
		fmt.Fprintf(w, "Receiver hasn't signed up yet")
		return
	}
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
	tx, err := DBC.Begin()
	if err != nil {
		fmt.Println(err)
	}
	count := 0
	var aw int
	var trn int
	rows, err :=
		DBT.Query("SELECT award, torn FROM transactions")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
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
	res, err := DBC.Exec("UPDATE balance SET coins = coins - ?  WHERE rollno = ? AND coins - ? >= 0", c, rnf, c)
	if err != nil {
		fmt.Println("Error")
		tx.Rollback()
		return
	}
	n, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error")
		tx.Rollback()
		return
	}
	if n == 0 {
		fmt.Fprintf(w, "Sender doesn't have enough balance")
		tx.Rollback()
		return
	}
	_, err = DBC.Exec("UPDATE balance SET coins = coins + ? - ? WHERE rollno = ?", c, tax, rnt)
	if err != nil {
		fmt.Println("Error")
		tx.Rollback()
		return
	}
	statement, err :=
		DBT.Prepare("INSERT INTO transactions (award, from , to, coins, tax, timestamp) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(0, rnf, rnt, c, tax, time.Now().String())
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return
	}
	fmt.Fprintf(w, "Coin Transfer successful\n")
}

func checkBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("this route can only handle GET request")
		return
	}
	rolln := r.FormValue("rollno")
	rn, _ := strconv.Atoi(rolln)
	fmt.Println("GET request successful")
	if isLoggedIn(w, r) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value

		claims := &Claims{}

		_, _ = jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
		if claims.Rollno != rolln {
			w.Write([]byte(fmt.Sprintf("Login as roll number %s to check your balance", rolln)))
			return
		}
		rows, err := DBC.Query("SELECT rollno ,coins FROM balance")
		if err != nil {
			fmt.Println(err)
		}
		var rollno int
		var coins int
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&rollno, &coins)
			if rollno == rn {
				fmt.Fprintf(w, "you have %d coins", coins)
				return
			}
		}
	} else {
		fmt.Fprintf(w, "Login to check your balance")
		return
	}
}
