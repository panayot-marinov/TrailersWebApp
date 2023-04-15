package src

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	//TODO: remove that
	from := time.Date(2022, 10, 16, 10, 10, 10, 10, time.UTC)
	to := time.Now()
	//trailerData, err := GetTrailersDataFromDb(db, from, to)
	//print("a2\n")
	//if err != nil {
	//	print("a3\n")
	//	w.WriteHeader(http.StatusBadRequest)
	//	http.Redirect(w, r, prevUrl, http.StatusBadRequest)
	//	return
	//}
	//print("a4\n")

	trailers, err := GetTrailersListFromDb(db, session.username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Redirect(w, r, prevUrl, http.StatusBadRequest)
		return
	}

	var trailerDataMap map[string][]TrailerData
	trailerDataMap = make(map[string][]TrailerData)
	for _, trailer := range trailers {
		trailerData, err := GetTrailerDataFromDb(db, from, to, trailer.RegistrationPlate, session.username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Redirect(w, r, prevUrl, http.StatusBadRequest)
			return
		}
		print("a4\n")
		trailerDataMap[trailer.RegistrationPlate] = trailerData
	}

	jTrailerData, err := json.Marshal(trailerDataMap)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	print(string(jTrailerData[:]))
	print("\n")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jTrailerData)
	print("a5\n")
}

func (authHandler *AuthHandler) GetTrailersList(w http.ResponseWriter, r *http.Request) {
	prevUrl := r.Header.Get("Referer")
	print("a0")
	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	print("a1\n")

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	var trailers []Trailer
	trailers, err := GetTrailersListFromDb(db, session.username)
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

	//trailerRegPlate := r.FormValue("trailerNumber")
	//trailerName := r.FormValue("trailerName")

	var trailer Trailer
	trailer.Brand = r.FormValue("brand")
	trailer.Model = r.FormValue("model")
	trailer.Name = r.FormValue("nameField")
	trailer.RegistrationPlate = r.FormValue("registrationPlate")
	trailer.SerialNumber = r.FormValue("serialNumber")
	trailer.City = r.FormValue("city")
	trailer.Area = r.FormValue("area")
	trailer.AddressLine = r.FormValue("addressLine")
	zipCode, err := strconv.Atoi(r.FormValue("zipCode"))
	if err != nil {
		print("Cannot parse zipCode to int")
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	trailer.ZipCode = zipCode

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	user, err := GetUserInfoFromDbWithUsername(db, session.username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	print("userid=")
	print(user.Id)
	trailer.UserId = user.Id
	//trailer := Trailer{
	//	RegistrationPlate: trailerRegPlate,
	//	Name:              trailerName,
	//	UserId:            user.Id,
	//}

	err = RegisterNewTrailerToDb(db, trailer)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
		print("\n")
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}
}

func (authHandler *AuthHandler) Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		//destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}
	prevUrl := r.Header.Get("Referer")

	var trailer Trailer
	trailer.Brand = r.FormValue("brand")
	trailer.Model = r.FormValue("model")
	trailer.Name = r.FormValue("nameField")
	trailer.RegistrationPlate = r.FormValue("registrationPlate")
	trailer.SerialNumber = r.FormValue("serialNumber")
	isActiveStr := r.FormValue("isActive")
	fmt.Println("isActiveStr = " + isActiveStr)
	isActive, err := strconv.ParseBool(isActiveStr)
	if err != nil {
		print("Cannot parse isActive to bool")
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	trailer.IsActive = isActive
	trailer.City = r.FormValue("city")
	trailer.Area = r.FormValue("area")
	trailer.AddressLine = r.FormValue("addressLine")
	zipCode, err := strconv.Atoi(r.FormValue("zipCode"))
	if err != nil {
		print("Cannot parse zipCode to int")
		print("zipCode = ")
		print(zipCode)
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	trailer.ZipCode = zipCode

	hasValidData, session, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	user, err := GetUserInfoFromDbWithUsername(db, session.username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}
	print("userid=")
	print(user.Id)
	trailer.UserId = user.Id

	print("updating trailer")
	err = UpdateTrailerIntoDb(db, trailer)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
		print("\n")
		w.WriteHeader(http.StatusUnauthorized)
		panic(err)
	}
}

func (authHandler *AuthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleting trailers...")
	if r.Method != "POST" {
		//destUrl := "http://localhost:8080/"
		w.WriteHeader(http.StatusUnauthorized)
		//http.Redirect(w, r, destUrl, http.StatusAccepted)
		return
	}
	prevUrl := r.Header.Get("Referer")

	hasValidData, _, _ := authHandler.CheckCookie(w, r)
	if !hasValidData {
		print("no valid data changing password")
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		fmt.Println("cannot parse form")
		fmt.Println(err)
	}

	trailersCheckboxes := r.Form["trailerCheckbox"]
	for _, registrationPlate := range trailersCheckboxes {
		print("deleting trailer")
		err := DeleteTrailerFromDb(db, registrationPlate)
		if err != nil {
			fmt.Println("Error executing insert statement")
			print(err)
			print("\n")
			w.WriteHeader(http.StatusUnauthorized)
			panic(err)
		}
	}
}
