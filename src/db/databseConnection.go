package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"errors"
	"os"
)

var (
	DB *sql.DB
	DB_NAME string
	DB_AUTHENTICATION string
)

type User struct {
	UserId string
	Name string
}

type UserAccountDB struct {
	UserId string
	Balance        float64
	Available      float64
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

func GetUser(userId string) (User, error) {
	user := User { UserId: userId }

	err := DB.QueryRow("SELECT name FROM users WHERE use_id = ?", userId).Scan(&user.Name, &user.UserId)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		return user, errors.New("User does not exist.")
		//TODO: is there a way to return nil here?
	}

	return user, nil
}

func GetAccount(userId string) (UserAccountDB, error) {
	account := UserAccountDB{}
	err := DB.QueryRow("SELECT user_id, balance, available_balance FROM accounts").Scan(&account.UserId, &account.Balance, &account.Available)

	if err != nil {
		glog.Error("Can not find the user account in the database: ", userId)
		return account, errors.New("User account does not exist.")
	}

	return account, nil

}

func UpdateAccountBalance(userId string, val float64) error {
	stmt, err := DB.Prepare("UPDATE accounts SET balance=? where user_id =?")

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot create an update query")
	}

	_, err = stmt.Exec(val, userId)

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot execute an update query")
	}

	return nil
}

func UpdateAvailableAccountBalance(userId string, val float64) error {
	stmt, err := DB.Prepare("UPDATE accounts SET available_balance=? where user_id =?")

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot create an update query")
	}

	_, err = stmt.Exec(val, userId)

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot execute an update query")
	}

	return nil

}

func UpdateUserStock(userId string, stock string, amount float64) error {
	stmt, err := DB.Prepare("INSERT INTO stock(user_id, symbol, amount) VALUES(?,?,?) ON DUPLICATE KEY UPDATE amount=?")

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot create an update stock query")
	}

	_, err = stmt.Exec(userId, stock, amount, amount)

	if err != nil {
		glog.Error(err, " ", userId)
		return errors.New("Cannot execute an update stock query")
	}

	return nil
}
