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

func GetUser(userId string) (User, error) {
	user := User { UserId: userId }

	err := DB.QueryRow("SELECT name FROM users WHERE useId = ?", userId).Scan(&user.Name, &user.AccountNumber)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		return user, errors.New("User does not exist.")
		//TODO: is there a way to return nil here?
	}

	return user, nil
}

func CreateNewUser(userId string, name string, email string, address string) {
	stmt, err := DB.Prepare("INSERT users(user_id, user_name, user_address, user_email) VALUES(?,?,?,?)")
	
	if err != nil {
		glog.Error(err)
		return
	}

	_, err = stmt.Exec(userId, name, address, email)

	if err != nil {
		glog.Error(err)
		return
	}
}
