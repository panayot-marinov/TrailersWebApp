package src

import (
	"database/sql"
	"fmt"
	"time"
)

func GetAccountInfoFromDb(db *sql.DB, username string) (Account, error) {
	query := "SELECT email FROM \"Users\" WHERE username = $1"
	print("d0\n")
	print("u=")
	print(username)
	print("\n")
	row := db.QueryRow(query, username)

	var account Account
	err := row.Scan(&account.Email)
	if err != nil {
		fmt.Println("Error executing select statement")
		return Account{}, err
	}

	print("email=")
	print(account.Email)
	account.Username = username

	return account, nil
}

func GetUserInfoFromDbWithUsername(db *sql.DB, username string) (User, error) {
	query := "SELECT id, email, company, is_verified, created_at, updated_at, password FROM \"Users\" WHERE username = $1"
	row := db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.Id, &user.Email, &user.Company, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		fmt.Println("Error executing select statement")
		return User{}, err
	}

	user.Username = username

	return user, nil
}

func GetUserInfoFromDbWithEmail(db *sql.DB, email string) (User, error) {
	query := "SELECT id, username, company, is_verified, created_at, updated_at, password FROM \"Users\" WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Company, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		fmt.Println("Error executing select statement")
		return User{}, err
	}

	user.Email = email

	return user, nil
}

func RegisterNewAccountToDb(db *sql.DB, user User) error {
	query := "INSERT INTO \"Users\" (username, email, password, company, is_verified, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.Exec(query, user.Username, user.Email, user.Password, user.Company, user.IsVerified, user.UpdatedAt, user.CreatedAt)
	if err != nil {
		fmt.Println("Error executing insert statement")
		fmt.Println(err)
		print("\n")
		return err
	}

	return nil
}

func UpdateAccountPasswordToDb(db *sql.DB, username string, password string) error {
	query := "UPDATE \"Users\" SET password=$1 WHERE username=$2"
	_, err := db.Exec(query, password, username)
	fmt.Println(query)

	if err != nil {
		fmt.Println("Error executing insert statement")
		return err
	}

	return nil
}

func DeleteAccountFromDb(db *sql.DB, username string) error {
	query := "DELETE FROM \"Users\" WHERE username=$1"
	_, err := db.Exec(query, username)
	fmt.Println(query)

	if err != nil {
		fmt.Println("Error executing DELETE statement")
		return err
	}

	return nil
}

func SendMqttMessageToDb(db *sql.DB, mqttMessage string, mqttTopic string) error {
	query := "INSERT INTO \"mqtt\" (message, topic) VALUES ($1, $2)"
	_, err := db.Exec(query, mqttMessage, mqttTopic)
	fmt.Println(query)
	if err != nil {
		fmt.Println("Error executing insert statement")
		return err
	}

	return nil
}

// func GetTrailersDataFromDb(db *sql.DB, username string) (Account, error) {
// 	query := "SELECT * FROM \"TrailersData\" WHERE username = $1"
// 	print("d0\n")
// 	print("u=")
// 	print(username)
// 	print("\n")
// 	row := db.QueryRow(query, username)

// 	var account Account
// 	err := row.Scan(&account.Email)
// 	if err != nil {
// 		fmt.Println("Error executing select statement")
// 		return Account{}, err
// 	}

// 	print("email=")
// 	print(account.Email)
// 	return account, nil
// }

func GetTrailerRegPlateBySerialNumberFromDb(db *sql.DB, serialNumber string) (string, error) {
	query := "SELECT registration_plate FROM \"Trailers\" WHERE serial_number = $1"
	row := db.QueryRow(query, serialNumber)

	var registrationPlate string
	err := row.Scan(&registrationPlate)
	if err != nil {
		fmt.Println("Error executing select statement")
		return "", err
	}

	return registrationPlate, nil
}

func InsertTrailerDataIntoDb(db *sql.DB, trailerData TrailerData) error {
	registrationPlate, err := GetTrailerRegPlateBySerialNumberFromDb(db, trailerData.SerialNumber)
	if err != nil {
		fmt.Println("Error getting registration plate by serial number")
		return err
	}

	query := "INSERT INTO \"TrailersData\" (lattitude, longtitude, gps_time, os_time, weight, weight_status, shunt_voltage, power_supply_voltage, serial_number, cpu_temp, registration_plate) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err = db.Exec(query, trailerData.Latt, trailerData.Longt, trailerData.GpsTime, trailerData.OsTime, trailerData.Weight,
		trailerData.WeightStatus, trailerData.ShuntVoltage, trailerData.PowerSupplyVoltage,
		trailerData.SerialNumber, trailerData.CpuTemp, registrationPlate)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error executing insert statement")
		return err
	}

	return nil
}

func GetVerificationDataFromDb(db *sql.DB, email string) (VerificationData, error) {
	query := "SELECT * FROM \"Verifications\" WHERE email = $1"
	row := db.QueryRow(query, email)

	var verificationData VerificationData
	err := row.Scan(&verificationData.Email, &verificationData.Code, &verificationData.ExpiresAt, &verificationData.Type)
	if err != nil {
		fmt.Println("Error executing select statement")
		return VerificationData{}, err
	}

	return verificationData, nil
}

func StoreVerificationData(db *sql.DB, verificationData VerificationData) error {
	query := "INSERT INTO \"Verifications\" (email, code, expiresAt, type) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(query, verificationData.Email, verificationData.Code, verificationData.ExpiresAt, verificationData.Type)
	fmt.Println(query)
	if err != nil {
		fmt.Println("Error executing insert statement")
		return err
	}

	return nil
}

func DeleteVerificationDataFromDb(db *sql.DB, email string, verificationType MailType) error {
	query := "DELETE FROM \"Verifications\" WHERE email=$1 AND type=$2"
	_, err := db.Exec(query, email, verificationType)
	fmt.Println(query)

	if err != nil {
		fmt.Println("Error executing DELETE statement")
		return err
	}

	return nil
}

func UpdateAccountVerificationStatus(db *sql.DB, username string, isVerified bool) error {
	query := "UPDATE \"Users\" SET is_verified=$1 WHERE username=$2"
	_, err := db.Exec(query, isVerified, username)

	if err != nil {
		fmt.Println("Error executing update statement")
		return err
	}

	return nil
}

func GetTrailersDataFromDb(db *sql.DB, from time.Time, to time.Time) ([]TrailerData, error) {
	query := "SELECT lattitude, longtitude, weight, weight_status, shunt_voltage, power_supply_voltage, gps_time, os_time FROM \"TrailersData\" WHERE \"os_time\" >= $1 AND \"os_time\" <= $2"
	rows, err := db.Query(query, from, to)
	if err != nil {
		print("query err")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var trailerData []TrailerData

	for rows.Next() {
		var currentTailerData TrailerData
		print("row")
		err := rows.Scan(&currentTailerData.Latt, &currentTailerData.Longt, &currentTailerData.Weight,
			&currentTailerData.WeightStatus, &currentTailerData.ShuntVoltage, &currentTailerData.PowerSupplyVoltage,
			&currentTailerData.GpsTime, &currentTailerData.OsTime)
		if err != nil {
			fmt.Println("Error scanning row")
			return trailerData, err
		}

		trailerData = append(trailerData, currentTailerData)
	}
	if err = rows.Err(); err != nil {
		return trailerData, err
	}

	return trailerData, nil
}

func GetTrailerDataFromDb(db *sql.DB, from time.Time, to time.Time, registrationPlate string) ([]TrailerData, error) {
	query := "SELECT lattitude, longtitude, weight, weight_status, shunt_voltage, power_supply_voltage, gps_time, os_time FROM \"TrailersData\" WHERE \"registration_plate\" = $1 AND \"os_time\" >= $2 AND \"os_time\" <= $3"
	rows, err := db.Query(query, registrationPlate, from, to)
	if err != nil {
		print("query err")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var trailerData []TrailerData

	for rows.Next() {
		var currentTailerData TrailerData
		print("row")
		err := rows.Scan(&currentTailerData.Latt, &currentTailerData.Longt, &currentTailerData.Weight,
			&currentTailerData.WeightStatus, &currentTailerData.ShuntVoltage, &currentTailerData.PowerSupplyVoltage,
			&currentTailerData.GpsTime, &currentTailerData.OsTime)
		if err != nil {
			fmt.Println("Error scanning row")
			return trailerData, err
		}

		trailerData = append(trailerData, currentTailerData)
	}
	if err = rows.Err(); err != nil {
		return trailerData, err
	}

	return trailerData, nil
}

func GetTrailersListFromDb(db *sql.DB) ([]Trailer, error) {
	query := "SELECT registration_plate, name, user_id, serial_number, brand, model, city, area, address_line, zip_code, is_active FROM \"Trailers\""
	rows, err := db.Query(query)
	if err != nil {
		print("query err")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var trailers []Trailer

	for rows.Next() {
		var currentTrailer Trailer
		print("row")
		err := rows.Scan(&currentTrailer.RegistrationPlate, &currentTrailer.Name, &currentTrailer.UserId,
			&currentTrailer.SerialNumber, &currentTrailer.Brand, &currentTrailer.Model,
			&currentTrailer.City, &currentTrailer.Area, &currentTrailer.AddressLine, &currentTrailer.ZipCode,
			&currentTrailer.IsActive)
		if err != nil {
			fmt.Println("Error scanning row")
			fmt.Println(err)
			return trailers, err
		}

		trailers = append(trailers, currentTrailer)
	}
	if err = rows.Err(); err != nil {
		return trailers, err
	}

	return trailers, nil
}

func RegisterNewTrailerToDb(db *sql.DB, trailer Trailer) error {
	//query := "INSERT INTO \"Trailers\" (registration_plate, name, user_id) VALUES ($1, $2, $3)"
	query := "INSERT INTO \"Trailers\" (user_id, brand, model, name, registration_plate, serial_number, city, area, address_line, zip_code) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

	print("user_id")
	print(trailer.UserId)
	//_, err := db.Exec(query, trailer.RegistrationPlate, trailer.Name, trailer.UserId)
	_, err := db.Exec(query, trailer.UserId, trailer.Brand, trailer.Model, trailer.Name, trailer.RegistrationPlate,
		trailer.SerialNumber, trailer.City, trailer.Area, trailer.AddressLine, trailer.ZipCode)
	if err != nil {
		fmt.Println("Error executing insert statement")
		fmt.Println(err)
		print("\n")
		return err
	}
	fmt.Println("Successfully inserted into db")

	return nil
}

func UpdateTrailerIntoDb(db *sql.DB, trailer Trailer) error {
	fmt.Println("Updating trailer")
	fmt.Println("Name = " + trailer.Name)
	query := "UPDATE \"Trailers\" SET name=$1, brand=$2, model=$3, serial_number=$4, city=$5, area=$6, address_line=$7, zip_code=$8, is_active=$9 WHERE registration_plate=$10"
	_, err := db.Exec(query, trailer.Name, trailer.Brand, trailer.Model, trailer.SerialNumber, trailer.City, trailer.Area, trailer.AddressLine, trailer.ZipCode, trailer.IsActive, trailer.RegistrationPlate)

	if err != nil {
		fmt.Println("Error executing update statement")
		fmt.Println(err)
		return err
	}

	return nil
}

func DeleteTrailerFromDb(db *sql.DB, registrationPlate string) error {
	fmt.Println("Deleting trailer")
	query := "DELETE FROM \"Trailers\" WHERE registration_plate=$1"
	_, err := db.Exec(query, registrationPlate)

	if err != nil {
		fmt.Println("Error executing delete statement")
		fmt.Println(err)
		return err
	}

	return nil
}

func newNullString(s string) sql.NullString {
	if len(s) == 0 {
		fmt.Println("returning null string")
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
