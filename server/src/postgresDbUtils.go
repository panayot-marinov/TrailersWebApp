package src

import (
	"database/sql"
	"fmt"
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
	return account, nil
}

func RegisterNewAccountToDb(db *sql.DB, username string, email string, password string, hashedPassword string) error {
	query := "INSERT INTO \"Users\" (username, email, password) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, username, email, hashedPassword)
	if err != nil {
		fmt.Println("Error executing insert statement")
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
