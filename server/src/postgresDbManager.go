package src

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func ConnectToDb() *sql.DB {
	//var connStr = os.Getenv("CONNSTR")
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var dbname = os.Getenv("POSTGRES_DB")
	var user = os.Getenv("POSTGRES_USER")
	var password = os.Getenv("POSTGRES_PASSWORD")
	var connect_timeout = os.Getenv("POSTGRES_CONNECT_TIMEOUT")
	var sslmode = os.Getenv("POSTGRES_SSLMODE")

	var connStr = "user=" + user + " " +
		"password=" + password + " " +
		"host=" + host + " " +
		"port=" + port + " " +
		"dbname=" + dbname + " " +
		"connect_timeout=" + connect_timeout + " " +
		"sslmode=" + sslmode

	println("CONNSTR = " + connStr)

	db, err := sql.Open("postgres", connStr) //Only checking arguments
	if err != nil {
		print("Cannot connect to db")
		panic(err)
	}

	err = db.Ping() //Actually opening up a connection
	if err != nil {
		panic(err)
	}

	return db
}
