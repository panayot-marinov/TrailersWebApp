package src

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var sessions = map[string]session{}

func Welcome(w http.ResponseWriter, r *http.Request) (bool, string, session) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return false, "", session{}
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return false, "", session{}
	}
	sessionToken := c.Value

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return false, "", session{
			"notExists",
			time.Unix(0, 0),
		}
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return false, "", userSession
	}

	return true, sessionToken, userSession
	// If the session is valid, return the welcome message to the user
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) *http.Cookie {
	// (BEGIN) The code from this point is the same as the first part of the `Welcome` route
	sessionGotSuccessfully, sessionToken, userSession := Welcome(w, r)
	if !sessionGotSuccessfully {
		return &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Unix(0, 0),
		}
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// If the previous session is valid, create a new session token for the current user
	newSessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the user whom it represents
	newSession := session{
		username: userSession.username,
		expiry:   expiresAt,
	}
	sessions[newSessionToken] = newSession

	// Delete the older session token
	delete(sessions, sessionToken)

	// Set the new token as the users `session_token` cookie
	newCookie := &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(120 * time.Second),
	}
	http.SetCookie(w, newCookie)

	return newCookie
}

func CheckCookie(w http.ResponseWriter, r *http.Request) (bool, *http.Cookie) {
	sessionGotSuccessfully, sessionToken, userSession := Welcome(w, r)
	if !sessionGotSuccessfully && userSession.expiry.Equal(time.Unix(0, 0)) {
		return false, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Unix(0, 0),
		}
	} else {
		refreshedCookie := Refresh(w, r)
		return true, refreshedCookie
	}
}

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
	w.WriteHeader(http.StatusFound)
	//http.Redirect(w, r, url, http.StatusFound)
	print("Logged in")
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}

	db := ConnectToDb()
	defer db.Close()
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword := hashSha256(password)

	query := "INSERT INTO \"Users\" (username, email, password) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, username, email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing insert statement")
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	print("0\n")
	prevUrl := r.Header.Get("Referer")
	if r.Method != "POST" {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	print("1\n")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("2\n")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	if _, ok := sessions[cookie.Value]; ok {
		print("3\n")
		delete(sessions, cookie.Value)
		cookie = &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),

			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusSeeOther)
		http.Redirect(w, r, prevUrl, http.StatusSeeOther)
		print("4\n")
		return
	} else {
		print("5\n")
		cookie = &http.Cookie{
			Name:    "session_token",
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),

			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
}
