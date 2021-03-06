package main

import (
	"fmt"
	"log"
	// "math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	//gomail "gopkg.in/mail.v2"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

var jwtKey = []byte("utkarshs")

type Claims struct {
	Name   string `json:"name"`
	Rollno string `json:"rollno"`
	jwt.StandardClaims
}

func encrypt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {

	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
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
			return false
		}
		w.WriteHeader(http.StatusBadRequest)
		return false
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	return true
}

func generateJWT(w http.ResponseWriter, r *http.Request) {
	rollno := r.FormValue("rollno")
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Rollno: rollno,
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
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func userExists(rn int) bool {
	rows, _ :=
		DBU.Query("SELECT rollno from USER")
	var rollno int
	f := 0
	for rows.Next() {
		rows.Scan(&rollno)
		{
			if rollno == rn {
				f = 1
			}
		}
	}
	if f == 1 {
		return true
	} else {
		return false
	}
}

func isAdmin(rn int) bool {
	rows, _ :=
		DBU.Query("SELECT rollno, isAdmin from USER")
	var rollno int
	var adm int
	f := 0
	for rows.Next() {
		rows.Scan(&rollno, &adm)
		{
			if rollno == rn && adm == 1 {
				f = 1
			}
		}
	}
	if f == 1 {
		return true
	} else {
		return false
	}
}

func makeAdmin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Println("POST request successful")
	rollno := r.FormValue("rollno")
	rn, _ := strconv.Atoi(rollno)
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
		rolln, _ := strconv.Atoi(claims.Rollno)
		if !isAdmin(rolln) {
			fmt.Fprintf(w, "Only ADMIN can access this page")
			return
		}

		tx, err := DBU.Begin()
		if err != nil {
			fmt.Println(err)
		}
		_, err = DBU.Exec("UPDATE user SET isAdmin = ? WHERE rollno = ?", 1, rn)
		if err != nil {
			fmt.Println(err)
			return
		}
		err = tx.Commit()
		if err != nil {
			fmt.Println("Error")
			return
		}
		fmt.Fprintf(w, "%d is now an admin", rn)
	} else {
		fmt.Fprintf(w, "Login as ADMIN to make new admins")
	}
}

func isAdminLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return false
	}
	if isLoggedIn(w, r) {
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return false
			}
			w.WriteHeader(http.StatusBadRequest)
			return false
		}

		tokenStr := cookie.Value

		claims := &Claims{}

		_, _ = jwt.ParseWithClaims(tokenStr, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
		rolln, _ := strconv.Atoi(claims.Rollno)
		return isAdmin(rolln)
	} else {
		return false
	}
}

// func generateOTP(rn int) {
// 	otp := 0
// 	for i := 0; i < 4; i++ {
// 		d := rand.Intn(10)
// 		otp = otp*10 + d
// 	}

// 	dbo(rn,otp)
// 	new := gomail.NewMessage()

// 	new.SetHeader("From", "utkarshs20@iitk.ac.in")
// 	new.SetHeader("To", strconv.Itoa(rn)+"@iitk.ac.in")
// 	new.SetHeader("Subject", "OTP")
// 	new.SetBody("text/plain", "Your OTP is "+strconv.Itoa(otp))

// 	a := gomail.NewDialer("mmtp.iitk.ac.in", 25, "utkarshs20@gmail.com", "<password>")

// 	if err := a.DialAndSend(new); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

// func validateOTP(w http.ResponseWriter, r *http.Request){
// 	if err := r.ParseForm(); err != nil {
// 		fmt.Fprintf(w, "ParseForm() err: %v", err)
// 		return
// 	}
// 	rollno := r.FormValue("rollno")
// 	rn, _ := strconv.Atoi(rollno)
// 	ot := r.FormValue("otp")
// 	otp, _ := strconv.Atoi(ot)
//     rows,err:= 
// 	    DBO.Query("SELECT rollno, otp FROM otp")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	var rolln int
// 	var o int
// 	for rows.Next(){
//         rows.Scan(&rolln,&o)
// 		{
// 			if rolln==rn {
// 				if (o==otp){
// 					http.Redirect(w, r, "/redeem", http.StatusSeeOther)
// 				}else{
// 					fmt.Fprintf(w,"Incorrect OTP")
// 					return 
// 				}
// 			}
// 		}
// 	}
// }
