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

func TrailersAdd(w http.ResponseWriter, r *http.Request) {
	print("aa1")
	_, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusSeeOther)
		return
	}

	print("before template, there is a cookie")
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "trailersAdd.html", nil)
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
		//w.WriteHeader(http.StatusInternalServerError)
		//tpl.ExecuteTemplate(w, "trailersManager.html", nil)
		//return
		destUrl := "http://localhost:8080/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
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
	fmt.Println("Making add request...")
	brand := r.FormValue("brand")
	model := r.FormValue("model")
	name := r.FormValue("nameField")
	registrationPlate := r.FormValue("registrationPlate")
	serialNumber := r.FormValue("serialNumber")
	city := r.FormValue("city")
	area := r.FormValue("area")
	addressLine := r.FormValue("addressLine")
	zipCode := r.FormValue("zipCode")

	client := &http.Client{}
	params := url.Values{}
	params.Add("brand", brand)
	params.Add("model", model)
	params.Add("nameField", name)
	params.Add("registrationPlate", registrationPlate)
	params.Add("serialNumber", serialNumber)
	params.Add("city", city)
	params.Add("area", area)
	params.Add("addressLine", addressLine)
	params.Add("zipCode", zipCode)

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

func MakeEditRequest(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("editNameModal")
	brand := r.FormValue("editBrandModal")
	model := r.FormValue("editModelModal")
	registrationPlate := r.FormValue("editRegPlateModal")
	serialNumber := r.FormValue("editSerialNumberModal")
	isActive := r.FormValue("editIsActiveSwitch")
	if isActive == "on" {
		isActive = "true"
	} else {
		isActive = "false"
	}
	city := r.FormValue("editCityModal")
	area := r.FormValue("editAreaModal")
	addressLine := r.FormValue("addressLine")
	zipCode := r.FormValue("zipCode")

	fmt.Println()
	fmt.Println(name)
	fmt.Println(brand)
	fmt.Println(model)
	fmt.Println(registrationPlate)
	fmt.Println(serialNumber)
	fmt.Println(isActive)
	fmt.Println(city)
	fmt.Println(area)
	fmt.Println(addressLine)
	fmt.Println(zipCode)

	client := &http.Client{}
	params := url.Values{}
	params.Add("nameField", name)
	params.Add("brand", brand)
	params.Add("model", model)
	params.Add("registrationPlate", registrationPlate)
	params.Add("serialNumber", serialNumber)
	params.Add("isActive", isActive)
	params.Add("city", city)
	params.Add("area", area)
	params.Add("addressLine", addressLine)
	params.Add("zipCode", zipCode)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/trailers/edit", strings.NewReader(params.Encode()))
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
}

func MakeDeleteRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Cannot parse form")
		fmt.Println(err)
	}

	fmt.Println("deleting trailers...")

	client := &http.Client{}
	params := url.Values{}

	var checkboxList []string

	for key, values := range r.Form {
		if strings.HasPrefix(key, "trailerCheckbox") {
			if values[0] == "on" {
				checkboxID := strings.TrimPrefix(key, "trailerCheckbox")
				checkboxList = append(checkboxList, checkboxID)
				params.Add("trailerCheckbox", checkboxID)
			}
		}
	}
	params.Add("a", "b")

	for _, element := range checkboxList {
		fmt.Println(element)
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:8081/api/v1/trailers/delete", strings.NewReader(params.Encode()))
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
}
