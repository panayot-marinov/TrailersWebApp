package src

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var config Configuration

func SetupRoutes(currentConfig Configuration) {
	config = currentConfig

	r := mux.NewRouter()
	r.HandleFunc("/", Get).Methods(http.MethodGet)
	r.HandleFunc("/login", Login).Methods(http.MethodGet)
	r.HandleFunc("/makeLoginRequest", MakeLoginRequest).Methods(http.MethodPost)
	r.HandleFunc("/register", Register).Methods(http.MethodGet)
	r.HandleFunc("/makeRegisterRequest", MakeRegisterRequest).Methods(http.MethodPost)
	r.HandleFunc("/logout", MakeLogoutRequest).Methods(http.MethodGet)
	r.HandleFunc("/myUserProfile", MyUserProfile).Methods(http.MethodGet)
	r.HandleFunc("/makeChangePasswordRequest", MakeChangePasswordRequest).Methods(http.MethodPost)
	r.HandleFunc("/makeDeleteAccountRequest", MakeDeleteAccountRequest).Methods(http.MethodPost)
	r.HandleFunc("/sendPasswordResetRequest", SendPasswordResetRequestPage).Methods(http.MethodGet)
	r.HandleFunc("/makePasswordResetSendEmailRequest", MakePasswordResetSendEmailRequest).Methods(http.MethodGet)
	r.HandleFunc("/passwordReset", PasswordReset).Methods(http.MethodGet)
	r.HandleFunc("/makePasswordResetRequest", MakePasswordResetRequest).Methods(http.MethodPost)
	r.HandleFunc("/vf", VerifyMail).Methods(http.MethodGet)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./src/templates/assets"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./src/templates/static"))))

	mailR := r.PathPrefix("/verify").Subrouter()
	mailR.HandleFunc("/mail", VerifyMail).Methods(http.MethodGet)
	//mailR.HandleFunc("/password-reset", VerifyPasswordReset)

	r.PathPrefix("/verify/assets/").Handler(http.StripPrefix("/verify/assets/",
		http.FileServer(http.Dir("./src/templates/assets"))))
	r.PathPrefix("/verify/static/").Handler(http.StripPrefix("/verify/static/",
		http.FileServer(http.Dir("./src/templates/static"))))

	trailersR := r.PathPrefix("/trailers").Subrouter()
	trailersR.HandleFunc("/data", TrailerDataDetails).Methods(http.MethodGet)
	trailersR.HandleFunc("/add", TrailersAdd).Methods(http.MethodGet)
	trailersR.HandleFunc("/manager", TrailersManager).Methods(http.MethodGet)
	trailersR.HandleFunc("/makeAddRequest", MakeAddRequest).Methods(http.MethodPost)
	trailersR.HandleFunc("/makeEditRequest", MakeEditRequest).Methods(http.MethodPost)
	trailersR.HandleFunc("/makeDeleteRequest", MakeDeleteRequest).Methods(http.MethodPost)

	r.PathPrefix("/trailers/assets/").Handler(http.StripPrefix("/trailers/assets/",
		http.FileServer(http.Dir("./src/templates/assets"))))
	r.PathPrefix("/trailers/static/").Handler(http.StripPrefix("/trailers/static/",
		http.FileServer(http.Dir("./src/templates/static"))))

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.ClientPort), r))
}

func Get(w http.ResponseWriter, r *http.Request) {
	print("aa1")
	cookie, err := r.Cookie("session_token")

	hostname := strings.Split(r.Host, ":")[0]
	if err != nil {
		print("cannot get cookie\n")

		destUrl := config.Protocol + "://" + hostname + ":" + strconv.Itoa(config.ClientPort) + "/login"
		fmt.Println("desturl = " + destUrl)
		http.Redirect(w, r, destUrl, http.StatusFound)
		return
	}

	print("aa0")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet,
		config.Protocol+"://"+config.ServerHost+":"+strconv.Itoa(config.ServerPort)+"/api/v1/trailers/data",
		nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
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
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}
	print("aa3")

	defer resp.Body.Close()

	trailersData, err := io.ReadAll(resp.Body)
	if err != nil {
		print("aa6")
		w.WriteHeader(http.StatusBadRequest)
		tpl.ExecuteTemplate(w, "index.html", templatesData)
		return
	}

	req, err = http.NewRequest(http.MethodGet,
		config.Protocol+"://"+config.ServerHost+":"+strconv.Itoa(config.ServerPort)+"/api/v1/userProfile",
		nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	req.AddCookie(cookie)
	resp, err = client.Do(req)
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

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "index.html",
		template.FuncMap{"jsonTrailersData": string(trailersData[:]),
			"Username": user.Username, "Email": user.Email, "Company": user.Company})
}
