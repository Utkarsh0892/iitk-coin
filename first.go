package main

import (
	"database/sql"
	"fmt"
	"strconv"
  _	"github.com/mattn/go-sqlite3"
)

func main() {
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, name TEXT)")
statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO user (rollno, name) VALUES (?, ?)")
statement.Exec(1, "yash")
	rows, _ :=
		database.Query("SELECT rollno, name FROM user")
	var rollno int
	var name string
	for rows.Next() {
		rows.Scan(&rollno, &name)
		fmt.Println(strconv.Itoa(rollno) + ": " + name)
	}
}