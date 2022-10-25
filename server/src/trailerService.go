package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (authHandler *AuthHandler) GetTrailerData(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")
	print("a0")
	hasValidData, _, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	print("a1\n")

	db := ConnectToDb()
	defer db.Close()
	//TODO: remove that
	from := time.Date(2022, 10, 16, 10, 10, 10, 10, time.UTC)
	to := time.Now()
	trailerData, err := GetTrailerDataFromDb(db, from, to)
	print("a2\n")
	if err != nil {
		print("a3\n")
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}
	print("a4\n")

	jData, err := json.Marshal(trailerData)
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

func (authHandler *AuthHandler) GetTrailersList(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")
	print("a0")
	hasValidData, _, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	print("a1\n")

	db := ConnectToDb()
	defer db.Close()
	var trailers []Trailer
	trailers, err := GetTrailersListFromDb(db)
	print("a2\n")
	if err != nil {
		print("a3\n")
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}
	print("a4\n")

	jData, err := json.Marshal(trailers)
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

func (authHandler *AuthHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}
	prevUrl := r.Header.Get("Referer")

	trailerNumber := r.FormValue("trailerNumber")
	trailerName := r.FormValue("trailerName")

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb()
	defer db.Close()

	user, err := GetUserInfoFromDbWithUsername(db, session.username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	print("userid=")
	print(user.Id)
	trailer := Trailer{
		Number: trailerNumber,
		Name:   trailerName,
		UserId: user.Id,
	}

	err = RegisterNewTrailerToDb(db, trailer)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
		print("\n")
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}
}
