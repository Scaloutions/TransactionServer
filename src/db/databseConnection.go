package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/joho/godotenv"
	"strconv"
	"errors"
	"os"
	//"api"
)

var (
	DB *sql.DB
	DB_NAME string
	DB_AUTHENTICATION string
	DB_SERVER_ADDRESS string
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

type BuyObj struct {
	UserId		string
	Stock       string
	StockAmount float64
	MoneyAmount float64
	TransactionNum int	
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

	testMode, _ := strconv.ParseBool(os.Getenv("DEV_ENVIRONMENT"))
	if testMode {
		DB_SERVER_ADDRESS = os.Getenv("DB_SERVER_ADDRESS_DEV")
	} else {
		DB_SERVER_ADDRESS = os.Getenv("DB_SERVER_ADDRESS_PROD")
	}
	
	glog.Info(DB_NAME, " ", DB_AUTHENTICATION)
}

func databaseConnection() (db *sql.DB) {
	// make sure we're accessing mysql running in a docker container
	// db, err := sql.Open("mysql", DB_AUTHENTICATION + "@tcp(172.18.0.2:3306)/" + DB_NAME)
	db, err := sql.Open("mysql", DB_AUTHENTICATION + "@tcp("+DB_SERVER_ADDRESS+")/" + DB_NAME)

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
	glog.Info("DB:\tExecuting INSERT user for:", userId, " ", name, " ", address, " ", email)
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

	glog.Info("DB:\tExecuting SELECT username for:", userId)
	err := DB.QueryRow("SELECT user_name FROM users WHERE user_id =?", userId).Scan(&user.Name)
	if err != nil {
		glog.Error("Can not find the user in the database: ", userId)
		glog.Info("Error from authentication: ", err)
		return user, errors.New("User does not exist.")
		//TODO: is there a way to return nil here?
	}

	return user, nil
}

func GetAccount(userId string) (UserAccountDB, error) {
	var account UserAccountDB
	glog.Info("DB:\tExecuting SELECT account for:", userId)
	err := DB.QueryRow("SELECT user_id, balance, available_balance FROM accounts WHERE user_id=?", userId).Scan(&account.UserId, &account.Balance, &account.Available)

	if err != nil {
		glog.Error("Can not find the user account in the database: ", userId)
		return account, errors.New("User account does not exist.")
	}

	return account, nil

}

func CreateNewAccount(userId string) {
	glog.Info("DB:\tExecuting INSERT account for:", userId)
	stmt, err := DB.Prepare("INSERT accounts(user_id, balance, available_balance) VALUES(?,?,?)")
	
	if err != nil {
		glog.Error(err)
		return
	}

	_, err = stmt.Exec(userId, 0, 0)

	if err != nil {
		glog.Error(err)
		return
	}
}

func UpdateAccountBalance(userId string, val float64) error {
	glog.Info("DB:\tUpdating balance for user: ", userId, " with value ", val)
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
	glog.Info("DB:\tExecuting UPDATE for ", userId, " available balance: ", val)
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

/*
	Updates or creates a new stock record for a user
*/
func UpdateUserStock(userId string, stock string, amount float64) error {
	glog.Info("DB:\tExecuting INSERT for ", userId, " stock: ", stock, " amount: ", amount)
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

func GetUserStockAmount(userId string, stock string) (float64, error){
	var stockAmount float64 = 0.0
	glog.Info("DB:\tExecuting SELECT amount on stock for ", userId, " and stock symbol: ", stock)
	err := DB.QueryRow("SELECT amount FROM stock WHERE user_id = ? AND symbol=?", userId, stock).Scan(&stockAmount)
	if err != nil {
		// Do not return error since it only means that the user does not have that stock so 
		// stockAmount is just zero because there is not entry in the db
		glog.Error("Can not find user stock in the database: ", userId, " ", stock)
		// return stockAmount, errors.New("User does not exist.")
	}

	return stockAmount, nil
}

func CreateNewBuy(buyObj BuyObj) error {
	glog.Info("DB:\tExecuting  CREATE BUY for user: ", buyObj.UserId)

	stmt, err := DB.Prepare("INSERT INTO buy(user_id, stock, stock_amount, money_amount, transaction_num) VALUES(?,?,?,?,?)")

	if err != nil {
		glog.Error(err)
		return err
	}

	_, err = stmt.Exec(buyObj.UserId, buyObj.Stock, buyObj.StockAmount, buyObj.MoneyAmount, buyObj.TransactionNum)

	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func GetBuy(user_id string) (BuyObj, error) {
	buyObj := BuyObj{}

	glog.Info("DB:\tExecuting SELECT BUY for", user_id)
	err := DB.QueryRow("SELECT * FROM buy WHERE user_id = ? LIMIT 1", user_id).Scan(&buyObj.UserId, &buyObj.Stock, &buyObj.StockAmount, &buyObj.MoneyAmount, &buyObj.TransactionNum)
	if err != nil {
		glog.Error("Can not find BUY in buy table for: ", user_id)
		return buyObj, errors.New("Can not find BUY in buy table.")
	}
	glog.Info("DB:\tRetrived BUY for ", user_id, " as: ", buyObj)

	return buyObj, nil
}

func DeleteBuy(user_id string) error {
	glog.Info("DB:\tExecuting  DELETE BUY for user: ", user_id)

	stmt, err := DB.Prepare("DELETE FROM buy WHERE user_id=? ORDER BY transaction_num ASC LIMIT 1")

	if err != nil {
		glog.Error(err)
		return err
	}

	_, err = stmt.Exec(user_id)

	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}

func CreateNewSell() {

}

func GetSell(){
	
}