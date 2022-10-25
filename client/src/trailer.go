package src

// User is the data type for user object
type Trailer struct {
	Id     int    `sql:"id"`
	Number string `sql:"trailer_number"`
	Name   string `sql:"trailer_name"`
	UserId int    `json:"username" sql:"user_id"`
}
