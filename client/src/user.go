package src

// User is the data type for user object
type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	Company  string `json:"company"`
}
