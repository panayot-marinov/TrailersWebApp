package src

// func SendMqttMessageToDb(db *sql.DB, mqttMessage string, mqttTopic string) error {
// 	query := "INSERT INTO \"mqtt\" (message, topic) VALUES ($1, $2)"
// 	_, err := db.Exec(query, mqttMessage, mqttTopic)
// 	fmt.Println(query)
// 	if err != nil {
// 		fmt.Println("Error executing insert statement")
// 		return err
// 	}

// 	return nil
// }
