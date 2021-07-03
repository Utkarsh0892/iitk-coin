package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

var DBU *sql.DB
var DBC *sql.DB
var DBT *sql.DB

func OpenDB() {
	dbu, err := sql.Open("sqlite3", "./user.db")
	statement, _ :=
		dbu.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, name TEXT,password TEXT,email TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
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
		dbt.Prepare("CREATE TABLE IF NOT EXISTS transactions (award INTEGER, fromrn INTEGER, torn INTEGER, coins INTEGER, tax INTEGER, timestamp TEXT)")
	statement.Exec()
	if err != nil {
		fmt.Println(err)
	}
	DBT = dbt
}

func db(rn int, name string, password string, email string) {	
	statement, _ :=
		DBU.Prepare("INSERT INTO user (rollno, name, password, email) VALUES (?, ?, ?, ?)")
	statement.Exec(rn, name, password, email)
}

func dbc(rn int, coin int) {
	
	statement, _ :=
		DBC.Prepare("INSERT INTO balance (rollno, coins) VALUES (?, ?)")
	statement.Exec(rn, coin)
}

func dbt(aw int, frn int, trn int, coin int, tax int, ts string) {
	statement, _ :=
		DBT.Prepare("INSERT INTO transactions (award, fromrn ,torn , coins, tax, timestamp) VALUES (?, ?, ?, ?, ?, ?)")
	statement.Exec(aw, frn, trn, coin, tax, ts)
}