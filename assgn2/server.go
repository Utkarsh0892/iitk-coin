package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

var jwtKey = []byte("secret_key")

type Claims struct {
	rollno string `json:"rollno"`
	jwt.StandardClaims
}

func db(rn int, name string, password string, email string) {
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno INTEGER, name TEXT,password TEXT,email TEXT)")
	statement.Exec()
	statement, _ =
		database.Prepare("INSERT INTO user (rollno, name, password, email) VALUES (?, ?, ?, ?)")
	statement.Exec(rn, name, password, email)
}
func encrypt(string) (password string) {
	salt := make([]byte, PW_SALT_BYTES)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
	rollno := r.FormValue("rollno")
	password := r.FormValue("password")
	email := r.FormValue("email")
	name := r.FormValue("name")
	password = encrypt(password)
	rn, _ := strconv.Atoi(rollno)
	db(rn, name, password, email)
}

func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	rollno := r.FormValue("rollno")
	password := r.FormValue("password")
	password = encrypt(password)
	rn, _ := strconv.Atoi(rollno)
	database, _ :=
		sql.Open("sqlite3", "./user.db")
	rows, _ :=
		database.Query("SELECT rollno, password FROM user")
	var rolln int
	var pwd string
	for rows.Next() {
		rows.Scan(&rolln, &pwd)
		if rn == rolln && password == pwd {

			expirationTime := time.Now().Add(time.Minute * 5)

			claims := &Claims{
				rollno: rollno,
				StandardClaims: jwt.StandardClaims{
					ExpiresAt: expirationTime.Unix(),
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(jwtKey)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			http.SetCookie(w,
				&http.Cookie{
					Name:    "token",
					Value:   tokenString,
					Expires: expirationTime,
				})
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
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

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Hello, %s", claims.rollno)))

}

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
