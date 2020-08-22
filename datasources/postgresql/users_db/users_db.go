package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	postgresHost     = "postgresHost"
	postgresPort     = "postgresPort"
	postgresUser     = "postgresUser"
	postgresPassword = "postgresPassword"
	postgresDbname   = "postgresDbname"
)

var (
	Client *sql.DB

	host     = os.Getenv(postgresHost)
	port     = os.Getenv(postgresPort)
	user     = os.Getenv(postgresUser)
	password = os.Getenv(postgresPassword)
	dbname   = os.Getenv(postgresDbname)
)

func init() {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	Client, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}
	log.Println("Successfully connected!")
}
