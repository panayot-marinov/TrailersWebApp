package src

import (
	"html/template"
	"io"
	"log"
	"net/http"

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
	r.HandleFunc("/account", AccountDetails).Methods(http.MethodGet)
	r.HandleFunc("/makeChangePasswordRequest", MakeChangePasswordRequest).Methods(http.MethodPost)
	r.HandleFunc("/makeDeleteAccountRequest", MakeDeleteAccountRequest).Methods(http.MethodPost)
	r.HandleFunc("/sendPasswordResetRequest", SendPasswordResetRequestPage).Methods(http.MethodGet)
	r.HandleFunc("/makePasswordResetSendEmailRequest", MakePasswordResetSendEmailRequest).Methods(http.MethodGet)
	r.HandleFunc("/passwordReset", PasswordReset).Methods(http.MethodGet)
	r.HandleFunc("/makePasswordResetRequest", MakePasswordResetRequest).Methods(http.MethodPost)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/",
		http.FileServer(http.Dir("./src/templates/assets"))))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		http.FileServer(http.Dir("./src/templates/static"))))

	mailR := r.PathPrefix("/verify").Methods(http.MethodGet).Subrouter()
	mailR.HandleFunc("/mail", VerifyMail)
	//mailR.HandleFunc("/password-reset", VerifyPasswordReset)

	trailersR := r.PathPrefix("/trailers").Subrouter()
	trailersR.HandleFunc("/data", TrailerDataDetails).Methods(http.MethodGet)
	trailersR.HandleFunc("/manager", TrailersManager).Methods(http.MethodGet)
	trailersR.HandleFunc("/makeAddRequest", MakeAddRequest).Methods(http.MethodPost)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", r))
}

// func Get(w http.ResponseWriter, r *http.Request) {
// 	//w.Header().Set("Content-type", "application/json")
// 	w.Header().Set("Content-type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// 	//w.Write([]byte(`{"message": "get called"}`))
// 	print("In index\n")

// 	tpl.ExecuteTemplate(w, "index.html", nil) //Read about nginx
// }

func Get(w http.ResponseWriter, r *http.Request) {
	print("aa0")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081/api/v1/trailers/data", nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	print("aa1")
	cookie, err := r.Cookie("session_token")
	if err != nil {
		print("cannot get cookie\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "index.html", nil)
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
		tpl.ExecuteTemplate(w, "index.html", nil)
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
		tpl.ExecuteTemplate(w, "index.html", nil)
		return
	}

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "index.html", template.FuncMap{"jsonData": string(b[:])})
}
