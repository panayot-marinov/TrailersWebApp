package src

// User is the data type for user object
type Trailer struct {
	Id     int    `sql:"id"`
	Number string `sql:"number"`
	Name   string `sql:"name"`
	UserId int    `sql:"user_id"`
}
