// Auth test task
package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"database/sql"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

const (
	SECRET = "tatl_tech_task"

	DB_NAME     = "tatltest"
	DB_USER     = "root"
	DB_PASSWORD = "123456aA"
)

var db *sql.DB

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	r.ParseForm()

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var user User
	err := db.QueryRow("SELECT id, email FROM user WHERE email = ? AND password = ?", email, password).Scan(&user.Id, &user.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized) // 401
		io.WriteString(w, `{"error":"invalid_credentials"}`)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.Id,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * time.Duration(1)).Unix(), // 1 hour
		"iat":   time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(SECRET))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error":"failed_to_generate_token"}`)
		return
	}
	io.WriteString(w, `{"token":"`+tokenString+`"}`)
	return
}

func Auth(next http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(next)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	io.WriteString(w, `{"status":"success"}`)
}

func main() {
	var err error
	db, err = sql.Open("mysql", DB_USER+":"+DB_PASSWORD+"@tcp(127.0.0.1:3306)/"+DB_NAME)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/auth", AuthHandler)
	http.Handle("/", Auth(http.HandlerFunc(IndexHandler)))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
