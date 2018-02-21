package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"os"
	"api"
)

var (
	DB *sql.DB
	DB_NAME string
	DB_AUTHENTICATION string
)

// const (
// 	// DB_URL = "localhost://127.0.0.1"
// 	DB_NAME = "DAYTRADING"
// 	DB_AUTHENTICATION = "root:root"
// )

type User struct {
	UserId string
	Name string
	AccountNumber string
}

func InitializeDB() {
	loadCredentials()
	DB = databaseConnection()
}

func loadCredentials() {
	err := godotenv.Load()
 	if err != nil {
 	   glog.Error("Error loading .env file")
	}
	
	DB_NAME = os.Getenv("DB_NAME")
	DB_AUTHENTICATION = os.Getenv("DB_USER_NAME") + ":" + os.Getenv("DB_PASSWORD")
	glog.Info(DB_NAME, " ", DB_AUTHENTICATION)
}

func databaseConnection() (db *sql.DB) {
	// connectionStr := DB_NAME+":"+DB_PASSWORD+"@tcp("+DB_URL+")"
	// db, err := sql.Open(DB_NAME, connectionStr)
	db, err := sql.Open("mysql", DB_AUTHENTICATION + "@/" + DB_NAME)

	if err != nil {
		glog.Error("Failed to establish connection with the Quote Server.")
		glog.Error(err)
	}
	return db
}

func Close() {
	DB.Close()
}

func GetUser(userId string) User {
	user := User { UserId: userId }

	err := DB.QueryRow("SELECT name, account_number FROM users WHERE useId = ?", userId).Scan(&user.Name, &user.AccountNumber)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		//TODO: is there a way to return nil here?
	}

	return user
}

func getAccount(accId string) api.Account {
	account := api.Account { AccountNumber: accId }

	//TODO: pull from the db
	return account
}

func CreateNewUser(userId string, name string, email string, address string, accId string) {
	stmt, err := DB.Prepare("INSERT users(user_id, user_name, account_number, user_address, user_email) VALUES(?,?,?,?,?)")
	
	if err != nil {
		glog.Error(err)
		return
	}

	_, err = stmt.Exec(userId, name, accId, address, email)

	if err != nil {
		glog.Error(err)
		return
	}
}
