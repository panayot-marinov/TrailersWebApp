package src

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var sessions = map[string]session{}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	db := ConnectToDb()
	defer db.Close()
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordHash := hashSha256(password)

	fmt.Println("username =" + username)

	var passwordHashDb []byte
	row := db.QueryRow("SELECT \"password\" FROM \"Users\" WHERE \"username\"=$1", username)
	if err := row.Scan(&passwordHashDb); err != nil {
		fmt.Println("ERROR! Cannot execute select query!")
		w.WriteHeader(http.StatusSeeOther)
		//http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	passwordsEqual := bytes.Compare(passwordHash, passwordHashDb)

	if passwordsEqual != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnauthorized)
		return
		//w.Write([]byte("Incorrect username or password!"))
	}

	//Valid password
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	sessions[sessionToken] = session{
		username: username,
		expiry:   expiresAt,
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
		Path:    "/",
	})

	//url := r.Header.Get("Referer")
	//url := "https://www.google.bg/"
	w.WriteHeader(http.StatusFound)
	//http.Redirect(w, r, url, http.StatusFound)
	print("Logged in")
	//w.Write([]byte("Successfully logged in!"))
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusAccepted)
		http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}

	db := ConnectToDb()
	defer db.Close()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword := hashSha256(password)

	fmt.Println("username =" + username)
	fmt.Println("email =" + email)
	//fmt.Println("password = " + hashedPassword)

	query := "INSERT INTO \"Users\" (username, email, password) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, username, email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing insert statement")
		panic(err)
	}
}
