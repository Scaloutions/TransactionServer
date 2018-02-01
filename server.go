package main

import (
	"github.com/golang/glog"
	"fmt"
	"log"
	"flag"
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

func usage() {
	fmt.Println("usage: example -logtostderr=true -stderrthreshold=[INFO|WARN|FATAL|ERROR] -log_dir=[string]\n")
	flag.PrintDefaults()
}


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
	glog.Info("Balance: ", account.Balance)
	glog.Info(account.Available)
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

	var account *Account
	if msg.Command != "authenticate" {
		account = getUser(msg.UserId)
	}
	
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
	flag.Usage = usage
	flag.Parse()

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
