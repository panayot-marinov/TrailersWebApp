package src

// using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go

// import (
// 	"fmt"
// 	"log"

// 	"github.com/sendgrid/sendgrid-go"
// 	"github.com/sendgrid/sendgrid-go/helpers/mail"
// )

// func SendMail() {
// 	from := mail.NewEmail("Example User", "panayot.marinov12@gmail.com")
// 	subject := "Sending with SendGrid is Fun"
// 	to := mail.NewEmail("Example User", "natural.medicine2k17@gmail.com")
// 	plainTextContent := "and easy to do anywhere, even with Go"
// 	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
// 	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
// 	//client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
// 	client := sendgrid.NewSendClient("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")
// 	response, err := client.Send(message)
// 	if err != nil {
// 		log.Println(err)
// 	} else {
// 		fmt.Println(response.StatusCode)
// 		fmt.Println(response.Body)
// 		fmt.Println(response.Headers)
// 	}
// }

import (
	//"github.com/d-vignesh/go-jwt-auth/utils"

	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// MailService represents the interface for our mail service.
type MailService interface {
	CreateMail(mailReq *Mail) []byte
	SendMail(mailReq *Mail) error
	NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail
}

type MailType int

// List of Mail Types we are going to send.
const (
	MailConfirmation MailType = iota + 1
	PassReset
)

// MailData represents the data to be sent to the template of the mail.
type MailData struct {
	Host     string
	Username string
	Code     string
}

// Mail represents a email request
type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}

// SGMailService is the sendgrid implementation of our MailService.
type SGMailService struct {
	logger hclog.Logger
	config Configuration
}

// NewSGMailService returns a new instance of SGMailService
func NewSGMailService(logger hclog.Logger, config Configuration) *SGMailService {
	return &SGMailService{logger, config}
}

// CreateMail takes in a mail request and constructs a sendgrid mail type.
func (ms *SGMailService) CreateMail(mailReq *Mail) []byte {

	m := mail.NewV3Mail()

	from := mail.NewEmail("Trailers project", mailReq.from)
	m.SetFrom(from)

	if mailReq.mtype == MailConfirmation {
		m.SetTemplateID(ms.config.MailApiConfig.MailVerifTemplateID)
	} else if mailReq.mtype == PassReset {
		m.SetTemplateID(ms.config.MailApiConfig.PassResetTemplateID)
	}

	p := mail.NewPersonalization()

	tos := make([]*mail.Email, 0)
	for _, to := range mailReq.to {
		tos = append(tos, mail.NewEmail("user", to))
	}

	p.AddTos(tos...)
	p.SetDynamicTemplateData("Host", mailReq.data.Host)
	p.SetDynamicTemplateData("Username", mailReq.data.Username)
	p.SetDynamicTemplateData("Code", mailReq.data.Code)
	p.Subject = "Trailers account verification"

	m.AddPersonalizations(p)
	return mail.GetRequestBody(m)
}

// SendMail creates a sendgrid mail from the given mail request and sends it.
func (ms *SGMailService) SendMail(mailReq *Mail) error {
	// from := mail.NewEmail("Example User", mailReq.from)
	// subject := "Sending with SendGrid is Fun"
	// to := mail.NewEmail("Example User", mailReq.to[0])
	// plainTextContent := "and easy to do anywhere, even with Go"
	// htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	// message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	// //client := sendgrid.NewSendClient(os.Getenv(ms.config.SendGridApiKey))
	// client := sendgrid.NewSendClient(ms.config.SendGridApiKey)
	// print("apikey =")
	// print(ms.config.SendGridApiKey)
	// response, err := client.Send(message)
	// if err != nil {
	// 	ms.logger.Error("unable to send mail", "error", err)
	// } else {
	// 	fmt.Println(response.StatusCode)
	// 	fmt.Println(response.Body)
	// 	fmt.Println(response.Headers)
	// }

	// return err

	request := sendgrid.GetRequest(ms.config.MailApiConfig.SendGridApiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = ms.CreateMail(mailReq)
	request.Body = Body
	response, err := sendgrid.API(request)
	if err != nil {
		ms.logger.Error("unable to send mail", "error", err)
		return err
	}
	ms.logger.Info("mail sent successfully", "sent status code", response.StatusCode)
	return nil

}

// NewMail returns a new mail request.
func (authHandler *AuthHandler) NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail {
	return &Mail{
		from:    from,
		to:      to,
		subject: subject,
		mtype:   mailType,
		data:    data,
	}
}

func (authHandler *AuthHandler) SendVerificationMail(db *sql.DB, user User, host string) error {
	from := "panayot.marinov12@gmail.com"
	to := []string{user.Email}
	print("to = ")
	print(to)
	print("\n")

	subject := "Email Verification for Bookite"
	mailType := MailConfirmation
	mailData := &MailData{
		Host:     host,
		Username: user.Username,
		Code:     GenerateRandomString(8),
	}
	fmt.Println("Code is " + mailData.Code)

	mailReq := authHandler.NewMail(from, to, subject, mailType, mailData)
	err := authHandler.mailService.SendMail(mailReq)
	if err != nil {
		authHandler.logger.Error("unable to send mail", "error", err)
		return err
	}

	verificationData := &VerificationData{
		Email:     user.Email,
		Code:      mailData.Code,
		Type:      MailConfirmation,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(3)),
	}

	err = StoreVerificationData(db, *verificationData)
	if err != nil {
		authHandler.logger.Error("unable to store verification data", "error", err)
		return err
	}
	return err
}

// VerifyMail verifies the provided confirmation code and set the User state to verified
func (authHandler *AuthHandler) VerifyMail(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	code := r.URL.Query().Get("code")
	print("code=")
	print(code)
	print("username=")
	print(username)

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	account, err := GetAccountInfoFromDb(db, username)
	if err != nil {
		authHandler.logger.Error("unable to fetch account data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		print(4)
		return
	}

	verificationData := &VerificationData{}
	verificationData.Email = account.Email
	verificationData.Code = code

	// err := FromJSON(verificationData, r.Body)
	// if err != nil {
	// 	//ah.logger.Error("deserialization of verification data failed", "error", err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	//data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
	// 	//w.WriteHeader(http.StatusSeeOther)
	// 	//http.Redirect(w, r, prevUrl, http.StatusBadRequest)
	// 	return
	// }

	authHandler.logger.Debug("verifying the confimation code")
	verificationData.Type = MailConfirmation

	actualVerificationData, err := GetVerificationDataFromDb(db, verificationData.Email)
	if err != nil {
		authHandler.logger.Error("unable to fetch verification data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		print(3)
		//data.ToJSON(&GenericResponse{Status: false, Message: "Unable to verify mail. Please try again later"}, w)
		return
	}

	valid, err := authHandler.Verify(db, &actualVerificationData, verificationData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		print(2)
		return
	}
	if !valid {
		w.WriteHeader(http.StatusNotAcceptable)
		print(0)
		//data.ToJSON(&GenericResponse{Status: false, Message: err.Error()}, w)
		return
	}

	// correct code, update user status to verified.
	print(account.Username)
	err = UpdateAccountVerificationStatus(db, account.Username, true)
	if err != nil {
		authHandler.logger.Error("unable to set user verification status to true")
		w.WriteHeader(http.StatusInternalServerError)
		print(1)
		//data.ToJSON(&GenericResponse{Status: false, Message: "Unable to verify mail. Please try again later"}, w)
		return
	}

	// delete the VerificationData from db
	err = DeleteVerificationDataFromDb(db, verificationData.Email, verificationData.Type)
	if err != nil {
		authHandler.logger.Error("unable to delete the verification data", "error", err)
	}

	authHandler.logger.Debug("user mail verification succeeded")

	w.WriteHeader(http.StatusAccepted)
	//data.ToJSON(&GenericResponse{Status: true, Message: "Mail Verification succeeded"}, w)
}

func (authHandler *AuthHandler) Verify(db *sql.DB, actualVerificationData *VerificationData, verificationData *VerificationData) (bool, error) {

	// check for expiration
	if actualVerificationData.ExpiresAt.Before(time.Now()) {
		authHandler.logger.Error("verification data provided is expired")
		err := DeleteVerificationDataFromDb(db, actualVerificationData.Email, actualVerificationData.Type)
		authHandler.logger.Error("unable to delete verification data from db", "error", err)
		return false, errors.New("Confirmation code has expired. Please try generating a new code")
	}

	if actualVerificationData.Code != verificationData.Code {
		authHandler.logger.Error("verification of mail failed. Invalid verification code provided")
		return false, errors.New("Verification code provided is Invalid. Please look in your mail for the code")
	}

	return true, nil
}

func (authHandler *AuthHandler) SendPasswordResetMail(db *sql.DB, user User, host string) error {
	from := "panayot.marinov12@gmail.com"
	to := []string{user.Email}

	subject := "Email Verification for Bookite"
	mailType := MailConfirmation
	mailData := &MailData{
		Host:     host,
		Username: user.Username,
		Code:     GenerateRandomString(8),
	}

	mailReq := authHandler.NewMail(from, to, subject, mailType, mailData)
	err := authHandler.mailService.SendMail(mailReq)
	if err != nil {
		authHandler.logger.Error("unable to send mail", "error", err)
		return err
	}

	verificationData := &VerificationData{
		Email:     user.Email,
		Code:      mailData.Code,
		Type:      PassReset,
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(3)),
	}

	err = StoreVerificationData(db, *verificationData)
	if err != nil {
		authHandler.logger.Error("unable to store verification data", "error", err)
		return err
	}
	return err
}

// GeneratePassResetCode generate a new secret code to reset password.
func (authHandler *AuthHandler) GeneratePasswordResetCode(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	host := r.FormValue("host")
	email := r.FormValue("email")
	fmt.Println("email = " + email)
	fmt.Println("host = " + host)

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()
	user, err := GetUserInfoFromDbWithEmail(db, email)
	if err != nil {
		authHandler.logger.Error("unable to fetch account data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		print(4)
		return
	}

	// Send verification mail
	from := "panayot.marinov12@gmail.com"
	to := []string{user.Email}
	subject := "Password Reset for trailers"
	mailType := PassReset
	mailData := &MailData{
		Host:     host,
		Username: user.Username,
		Code:     GenerateRandomString(8),
	}

	mailReq := authHandler.NewMail(from, to, subject, mailType, mailData)
	err = authHandler.mailService.SendMail(mailReq)
	if err != nil {
		authHandler.logger.Error("unable to send mail", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// store the password reset code to db
	verificationData := &VerificationData{
		Email:     user.Email,
		Code:      mailData.Code,
		Type:      PassReset,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(authHandler.mailService.config.MailApiConfig.PassResetCodeExpiration)),
	}

	err = StoreVerificationData(db, *verificationData)
	if err != nil {
		authHandler.logger.Error("unable to store verification data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		authHandler.logger.Error("unable to store password reset verification data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authHandler.logger.Debug("successfully mailed password reset code")
	w.WriteHeader(http.StatusOK)
	return
}

// VerifyPasswordReset verifies the code provided for password reset
func (authHandler *AuthHandler) VerifyPasswordReset(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	username := r.URL.Query().Get("username")
	code := r.URL.Query().Get("code")
	print("code=")
	print(code)
	print("username=")
	print(username)

	authHandler.logger.Debug("verifing password reset code")
	verificationData := VerificationData{}
	verificationData.Type = PassReset
	verificationData.Code = code

	db := ConnectToDb(authHandler.configuration.DbConfig)
	defer db.Close()

	account, err := GetAccountInfoFromDb(db, username)
	verificationData.Email = account.Email
	if err != nil {
		authHandler.logger.Error("unable to fetch account data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		print(4)
		return
	}

	actualVerificationData, err := GetVerificationDataFromDb(db, account.Email)
	if err != nil {
		authHandler.logger.Error("unable to fetch verification data", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	valid, err := authHandler.Verify(db, &actualVerificationData, &verificationData)
	if err != nil {
		authHandler.logger.Error("verification data is not valid", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !valid {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	authHandler.logger.Debug("password reset code verification succeeded")
	w.WriteHeader(http.StatusAccepted)
}

//func (authHandler *AuthHandler) SendPasswordResetEmail(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		//destUrl := "http://localhost:8080/"
//		w.WriteHeader(http.StatusUnauthorized)
//		//http.Redirect(w, r, destUrl, http.StatusAccepted)
//		return
//	}
//
//	host := r.FormValue("host")
//	email := r.FormValue("email")
//	print("email=")
//	print(email)
//
//	db := ConnectToDb(authHandler.configuration.DbConfig)
//	defer db.Close()
//	user, err := GetUserInfoFromDbWithEmail(db, email)
//	if err != nil {
//		fmt.Println("Error getting user info from db")
//		w.WriteHeader(http.StatusInternalServerError)
//		panic(err)
//	}
//
//	err = authHandler.SendVerificationMail(db, user, host)
//	if err != nil {
//		fmt.Println("Error sending verification email")
//		w.WriteHeader(http.StatusInternalServerError)
//		panic(err)
//	}
//
//	w.WriteHeader(http.StatusAccepted)
//}
