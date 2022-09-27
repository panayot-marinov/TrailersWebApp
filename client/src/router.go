package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func SetupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", Get).Methods(http.MethodGet)
	r.HandleFunc("/login", Login).Methods(http.MethodGet)
	r.HandleFunc("/makeLoginRequest", MakeLoginRequest).Methods(http.MethodPost)
	r.HandleFunc("/register", Register).Methods(http.MethodGet)
	r.HandleFunc("/makeRegisterRequest", MakeRegisterRequest).Methods(http.MethodPost)

	//api := r.PathPrefix("/api/v1").Subrouter()

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func Get(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-type", "application/json")
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(`{"message": "get called"}`))

	tpl.ExecuteTemplate(w, "index.html", nil) //Read about nginx
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func MakeLoginRequest(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)
	resp, _ := http.PostForm("http://localhost:8081/api/v1/login",
		params)

	fmt.Println("respBody")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))
	fmt.Println("code = " + strconv.Itoa(resp.StatusCode))

	if resp.StatusCode == http.StatusUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusUnauthorized)
	} else if resp.StatusCode != http.StatusFound {
		w.WriteHeader(http.StatusInternalServerError)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusInternalServerError)
	}

	//find cookie
	var cookie *http.Cookie = nil
	for _, currentCookie := range resp.Cookies() {
		if currentCookie.Name == "session_token" {
			cookie = currentCookie
		}
	}
	print("cookie.Value = " + cookie.Value)
	if cookie == nil {
		fmt.Println("there is no cookie in client")
	} else {
		fmt.Println("there is cookie in client")
		http.SetCookie(w, cookie)
	}

	//-----
	destUrl := "http://localhost:8080/"
	http.Redirect(w, r, destUrl, http.StatusSeeOther)

	//tpl.ExecuteTemplate(w, "index.html", cookie)

	// w.Header().Set("Content-type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// tpl.ExecuteTemplate(w, "index.html", nil)
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func MakeRegisterRequest(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	params := url.Values{}
	params.Add("username", username)
	params.Add("email", email)
	params.Add("password", password)
	resp, _ := http.PostForm("http://localhost:8081/api/v1/register",
		params)

	fmt.Println("respBody")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))

	destUrl := "http://localhost:8080/"
	http.Redirect(w, r, destUrl, http.StatusFound)
}
