package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	_ "github.com/mattn/go-sqlite3"
)

func redeem(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	i := r.FormValue("item")
	c, _ := strconv.Atoi(r.FormValue("coins"))
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
		rn, _ := strconv.Atoi(claims.Rollno)
		stat := "pending"
		rows, _ :=
			DBR.Query("SELECT rollno,status from redeem_requests")
		var rollno int
		var status string
		f := 0
		for rows.Next() {
			rows.Scan(&rollno, &status)
			{
				if rollno == rn {
					stat = status
					f = 1
				}
			}
		}
		if f == 1 {
			fmt.Fprintf(w, "Your request is %s", stat)
		} else {
			dbr(rn, i, c, "pending")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(w, "Your request has been sent to admin")
		}
	} else {
		fmt.Fprintf(w, "Login to redeem")
		return
	}
}

func manageRedeemRequests(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollno := r.FormValue("rollno")
	i := r.FormValue("item")
	stat := r.FormValue("action")
	rn, _ := strconv.Atoi(rollno)
	if isAdminLoggedIn(w, r) {
		txr, err := DBR.Begin()
		if err != nil {
			fmt.Println(err)
		}
		txc, err := DBC.Begin()
		if err != nil {
			fmt.Println(err)
		}
		DBR.Exec("PRAGMA journal_mode=WAL;")
		rows, err :=
			DBR.Query("SELECT rollno, item, coins FROM redeem_requests")
		if err != nil {
			fmt.Println(err)
			return
		}
		var rolln int
		var coin int
		var item string
		f := 0
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&rolln, &item, &coin)
			if rn == rolln && item == i {
				_, err := DBR.Exec("UPDATE redeem_requests SET status = ? WHERE rollno = ? AND item = ?", stat, rn, i)
				if err != nil {
					fmt.Println(err)
					return
				}
				if stat == "approved" {
					res, err := DBC.Exec("UPDATE balance SET coins = coins - ?  WHERE rollno = ? AND coins - ? >= 0", coin, rn, coin)
					if err != nil {
						fmt.Println("Error")
						txc.Rollback()
						return
					}
					n, err := res.RowsAffected()
					if err != nil {
						fmt.Println("Error")
						txc.Rollback()
						return
					}
					if n == 0 {
						fmt.Fprintf(w, "Requester doesn't have enough balance")
						txc.Rollback()
						return
					}
					dbt("redeem", rn, 1, coin, 0, time.Now().String())
				}
				f = 1
			}
		}
		err = txc.Commit()
		if err != nil {
			fmt.Println("Error")
			return
		}
		err = txr.Commit()
		if err != nil {
			fmt.Println("Error")
			return
		}
		if f == 0 {
			fmt.Fprintf(w, "Request doesn't exist")
			return
		}
		fmt.Fprintf(w, "The request has been %s", stat)
	} else {
		fmt.Fprintf(w, "Login as ADMIN to manage redeem requests")
	}
}

func viewRedeemRequests(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollno := r.FormValue("rollno")
	rn, _ := strconv.Atoi(rollno)
	rows, err :=
		DBR.Query("SELECT rollno, item, status FROM redeem_requests")
	if err != nil {
		fmt.Println(err)
		return
	}
	if !(userExists(rn) || (rn == 0)) {
		fmt.Fprintf(w, "User doesn't exist")
		return
	}
	var rolln int
	var item string
	var stat string
	f := 0
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&rolln, &item, &stat)
		if rn == rolln || rn == 0 {
			if f == 0 {
				fmt.Fprintf(w, "rollno\titem\tstatus\n")
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", rolln, item, stat)
			f = 1
		}
	}
	if f == 0 {
		fmt.Fprintf(w, "No Requests has been made from this user yet")
	}
}
