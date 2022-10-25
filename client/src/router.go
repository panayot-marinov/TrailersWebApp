package src

import (
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

func Get(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-type", "application/json")
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(`{"message": "get called"}`))
	print("In index\n")

	tpl.ExecuteTemplate(w, "index.html", nil) //Read about nginx
}
