package main

import (
	"github.com/golang/glog"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, there! Running test function..")
	testLogic()
}

// User map
var UserMap = make(map[string]*Account)

func authenticateUser(userId string) {
	account := initializeAccount(userId)
	UserMap[userId] = &account
	glog.Info("Retrieving user from the db..")
}

func testLogic(){
	account := initializeAccount("123")
	add(&account, 100)
	// glog.Info("Account balance after adding: ", account.getBalance())
	glog.Info(account.Balance==100)
	buy(&account, "Apple", 10)
	// glog.Info("Available account balance after BUY: ", account.getBalance())
	glog.Info(account.Balance==100)
	glog.Info(account.Available==90)
	// glog.Info("Account balance after BUY: ", account.Balance)
	commitBuy(&account)
	glog.Info(account.Balance==90)
	glog.Info(account.Available==90)
	glog.Info(account.StockPortfolio["Apple"]==10)
	// glog.Info("Available account balance after COMMIT BUY: ", account.getBalance())
	// glog.Info("Account balance after COMMIT BUY: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	// glog.Info("Selling 5 shares of Apple")
	sell(&account, "Apple", 5)
	glog.Info(account.StockPortfolio["Apple"]==5)
	glog.Info(account.Balance==90)
	commitSell(&account)
	glog.Info(account.StockPortfolio["Apple"]==5)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==95)
	// glog.Info("Available account balance after COMMIT SELL: ", account.getBalance())
	// glog.Info("Account balance after COMMIT SELL: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	//this should fail
	sell(&account, "Apple", 100)
	commitSell(&account)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	glog.Info(account.StockPortfolio["Apple"]==5)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==95)
	buy(&account, "Apple", 10000)
	commitBuy(&account)
	glog.Info(account.StockPortfolio["Apple"]==5)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==95)
	buy(&account, "Apple", 10)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==94)
	cancelBuy(&account)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==95)
	glog.Info(account.StockPortfolio["Apple"]==5)
	glog.Info("BEFORE TRIGGERS:", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["Apple"])
	setBuyAmount(&account, "Apple", 10)
	setBuyTrigger(&account,"Apple", 0.5)
	glog.Info("AFTER TRIGGERS:", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["Apple"])
	buy(&account, "Apple", 10)
	commitBuy(&account)
	glog.Info("Balance: ", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["Apple"])
	setBuyTrigger(&account,"Apple", 5)
	glog.Info("AFTER TRIGGERS:", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["Apple"])

}

type Response struct {
	UserId string
	PriceDollars float64
	PriceCents float64
	Command string
	CommandNumber int
	Stock string
}

func getUser(userId string) *Account {
	return UserMap[userId]
}

//Parse request and return Response Object
func parseRequest(w http.ResponseWriter, r *http.Request) {
	msg := Response{} //initialize empty user

	//Parse json request body and use it to set fields on user
	//Note that user is passed as a pointer variable so that it's fields can be modified
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil{
		panic(err)
	}

	//Marshal or convert user object back to json and write to response 
	// msgJson, err := json.Marshal(msg)
	// if err != nil{
	// 	panic(err)
	// }

	// glog.Info("USERID: ", msg.UserId)
	// glog.Info("PRICE: ", msg.PriceDollars)
	// glog.Info("Command: ", msg.Command)
	// glog.Info("CommandNo: ", msg.CommandNumber)
	// glog.Info("Stock: ", msg.Stock)


	var account *Account
	if msg.Command != "authenticate" {
		account = getUser(msg.UserId)
		// glog.Info("USER: ", account.AccountNumber, account.Balance)
	}
	
	// if msg.Command == "add" {
	// 	glog.Info("YAAY adding money!!!")
	// 	add(&account, msg.PriceDollars)
	// 	glog.Info("Account balance after adding: ", account.getBalance())
	// }

	//TODO: rewrite this!!
	switch(msg.Command) {
	case "authenticate":
		authenticateUser(msg.UserId)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "add":
		add(account, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "buy":
		buy(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "commit_sell":
		commitSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "commit_buy":
		commitBuy(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "cance_buy":
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "cancel_sell":
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "set_buy_amount":
		setSellAmount(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "set_sell_amount":
		setSellAmount(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "cancel_set_buy":
		cancelSetBuy(account, msg.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "cancel_set_sell":
		cancelSetSell(account, msg.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "set_buy_trigger":
		setBuyTrigger(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "set_sell_trigger":
		setSellTrigger(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["Apple"])
	case "dumplog":
		glog.Info("SAVING XML LOG FILE")
	default: 
		panic("Oh noooo we can't process this request :(")

	}


	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response 
	// w.Write(msgJson)
	// return msg
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/test", echoString).Methods("GET")
	// router.HandleFunc("/getQuote", echoString).Methods("GET")
	router.HandleFunc("/api/authenticate", parseRequest).Methods("POST")
	router.HandleFunc("/api/add", parseRequest).Methods("POST")
	router.HandleFunc("/api/sell", parseRequest).Methods("POST")
	router.HandleFunc("/api/buy", parseRequest).Methods("POST")
	router.HandleFunc("/api/commit_sell",parseRequest).Methods("POST")
	router.HandleFunc("/api/commit_buy",parseRequest).Methods("POST")
	router.HandleFunc("/api/cancel_buy", parseRequest).Methods("POST")
	router.HandleFunc("/api/cancel_sell", parseRequest).Methods("POST")
	router.HandleFunc("/api/set_buy_amount", parseRequest).Methods("POST")
	router.HandleFunc("/api/set_sell_amount", parseRequest).Methods("POST")
	router.HandleFunc("/api/cancel_set_buy", parseRequest).Methods("POST")
	router.HandleFunc("/api/cancel_set_sell", parseRequest).Methods("POST")
	router.HandleFunc("/api/set_buy_trigger", parseRequest).Methods("POST")
	router.HandleFunc("/api/set_sell_trigger", parseRequest).Methods("POST")
	// router.HandleFunc("/api/", ).Methods("POST")
	//router.HandleFunc("/api/", ).Methods("POST")

	log.Fatal(http.ListenAndServe(":9090", router))

}
