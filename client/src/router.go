package src

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func SetupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/", Get).Methods(http.MethodGet)
	r.HandleFunc("/login", Login).Methods(http.MethodGet)
	r.HandleFunc("/makeLoginRequest", MakeLoginRequest).Methods(http.MethodPost)
	r.HandleFunc("/register", Register).Methods(http.MethodGet)
	r.HandleFunc("/makeRegisterRequest", MakeRegisterRequest).Methods(http.MethodPost)
	r.HandleFunc("/logout", MakeLogoutRequest).Methods(http.MethodPost)
	r.HandleFunc("/account", GetAccountInfo).Methods(http.MethodGet)
	r.HandleFunc("/makeChangePasswordRequest", MakeChangePasswordRequest).Methods(http.MethodPost)

	//api := r.PathPrefix("/api/v1").Subrouter()

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func Get(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-type", "application/json")
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(`{"message": "get called"}`))
	print("In index\n")

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
		//w.WriteHeader(http.StatusFound)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	} else if resp.StatusCode != http.StatusOK {
		//w.WriteHeader(http.StatusFound)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
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

func MakeLogoutRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("before making request")

	prevUrl := r.Header.Get("Referer")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, prevUrl, http.StatusUnauthorized)
		return
	}
	fmt.Println("cookie_value = " + cookie.Value)

	jar, err := cookiejar.New(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Redirect(w, r, prevUrl, http.StatusInternalServerError)
		return
	}

	client := &http.Client{
		Jar: jar,
	}

	urlObj, _ := url.Parse("http://localhost:8081/api/v1/logout")
	client.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	resp, _ := client.PostForm("http://localhost:8081/api/v1/logout", url.Values{})

	fmt.Println("request made successfully")
	if resp.StatusCode == http.StatusUnauthorized {
		//w.WriteHeader(http.StatusOK)
		fmt.Println("redirect 0")
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusSeeOther)
		return
	}
	// } else if resp.StatusCode != http.StatusAccepted {
	// 	//w.WriteHeader(http.Sta)
	// 	fmt.Println("redirect 1")
	// 	destUrl := "http://localhost:8080/login"
	// 	http.Redirect(w, r, destUrl, http.StatusSeeOther)
	// 	return
	// }

	cookie = &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),

		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	fmt.Println("respBody")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))

	fmt.Println("redirect 2")
	destUrl := r.Header.Get("Referer")
	http.Redirect(w, r, destUrl, http.StatusSeeOther)
}

func GetAccountInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/api/v1/account", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}
	decoder := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	var account Account
	err = decoder.Decode(&account)
	print(account.Username)
	print(account.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	print("aaa")
	// 	tpl.ExecuteTemplate(w, "profile.html", nil)
	// 	return
	// }
	// print("body\n")
	// print(string(body[:]))
	// print("\n")
	// defer resp.Body.Close()

	// err = json.Unmarshal(body, &account)
	// print(account.Username)
	// print(account.Email)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	print("bbb")
	// 	tpl.ExecuteTemplate(w, "profile.html", nil)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)
	params := url.Values{}
	params.Add("username", account.Username)
	params.Add("email", account.Email)

	tpl.ExecuteTemplate(w, "account.html", account)
}

func MakeChangePasswordRequest(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("password " + password + " passwordRepeated " + passwordRepeated + "\n")

	client := &http.Client{}
	params := url.Values{}
	params.Add("password", password)
	params.Add("passwordRepeated", passwordRepeated)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/changePassword", strings.NewReader(params.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	print("request made \n")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}

	print("cookie got successfully\n")
	print(cookie)
	req.AddCookie(cookie)
	print("1")
	print("2")
	print("before making request")
	resp, err := client.Do(req)
	print("request made")
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "account.html", nil)
		return
	}
	print("request made successfully\n")
	print("resp body")
	print(resp.Body)

	//decoder := json.NewDecoder(resp.Body)

	//defer resp.Body.Close()

	fmt.Println("respBody")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))
	fmt.Println("code = " + strconv.Itoa(resp.StatusCode))

	if resp.StatusCode == http.StatusUnauthorized {
		print("1\n")
		w.WriteHeader(http.StatusUnauthorized)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	} else if resp.StatusCode != http.StatusOK {
		print("2\n")
		w.WriteHeader(http.StatusBadRequest)
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	}

	print("3\n")
	//-----
	destUrl := "http://localhost:8080/account"
	http.Redirect(w, r, destUrl, http.StatusSeeOther)

	//tpl.ExecuteTemplate(w, "index.html", cookie)

	// w.Header().Set("Content-type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// tpl.ExecuteTemplate(w, "index.html", nil)
}
