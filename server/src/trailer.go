package src

// Trailer is the data type for user object
type Trailer struct {
	Brand             string `sql:"brand"`
	Model             string `sql:"model"`
	RegistrationPlate string `sql:"registration_plate"`
	Name              string `sql:"name"`
	City              string `sql:"city"`
	Area              string `sql:"area"`
	AddressLine       string `sql:"address_line"`
	UserId            int    `sql:"user_id"`
	SerialNumber      string `sql:"serial_number"`
	ZipCode           int    `sql:"zip_code"`
	IsActive          bool   `sql:"active"`
}
