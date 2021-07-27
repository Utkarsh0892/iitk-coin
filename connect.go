package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
)

func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	rollno := r.FormValue("rollno")
	password := r.FormValue("password")
	rn, _ := strconv.Atoi(rollno)
	if !userExists(rn) {
		fmt.Fprintf(w, "You haven't signed up yet")
		return
	}
	rows, _ :=
		DBU.Query("SELECT rollno, password FROM user")
	var rolln int
	var pwd string
	for rows.Next() {
		rows.Scan(&rolln, &pwd)
		if rn == rolln {
			if comparePasswords(pwd, []byte(password)) {
				generateJWT(w, r)
			} else {
				fmt.Println("Incorrect Password")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollno := r.FormValue("rollno")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")
	fmt.Println("ok1")
	password = encrypt([]byte(password))
	fmt.Println("ok2")
	rn, _ := strconv.Atoi(rollno)
	fmt.Println("ok3")
	dbu(rn, name, password, email, 0)
	fmt.Println("ok4")
	dbc(rn, 0)
	fmt.Println("ok5")
	fmt.Fprintf(w, "SignUp Succesful")
}

func Home(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
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
		rows, _ :=
			DBU.Query("SELECT rollno, name FROM user")
		var rolln int
		var nm string
		var name string
		rn, _ := strconv.Atoi(claims.Rollno)
		for rows.Next() {
			rows.Scan(&rolln, &nm)
			if rn == rolln {
				name = nm
			}
		}
		w.Write([]byte(fmt.Sprintf("Hello, %s \n\nUse /logout endpoint to logout", name)))
	}
}

func logout(response http.ResponseWriter, request *http.Request) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}

func updateInfo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
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
		password := r.FormValue("password")
		email := r.FormValue("email")
		name := r.FormValue("name")
		cr, _ := strconv.Atoi(claims.Rollno)
		rows, err :=
			DBU.Query("SELECT rollno, name, password, email FROM user")
		if err != nil {
			fmt.Println("Error")
			return
		}
		DBU.Exec("PRAGMA journal_mode=WAL;")
		var rolln int
		var nam string
		var pwd string
		var em string
		for rows.Next() {
			rows.Scan(&rolln, &nam, &pwd, &em)
			if cr == rolln {
				if name == "" {
					name = nam
				}
				if comparePasswords(encrypt([]byte("")), []byte(password)) {
					password = pwd
				} else {
					password = encrypt([]byte(password))
				}
				if email == "" {
					email = em
				}
			}
		}
		tx, err := DBU.Begin()
		if err != nil {
			fmt.Println(err)
		}
		_, err = DBU.Exec("UPDATE user SET name = ? WHERE rollno = ?", name, cr)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = DBU.Exec("UPDATE user SET password = ? WHERE rollno = ?", password, cr)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = DBU.Exec("UPDATE user SET email = ? WHERE rollno = ?", email, cr)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = tx.Commit()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Fprintf(w, "Info Updated Succesfully")
	} else {
		fmt.Fprintf(w, "Login to update your info")
	}
}
