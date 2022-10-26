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
	query := "SELECT id, email, is_verified, created_at, updated_at, password FROM \"Users\" WHERE username = $1"
	row := db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.Id, &user.Email, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		fmt.Println("Error executing select statement")
		return User{}, err
	}

	user.Username = username

	return user, nil
}

func GetUserInfoFromDbWithEmail(db *sql.DB, email string) (User, error) {
	query := "SELECT id, username, is_verified, created_at, updated_at, password FROM \"Users\" WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		fmt.Println("Error executing select statement")
		return User{}, err
	}

	user.Email = email

	return user, nil
}

func RegisterNewAccountToDb(db *sql.DB, user User) error {
	query := "INSERT INTO \"Users\" (username, email, password, is_verified, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6)"
	fmt.Println(user.Username)
	fmt.Println(user.Email)
	fmt.Println(user.Password)
	fmt.Println(user.IsVerified)
	fmt.Println(user.CreatedAt)
	fmt.Println(user.UpdatedAt)
	_, err := db.Exec(query, user.Username, user.Email, user.Password, user.IsVerified, user.UpdatedAt, user.CreatedAt)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
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

func GetTrailerFromDb(db *sql.DB, serialNumber string) error {
	query := "SELECT id, number, name, user_id, serial, password FROM \"Users\" WHERE email = $1"
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.Password)
	if err != nil {
		fmt.Println("Error executing select statement")
		return User{}, err
	}

	user.Email = email

	return user, nil
}

func InsertTrailerDataIntoDb(db *sql.DB, trailerData TrailerData) error {
	print(trailerData.Latt)
	print("\n")
	print(trailerData.Longt)
	print("\n")
	print(trailerData.Weight)
	print("\n")
	print(trailerData.WeightStatus)
	print("\n")
	print(trailerData.ShuntVoltage)
	print("\n")
	print(trailerData.PowerSupplyVoltage)
	print("\n")
	query := "INSERT INTO \"TrailersData\" (lattitude, longtitude, gps_time, os_time, weight, weight_status, shunt_voltage, power_supply_voltage, trailer_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := db.Exec(query, trailerData.Latt, trailerData.Longt, trailerData.GpsTime, trailerData.OsTime, trailerData.Weight,
		trailerData.WeightStatus, trailerData.ShuntVoltage, trailerData.PowerSupplyVoltage)
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

func GetTrailerDataFromDb(db *sql.DB, from time.Time, to time.Time) ([]TrailerData, error) {
	query := "SELECT lattitude, longtitude, weight, weight_status, shunt_voltage, power_supply_voltage, gps_time, os_time FROM \"TrailersData\" WHERE \"os_time\" >= $1 AND \"os_time\" <= $2"
	rows, err := db.Query(query, from, to)
	if err != nil {
		print("query err")
		print(err)
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
	query := "SELECT id, number, name, user_id FROM \"Trailers\""
	rows, err := db.Query(query)
	if err != nil {
		print("query err")
		print(err)
		return nil, err
	}
	defer rows.Close()

	var trailers []Trailer

	for rows.Next() {
		var currentTrailer Trailer
		print("row")
		err := rows.Scan(&currentTrailer.Id, &currentTrailer.Number, &currentTrailer.Name, &currentTrailer.UserId)
		if err != nil {
			fmt.Println("Error scanning row")
			print(err)
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
	query := "INSERT INTO \"Trailers\" (number, name, user_id) VALUES ($1, $2, $3)"
	fmt.Println(trailer.Number)
	fmt.Println(trailer.Name)
	fmt.Println(trailer.UserId)

	_, err := db.Exec(query, trailer.Number, trailer.Name, trailer.UserId)
	if err != nil {
		fmt.Println("Error executing insert statement")
		print(err)
		print("\n")
		return err
	}

	return nil
}
