package src

import (
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
)

func ConnectToDb(dbConfig DbConfiguration) *sql.DB {
	// //var connStr = os.Getenv("CONNSTR")
	// var host = os.Getenv("POSTGRES_HOST")
	// var port = os.Getenv("POSTGRES_PORT")
	// //var dbname = os.Getenv("POSTGRES_DB")
	// var dbname = "MainDb"
	// var user = os.Getenv("POSTGRES_USER")
	// var password = os.Getenv("POSTGRES_PASSWORD")
	// var connect_timeout = os.Getenv("POSTGRES_CONNECT_TIMEOUT")
	// var sslmode = os.Getenv("POSTGRES_SSL_MODE")

	var connStr = "user=" + dbConfig.User + " " +
		"password=" + dbConfig.Pass + " " +
		"host=" + dbConfig.Host + " " +
		"port=" + strconv.Itoa(dbConfig.Port) + " " +
		"dbname=" + dbConfig.Name + " " +
		"connect_timeout=" + strconv.Itoa(dbConfig.ConnectTimeout) + " " +
		"sslmode=" + dbConfig.SslMode

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
