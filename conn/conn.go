package conn

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const (
	username = "root"
	password = "password"
	hostname = "127.0.0.1:3306"
	dbname   = "test"
)

func Connect() (*sql.DB, error) {
	fmt.Println(sql.Drivers())
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+hostname+")/"+dbname)
	if err == nil {
		fmt.Println("Connected to Database Successfully, ")
		return db, nil
	} else {
		fmt.Println("Error while creating db, ", err)
		return nil, err
	}
}
