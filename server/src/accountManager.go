package src

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-hclog"
)

var sessions = map[string]Session{}

type AuthHandler struct {
	mailService   *SGMailService
	logger        hclog.Logger
	configuration Configuration
}

// NewAuthHandler returns a new instance of AuthHandler
func NewAuthHandler(mailService *SGMailService, logger hclog.Logger, configuration Configuration) *AuthHandler {
	return &AuthHandler{mailService, logger, configuration}
}

func (authHandler *AuthHandler) Welcome(w http.ResponseWriter, r *http.Request) (bool, string, Session, *http.Cookie) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		print("Cookie error\n")
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return false, "", Session{}, &http.Cookie{}
		}
		print("cookie got but bad request")
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return false, "", Session{}, c
	}
	sessionToken := c.Value
	print("session token got")

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	print("USER SESSION GOT")
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return false, "", Session{
				"invalid",
				time.Unix(0, 0),
			},
			c
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		print("user session expired\n")
		sessionToken, userSession, _ = authHandler.Refresh(w, r)
		//delete(sessions, sessionToken)
		//w.WriteHeader(http.StatusUnauthorized)
		//return false, "", userSession
	}

	return true, sessionToken, userSession, c
	// If the session is valid, return the welcome message to the user
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", userSession.username)))
}

func (authHandler *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) (string, Session, *http.Cookie) {
	// (BEGIN) The code from this point is the same as the first part of the `Welcome` route
	sessionGotSuccessfully, sessionToken, userSession, _ := authHandler.Welcome(w, r)
	if !sessionGotSuccessfully {
		print("session not got successfully\n")
		return "", userSession, &http.Cookie{
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
	newSession := Session{
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

	return newSessionToken, newSession, newCookie
}

func (authHandler *AuthHandler) CheckCookie(w http.ResponseWriter, r *http.Request) (bool, Session, *http.Cookie) {
	sessionGotSuccessfully, sessionToken, userSession, cookie := authHandler.Welcome(w, r)
	if !sessionGotSuccessfully && userSession.expiry.Equal(time.Unix(0, 0)) {
		print("cookie invalid\n")
		return false, Session{
				"invalid",
				time.Unix(0, 0),
			}, &http.Cookie{
				Name:    "session_token",
				Value:   sessionToken,
				Expires: time.Unix(0, 0),
			}
	}
	return true, userSession, cookie
}

func (authHandler *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
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
	fmt.Println("user got from db successfully")

	passwordsEqual := bytes.Compare(passwordHash, passwordHashDb)

	if passwordsEqual != 0 {
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusUnauthorized)
		return
	}
	fmt.Println("passwords are equal")

	//Valid password
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(3 * 60 * (60 * time.Second))

	sessions[sessionToken] = Session{
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
	fmt.Println("Logged in")
}

func (authHandler *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	company := r.FormValue("company")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	hashedPassword := hashSha256(password)

	now := time.Now()

	user := User{
		Username:   username,
		Email:      email,
		Password:   hashedPassword,
		Company:    company,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := RegisterNewAccountToDb(db, user)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
		print("\n")
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}

	err = authHandler.SendVerificationMail(db, user)
	if err != nil {
		fmt.Println("Error sending verification email")
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func (authHandler *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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

func (authHandler *AuthHandler) GetUserProfileInfo(w http.ResponseWriter, r *http.Request) {
	print("a0\n")
	prevUrl := r.Header.Get("Referer")

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	print("a1\n")

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	user, err := GetUserInfoFromDbWithUsername(db, session.username)
	user.Username = session.username
	print("a2\n")
	if err != nil {
		print("a3\n")
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}
	print("a4\n")
	user.Password = nil

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(account)

	jData, err := json.Marshal(user)
	print(user.Username)
	print(user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	print(string(jData[:]))
	print("\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
	print("a5\n")
}

func (authHandler *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")

	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("password " + password + " passwordRepeated " + passwordRepeated + "\n")

	if password != passwordRepeated {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	hashedPassword := hashSha256(password)

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	err := UpdateAccountPasswordToDb(db, session.username, string(hashedPassword))
	if err != nil {
		fmt.Println("Error executing insert statement")
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func (authHandler *AuthHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")

	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("password " + password + " passwordRepeated " + passwordRepeated + "\n")

	if password != passwordRepeated {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	account, err := GetAccountInfoFromDb(db, session.username)
	account.Username = session.username
	if err != nil {
		print("a3\n")
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	err = DeleteVerificationDataFromDb(db, account.Email, MailConfirmation)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	err = DeleteVerificationDataFromDb(db, account.Email, PassReset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	err = DeleteAccountFromDb(db, session.username)
	if err != nil {
		fmt.Println("Error executing delete statement")
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}

func (authHandler *AuthHandler) PasswordReset(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")

	username := r.FormValue("username")
	code := r.FormValue("code")
	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("username " + username + "code" + code + " password " + password + " passwordRepeated " + passwordRepeated + "\n")

	if password != passwordRepeated {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}

	hashedPassword := hashSha256(password)

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	account, err := GetAccountInfoFromDb(db, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	var verificationData VerificationData
	verificationData.Email = account.Email
	verificationData.Code = code

	actualVerificationData, err := GetVerificationDataFromDb(db, account.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	valid, err := authHandler.Verify(db, &actualVerificationData, &verificationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}

	err = DeleteVerificationDataFromDb(db, account.Email, PassReset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	err = UpdateAccountPasswordToDb(db, username, string(hashedPassword))
	if err != nil {
		fmt.Println("Error executing insert statement")
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
}
