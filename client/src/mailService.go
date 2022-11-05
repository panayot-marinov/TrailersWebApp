package src

import (
	"net/http"
)

func VerifyMail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html")

	code := r.URL.Query().Get("code")
	username := r.URL.Query().Get("username")

	print("code=")
	print(code)

	requestURL := "http://localhost:8081/api/v1/verify/mail?code=" + string(code) + "&username=" + username
	print(requestURL)
	print("\n")
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "cannotVerifyEmail.html", nil)
		return
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusAccepted {
		print("cannot call api correctly\n")
		w.WriteHeader(http.StatusOK)
		tpl.ExecuteTemplate(w, "cannotVerifyEmail.html", nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	tpl.ExecuteTemplate(w, "accountVerified.html", nil)
}

func VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "text/html")

	code := r.URL.Query().Get("code")
	username := r.URL.Query().Get("username")

	print("code=")
	print(code)

	requestURL := "http://localhost:8081/api/v1/verify/passwordReset?code=" + string(code) + "&username=" + username
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
	if err != nil || resp.StatusCode != http.StatusAccepted {
		print("cannot call api correctly\n")
		w.WriteHeader(http.StatusInternalServerError)
		tpl.ExecuteTemplate(w, "error500.html", nil)
		return
	}

	w.WriteHeader(http.StatusOK)

	tpl.ExecuteTemplate(w, "password.html", nil)
}
