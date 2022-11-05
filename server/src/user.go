package src

import "time"

// User is the data type for user object
type User struct {
	Id         int       `json:"id" sql:"id"`
	Email      string    `json:"email" validate:"required" sql:"email"`
	Password   []byte    `json:"password" validate:"required" sql:"password"`
	Username   string    `json:"username" sql:"username"`
	Company    string    `json:"company" sql:"company"`
	IsVerified bool      `json:"isverified" sql:"is_verified"`
	CreatedAt  time.Time `json:"createdat" sql:"created_at"`
	UpdatedAt  time.Time `json:"updatedat" sql:"updated_at"`
}
