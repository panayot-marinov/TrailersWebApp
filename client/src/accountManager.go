package src

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

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
	hostname := strings.Split(r.Host, ":")[0]
	fmt.Println("Sending request to " + config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ServerPort) + "/api/v1/login")
	resp, _ := http.PostForm(config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/login",
		params)

	fmt.Println("respStatus = " + resp.Status)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println("respBody = ")
	fmt.Println(string(body))
	fmt.Println("code = " + strconv.Itoa(resp.StatusCode))

	if resp.StatusCode == http.StatusUnauthorized {
		//w.WriteHeader(http.StatusFound)
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusFound {
		//w.WriteHeader(http.StatusFound)
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
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
	destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort)
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
	var user User
	user.Company = r.FormValue("company")
	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")

	params := url.Values{}
	params.Add("company", user.Company)
	params.Add("username", user.Username)
	params.Add("email", user.Email)
	params.Add("password", user.Password)

	hostname := strings.Split(r.Host, ":")[0]
	resp, _ := http.PostForm(config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/register", params)

	fmt.Println("respBody")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body) // response body is []byte
	fmt.Println(string(body))

	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "verifyYourEmail.html", template.FuncMap{"Email": user.Email})
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
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	client := &http.Client{
		Jar: jar,
	}

	hostname := strings.Split(r.Host, ":")[0]
	urlObj, _ := url.Parse(config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ServerPort) + "/api/v1/logout")
	client.Jar.SetCookies(urlObj, []*http.Cookie{cookie})
	resp, _ := client.PostForm(config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/logout",
		url.Values{})

	fmt.Println("request made successfully")
	if resp.StatusCode == http.StatusUnauthorized {
		//w.WriteHeader(http.StatusOK)
		fmt.Println("redirect 0")
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
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

func MyUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		hostname := strings.Split(r.Host, ":")[0]
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	}

	client := &http.Client{}
	hostname := strings.Split(r.Host, ":")[0]
	req, err := http.NewRequest(http.MethodGet,
		config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/userProfile", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}
	decoder := json.NewDecoder(resp.Body)

	defer resp.Body.Close()

	var user User
	err = decoder.Decode(&user)
	print("username=")
	print(user.Username)
	print("email=")
	print(user.Email)
	print("company=")
	print(user.Company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "myUserProfile.html", nil)
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

	tpl.ExecuteTemplate(w, "myUserProfile.html",
		template.FuncMap{"Username": user.Username, "Email": user.Email, "Company": user.Company})
}

func MakeChangePasswordRequest(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("password " + password + " passwordRepeated " + passwordRepeated + "\n")

	client := &http.Client{}
	params := url.Values{}
	params.Add("password", password)
	params.Add("passwordRepeated", passwordRepeated)
	hostname := strings.Split(r.Host, ":")[0]
	req, err := http.NewRequest(http.MethodPost,
		config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/changePassword", strings.NewReader(params.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
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
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	} else if resp.StatusCode != http.StatusOK {
		print("2\n")
		w.WriteHeader(http.StatusBadRequest)
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	}

	print("3\n")
	//-----
	destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/account"
	http.Redirect(w, r, destUrl, http.StatusSeeOther)

	//tpl.ExecuteTemplate(w, "index.html", cookie)

	// w.Header().Set("Content-type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// tpl.ExecuteTemplate(w, "index.html", nil)
}

func MakeDeleteAccountRequest(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	print("password " + password + " passwordRepeated " + passwordRepeated + "\n")

	client := &http.Client{}
	params := url.Values{}
	params.Add("password", password)
	params.Add("passwordRepeated", passwordRepeated)
	hostname := strings.Split(r.Host, ":")[0]
	req, err := http.NewRequest(http.MethodPost,
		config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/deleteAccount",
		strings.NewReader(params.Encode()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	req.AddCookie(cookie)
	resp, err := client.Do(req)
	if err != nil {
		print("cannot call api\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
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
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	} else if resp.StatusCode != http.StatusOK {
		print("2\n")
		w.WriteHeader(http.StatusBadRequest)
		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	}

	print("3\n")
	//-----
	destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/account"
	http.Redirect(w, r, destUrl, http.StatusSeeOther)

	//tpl.ExecuteTemplate(w, "index.html", cookie)

	// w.Header().Set("Content-type", "text/html")
	// w.WriteHeader(http.StatusOK)
	// tpl.ExecuteTemplate(w, "index.html", nil)
}

func SendPasswordResetRequestPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "passwordResetRequest.html", nil)
}

func MakePasswordResetSendEmailRequest(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html")

	email := r.URL.Query().Get("email")

	hostname := strings.Split(r.Host, ":")[0]
	requestURL :=
		config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ServerPort) + "/api/v1/generatePasswordResetCode?email=" + email
	print(requestURL)
	print("\n")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		print("1")
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusAccepted {
		print("cannot call api correctly\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "index.html", nil)
}

func PasswordReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	code := r.URL.Query().Get("code")
	username := r.URL.Query().Get("username")

	print("code=")
	print(code)

	hostname := strings.Split(r.Host, ":")[0]
	requestURL :=
		config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ServerPort) + "/api/v1/verify/passwordReset?code=" + string(code) + "&username=" + username
	print(requestURL)
	print("\n")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	resp, err := client.Do(req)
	print(resp.StatusCode)
	if err != nil || resp.StatusCode != http.StatusAccepted {
		print(1)
		print("cannot call api correctly\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	print("password verified successfully")
	w.WriteHeader(http.StatusOK)

	tpl.ExecuteTemplate(w, "passwordReset.html", template.FuncMap{"username": username, "code": code})
}

func MakePasswordResetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	username := r.FormValue("username")
	code := r.FormValue("code")
	password := r.FormValue("password")
	passwordRepeated := r.FormValue("passwordRepeated")

	params := url.Values{}
	params.Add("username", username)
	params.Add("code", code)
	params.Add("password", password)
	params.Add("passwordRepeated", passwordRepeated)

	print("bb username=")
	print(username)
	print("\n")
	//change password
	hostname := strings.Split(r.Host, ":")[0]
	resp, err := http.PostForm(config.Protocol+"://"+hostname+":"+strconv.Itoa(config.ServerPort)+"/api/v1/passwordReset", params)
	print("aa")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}
	defer resp.Body.Close()

	print("request made successfully\n")
	print("resp body")
	print(resp.Body)

	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "passwordResetSuccessfully.html", nil)
}
