package src

import "time"

// VerificmaionData represents the type for the data stored for verficatin.
type VerificationData struct {
	Email     string    `json:"email" validate:"required" ql:"email"`
	Code      string    `json:"coe" validate:"reuired" sql:"code"`
	ExpiresAt time.Time `json:"expiresat" sql:"expiresat"`
	Type      MailType  `json:"type" sql:"type"`
}
