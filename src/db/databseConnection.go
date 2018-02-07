package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

const (
	DB_NAME = "mysql"
	DB_PASSWORD = "mypassw"
)

struct User {
	UserId string,
	Name string,
	AccountNumber string
}

func databaseConnection() (db *sql.DB){
	db, err := sql.Open(DB_NAME, DB_PASSWORD)

	if err != nil {
		glog.Error(err)
		panic(err)
	}
	return db
}

func GetUser(userId string) {
	db = databaseConnection()
	user User
	user.UserId = userId

	res, error db.Query("SELECT name, accountNumber FROM users WHERE useId = ?", 
		userId).Scan(&user.Name, &user.AccountNumber)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		panic(err)
	}
	defer db.Close()

	return User
}

func CreateNewUser(userId string, name string, email string, address string, accId string, uuid string) {
	db = databaseConnection()

	stmt, err := db.Prepare("INSERT users(user_id, user_name, account_number, user_address, user_email)
				VALUES(?,?,?,?,?)")
	
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(useId, name, email, address, accId, uuid)

	if err != nil {
		panic(err)
	}
}