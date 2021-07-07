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

func OpenDB() {
	dbu, err := sql.Open("sqlite3", "./user.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, _ :=
		dbu.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, name TEXT, password TEXT, email TEXT, isAdmin INTEGER)")
	statement.Exec()
	statement, _ =
		dbu.Prepare("INSERT INTO user (rollno, name, password, email, isAdmin) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(1, "SuperAdmin", "password", "email", 1)
	DBU = dbu
	dbc, err := sql.Open("sqlite3", "./balance.db")
	if err != nil {
		fmt.Println(err)
	}
	statement, _ =
		dbc.Prepare("CREATE TABLE IF NOT EXISTS balance ( rollno INTEGER, coins INTEGER)")
	statement.Exec()
	DBC = dbc
	dbt, err := sql.Open("sqlite3", "./transactions.db")
	statement, _ =
		dbt.Prepare("CREATE TABLE IF NOT EXISTS transactions (award INTEGER, from INTEGER, to INTEGER, coins INTEGER, tax INTEGER, timestamp TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	DBT = dbt
	dbr, err := sql.Open("sqlite3", "./redeem_requests.db")
	statement, _ =
		dbr.Prepare("CREATE TABLE IF NOT EXISTS redeem_requests (rollno INTEGER, item TEXT, coins  INTEGER, status TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	DBR = dbr
}

func dbu(rn int, name string, password string, email string, adm int) {
	statement, _ :=
		DBU.Prepare("INSERT INTO user (rollno, name, password, email, isAdmin) VALUES (?, ?, ?, ?, ?)")
	statement.Exec(rn, name, password, email, adm)
}

func dbc(rn int, coin int) {

	statement, _ :=
		DBC.Prepare("INSERT INTO balance (rollno, coins) VALUES (?, ?)")
	statement.Exec(rn, coin)
}

func dbt(aw int, frn int, trn int, coin int, tax int, ts string) {
	statement, _ :=
		DBT.Prepare("INSERT INTO transactions (award, from ,to , coins, tax, timestamp) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(aw, frn, trn, coin, tax, ts)
}

func dbr(rn int, i string, coins int, stat string) {
	statement, _ :=
		DBR.Prepare("INSERT INTO redeem_requests (rollno, item, coins, status) VALUES (?, ?, ?, ?)")
	statement.Exec(rn, i, coins, stat)
}
