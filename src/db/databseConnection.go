package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

var (
	DB *sql.DB
)

const (
	DB_URL = "localhost://127.0.0.1"
	DB_NAME = "mysql"
	DB_USER = "root"
	DB_PASSWORD = "root"
)

type User struct {
	UserId string
	Name string
	AccountNumber string
}

func InitializeDB() {
	DB = databaseConnection()
}

func databaseConnection() (db *sql.DB) {
	connectionStr := DB_NAME+":"+DB_PASSWORD+"@tcp("+DB_URL+")"
	db, err := sql.Open(DB_NAME, connectionStr)

	if err != nil {
		glog.Error("Failed to establish connection with the Quote Server.")
		glog.Error(err)
	}
	return db
}

func Close() {
	DB.Close()
}
/*
func GetUser(userId string) {
	user User
	user.UserId = userId

	res, err DB.Query("SELECT name, account_number FROM users WHERE useId = ?", 
		userId).Scan(&user.Name, &user.AccountNumber)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		return nil
	}

	return User
}

func CreateNewUser(userId string, name string, email string, address string, accId string, uuid string) {
	stmt, err := DB.Prepare("INSERT users(user_id, user_name, account_number, user_address, user_email)
				VALUES(?,?,?,?,?)")
	
	if err != nil {
		glog.Erron(err)
		return
	}

	_, err = stmt.Exec(useId, name, email, address, accId, uuid)

	if err != nil {
		glog.Error(err)
		return
	}
}
*/