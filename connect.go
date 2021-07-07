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
	password = encrypt([]byte(password))
	rn, _ := strconv.Atoi(rollno)
	dbu(rn, name, password, email, 0)
	dbc(rn, 0)
	fmt.Fprintf(w, "SignUp Succesful")
}

func Home(w http.ResponseWriter, r *http.Request) {
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
