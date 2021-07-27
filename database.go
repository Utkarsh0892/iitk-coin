package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DBU *sql.DB
var DBC *sql.DB
var DBT *sql.DB
var DBR *sql.DB
//var DBO *sql.DB

func OpenDB() {
	dbu, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, _ :=
		dbu.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, name TEXT, password TEXT, email TEXT, isAdmin INTEGER, unique(rollno))")
	statement.Exec()
	statement, _ =
		dbu.Prepare("INSERT INTO user (rollno, name, password, email, isAdmin) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(1, "SuperAdmin", encrypt([]byte("password")), "email", 1)
	DBU = dbu
	dbc, err := sql.Open("sqlite3", "./balance.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, _ =
		dbc.Prepare("CREATE TABLE IF NOT EXISTS balance ( rollno INTEGER, coins INTEGER, unique(rollno))")
	statement.Exec()
	DBC = dbc
	dbt, err := sql.Open("sqlite3", "./transactions.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, err =
		dbt.Prepare("CREATE TABLE IF NOT EXISTS transactions (type TEXT, fromrn INTEGER, torn INTEGER, coins INTEGER, tax INTEGER, timestamp TEXT)")
	if err != nil {
		fmt.Println(err)
	}
	statement.Exec()
	DBT = dbt
	dbr, err := sql.Open("sqlite3", "./redeem_requests.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, _ =
		dbr.Prepare("CREATE TABLE IF NOT EXISTS redeem_requests (rollno INTEGER, item TEXT, coins  INTEGER, status TEXT)")
	statement.Exec()
	DBR = dbr
	// dbo, err := sql.Open("sqlite3", "./otp.db")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// statement, _ =
	// 	dbo.Prepare("CREATE TABLE IF NOT EXISTS otp (rollno INTEGER, otp INTEGER,unique(rollno))")
	// statement.Exec()
	// DBO = dbo
}

func dbu(rn int, name string, password string, email string, adm int) {
	statement, _ :=
		DBU.Prepare("INSERT INTO user (rollno, name, password, email, isAdmin) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(rn, name, password, email, adm)
}

func dbc(rn int, coin int) {
	fmt.Println("ok6")

	statement, _ :=
		DBC.Prepare("INSERT INTO balance (rollno, coins) VALUES (?, ?)")
	fmt.Println("ok7")
	statement.Exec(rn, coin)
	fmt.Println("ok8")
}

func dbt(ty string, frn int, trn int, coin int, tax int, ts string) {
	statement, _ :=
		DBT.Prepare("INSERT INTO transactions (type, fromrn , torn, coins, tax, timestamp) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(ty, frn, trn, coin, tax, ts)
}

func dbr(rn int, i string, coins int, stat string) {
	statement, _ :=
		DBR.Prepare("INSERT INTO redeem_requests (rollno, item, coins, status) VALUES (?, ?, ?, ?)")
	statement.Exec(rn, i, coins, stat)
}

// func dbo(rn int, otp int) {
// 	statement, _ :=
// 		DBO.Prepare("INSERT INTO otp (rollno, otp) VALUES (?, ?)")
// 	statement.Exec(rn, otp)
// }