package src

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func TrailerDataDetails(w http.ResponseWriter, r *http.Request) {
	print("aa0")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/api/v1/trailers/data", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersData.html", nil)
		return
	}

	print("aa1")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersData.html", nil)
		return
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	print("code=")
	print(resp.StatusCode)
	print("aa2")
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersData.html", nil)
		return
	}
	print("aa3")
	// decoder := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	// var trailerData []TrailerData
	// err = decoder.Decode(&trailerData)
	// print("aa4")
	// // print(account.Username)
	// // print(account.Email)
	// if err != nil {
	// 	print("aa5")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	tpl.ExecuteTemplate(w, "trailersData.html", nil)
	// 	return
	// }

	//Array := [5]int{1, 2, 3, 4, 5}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		print("aa6")
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "trailersData.html", nil)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "trailersData.html", template.FuncMap{"jsonData": string(b[:])})
}

func TrailersManager(w http.ResponseWriter, r *http.Request) {
	print("aa0")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/api/v1/trailers/list", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}

	print("aa1")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}
	req.AddCookie(cookie)
	resp, err := client.Do(req)
	print("code=")
	print(resp.StatusCode)
	print("aa2")
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}
	print("aa3")

	//decoder := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	// var trailers []Trailer
	// err = decoder.Decode(&trailers)
	// print("aa4")
	// if err != nil {
	// 	print("aa5")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	tpl.ExecuteTemplate(w, "trailersManager.html", nil)
	// 	return
	// }

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		print("aa6")
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}
	print(string(b[:]))

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "trailersManager.html", template.FuncMap{"jsonData": string(b[:])})
}

func MakeAddRequest(w http.ResponseWriter, r *http.Request) {
	trailerNumber := r.FormValue("trailerNumber")
	trailerName := r.FormValue("trailerName")

	client := &http.Client{}
	params := url.Values{}
	params.Add("trailerNumber", trailerNumber)
	params.Add("trailerName", trailerName)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/trailers/add", strings.NewReader(params.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		return
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "trailersManager.html", nil)
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
	destUrl := "http://localhost:8080/trailers/manager"
	http.Redirect(w, r, destUrl, http.StatusSeeOther)

	//tpl.ExecuteTemplate(w, "index.html", cookie)

	// w.Header().Set("Content-type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// tpl.ExecuteTemplate(w, "index.html", nil)
}
