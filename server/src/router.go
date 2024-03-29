package src

import (
	"crypto/sha256"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(config Configuration) {
	logger := NewLogger()
	mailService := NewSGMailService(logger, config)
	authHandler := NewAuthHandler(mailService, logger, config)

	r := mux.NewRouter()
	//r.HandleFunc("/", Get).Methods(http.MethodGet)
	api := r.PathPrefix("/api/v1").Subrouter()
	//api.HandleFunc("/sendData", sendData).Methods(http.MethodPost)
	api.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)
	api.HandleFunc("/userProfile", authHandler.GetUserProfileInfo).Methods(http.MethodGet)
	api.HandleFunc("/changePassword", authHandler.ChangePassword).Methods(http.MethodPost)
	api.HandleFunc("/deleteAccount", authHandler.DeleteAccount).Methods(http.MethodPost)
	api.HandleFunc("/generatePasswordResetCode", authHandler.GeneratePasswordResetCode).Methods(http.MethodPost)
	api.HandleFunc("/passwordReset", authHandler.PasswordReset).Methods(http.MethodPost)

	mailR := api.PathPrefix("/verify").Methods(http.MethodGet).Subrouter()
	mailR.HandleFunc("/mail", authHandler.VerifyMail)
	mailR.HandleFunc("/passwordReset", authHandler.VerifyPasswordReset)

	trailersR := api.PathPrefix("/trailers").Subrouter()
	trailersR.HandleFunc("/data", authHandler.GetTrailerData).Methods(http.MethodGet)
	trailersR.HandleFunc("/list", authHandler.GetTrailersList).Methods(http.MethodGet)
	trailersR.HandleFunc("/add", authHandler.Add).Methods(http.MethodPost)
	trailersR.HandleFunc("/edit", authHandler.Edit).Methods(http.MethodPost)
	trailersR.HandleFunc("/delete", authHandler.Delete).Methods(http.MethodPost)

	// api.HandleFunc("/file/{fileID}", GetFile).Methods(http.MethodGet)
	// api.HandleFunc("/searchFile", SearchFile).Methods(http.MethodGet)

	// r.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/",
	// 	http.FileServer(http.Dir("./sources/templates/styles"))))
	// r.PathPrefix("/images/").Handler(http.StripPrefix("/images/",
	// 	http.FileServer(http.Dir("./sources/templates/images"))))
	// r.PathPrefix("/api/v1/styles/").Handler(http.StripPrefix("/api/v1/styles/",
	// 	http.FileServer(http.Dir("./sources/templates/styles"))))
	// r.PathPrefix("/api/v1/images/").Handler(http.StripPrefix("/api/v1/images/",
	// 	http.FileServer(http.Dir("./sources/templates/images"))))

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8081", r))
}

// func Get(w http.ResponseWriter, r *http.Request) {
// 	//w.Header().Set("Content-type", "application/json")
// 	w.Header().Set("Content-type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// 	//w.Write([]byte(`{"message": "get called"}`))
// 	tpl.ExecuteTemplate(w, "index.html", nil) //Read about nginx
// }

// func sendData(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Redirect(w, r, "/", http.StatusSeeOther)
// 		return
// 	}

// 	db := ConnectToDb(config.)
// 	defer db.Close()
// 	text1 := r.FormValue("text1")
// 	text2 := r.FormValue("text2")

// 	query := "INSERT INTO \"Texts\" (text1, text2) VALUES ($1, $2)"
// 	_, err := db.Exec(query, text1, text2)
// 	if err != nil {
// 		fmt.Println("Error executing insert statement")
// 		panic(err)
// 	}

// 	fmt.Println("text1 =" + text1)
// 	fmt.Println("text2 =" + text2)
// }

func hashSha256(str string) []byte {
	data := []byte(str)
	hashBytes := sha256.Sum256(data)
	return hashBytes[:]
}
