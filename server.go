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
	glog.Info("##### Account Balance: ", account.Balance, " Available: ", account.Available)
	// glog.Info("Account Stocks: ", account.StockPortfolio["S"])
	glog.Info("INFO: Retrieving user from the db..")
}

func testLogic(){
	account := initializeAccount("123")
	add(&account, 100)
	/*
		Test core commands
	*/
	// glog.Info(account.Balance==100)
	// buy(&account, "S", 10)
	// glog.Info(account.Balance==100)
	// glog.Info(account.Available==90)
	// commitBuy(&account)
	// glog.Info(account.Balance==90)
	// glog.Info(account.Available==90)
	// glog.Info(account.StockPortfolio["S"]==10)
	// sell(&account, "S", 5)
	// glog.Info(account.StockPortfolio["S"]==5)
	// glog.Info(account.Balance==90)
	// glog.Info(account.Available==90)
	// commitSell(&account)
	// glog.Info(account.StockPortfolio["S"]==5)
	// glog.Info(account.Balance==95)
	// glog.Info(account.Available==95)
	// //this should fail
	// sell(&account, "S", 100)
	// commitSell(&account)
	// glog.Info(account.StockPortfolio["S"]==5)
	// glog.Info(account.Balance==95)
	// glog.Info(account.Available==95)
	// buy(&account, "S", 10000)
	// commitBuy(&account)
	// glog.Info(account.StockPortfolio["S"]==5)
	// glog.Info(account.Balance==95)
	// glog.Info(account.Available==95)
	// buy(&account, "S", 10)
	// glog.Info(account.Balance==95)
	// glog.Info(account.Available==85)
	// cancelBuy(&account)
	// glog.Info(account.Balance==95)
	// glog.Info(account.Available==95)
	// glog.Info(account.StockPortfolio["S"]==5)
	// /*
	//     Test Buy Triggers
	// */
	// glog.Info("BEFORE TRIGGERS:", account.Balance, " Available: ", account.Available)
	// glog.Info(account.StockPortfolio["S"])
	// setBuyAmount(&account, "S", 10)
	// setBuyTrigger(&account,"S", 0.5)
	// glog.Info("AFTER TRIGGERS:", account.Balance)
	// glog.Info(account.Available)
	// glog.Info(account.StockPortfolio["S"])
	// buy(&account, "S", 10)
	// commitBuy(&account)
	// glog.Info("Balance: ", account.Balance)
	// glog.Info(account.Available)
	// glog.Info(account.StockPortfolio["S"])
	// setBuyAmount(&account, "S", 5)
	// setBuyTrigger(&account,"S", 5)
	// glog.Info("AFTER TRIGGERS:", account.Balance)
	// glog.Info(account.Available)
	// glog.Info(account.StockPortfolio["S"])
	/*
	    Test Sell Triggers
	*/
	buy(&account, "S", 50)
	commitBuy(&account)
	glog.Info("TESTING SELL TRIGGER")
	setSellAmount(&account, "S", 10)
	setSellTrigger(&account,"S", 5)
	glog.Info("AFTER TRIGGERS:", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["S"])
	buy(&account, "S", 10)
	commitBuy(&account)
	glog.Info("Balance: ", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["S"])
	setSellAmount(&account, "S", 5)
	setSellTrigger(&account,"S", 1)
	glog.Info("AFTER TRIGGERS:", account.Balance)
	glog.Info(account.Available)
	glog.Info(account.StockPortfolio["S"])

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

// func loggingMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         log.Println(r.RequestURI)
//         next.ServeHTTP(w, r)
//     })
// }

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


	var account *Account
	if msg.Command != "authenticate" {
		account = getUser(msg.UserId)
		// glog.Info("USER: ", account.AccountNumber, account.Balance)
	}
	
	//TODO: rewrite this!!
	switch(msg.Command) {
	case "authenticate":
		glog.Info("\n\n############################### INFO: Executing authenticate... ", msg.CommandNumber)
		authenticateUser(msg.UserId)
		glog.Info("\n############################### SUCCESS: Authentication Successful!")
	case "add":
		glog.Info("\n\n############################### INFO: Executing ADD FOR... ", msg.PriceDollars, msg.CommandNumber)
		add(account, msg.PriceDollars)
		glog.Info("SUCCESS: Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("SUCCESS: Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: ADD Successful \n")
	case "buy":
		glog.Info("\n\n############################### INFO: Executing BUY FOR... ", msg.PriceDollars, msg.CommandNumber)
		buy(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: BUY Successful")
	case "commit_sell":
		glog.Info("\n\n############################### INFO: Executing COMMIT SELL ", msg.CommandNumber)
		commitSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: COMMIT SELL Successful")
	case "commit_buy":
		glog.Info("\n\n############################### INFO: Executing COMMIT BUY ", msg.CommandNumber)
		commitBuy(account)
		glog.Info("\n############################### SUCCESS: COMMIT BUY Successful")
		// glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		// glog.Info("Account Stocks: ", account.StockPortfolio["S"])
	case "cancel_buy":
		glog.Info("\n\n############################### INFO: Executing CANCEL BUY ", msg.CommandNumber)
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: CANCEL BUY Successful")
	case "cancel_sell":
		glog.Info("\n\n############################### INFO: Executing CANCEL SELL ", msg.CommandNumber)
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: CANCEL SELL Successful")
	case "set_buy_amount":
		glog.Info("\n\n############################### INFO: Executing SET BUY AMOUNT ", msg.CommandNumber)
		setSellAmount(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: SET BUY AMOUNT Successful")
	case "set_sell_amount":
		glog.Info("\n\n############################### INFO: Executing SET SELL AMOUNT ", msg.CommandNumber)
		setSellAmount(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: SET SELL AMOUNT Successful")
	case "cancel_set_buy":
		glog.Info("\n\n############################### INFO: Executing CANCEL SET BUY ", msg.CommandNumber)
		cancelSetBuy(account, msg.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: CANCEL SET BUY Successful")
	case "cancel_set_sell":
		glog.Info("\n\n############################### INFO: Executing CANCEL SET SELL ", msg.CommandNumber)
		cancelSetSell(account, msg.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: CANCEL SET SELL Successful")
	case "set_buy_trigger":
		glog.Info("\n\n############################### INFO: Executing SET BUY TRIGGER ", msg.CommandNumber)
		setBuyTrigger(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: SET BUY TRIGGER Successful")
	case "set_sell_trigger":
		glog.Info("\n\n############################### INFO: Executing SET SELL TRIGGER ", msg.CommandNumber)
		setSellTrigger(account, msg.Stock, msg.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: SET SELL TRIGGER Successful")
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

	// router.Use(loggingMiddleware)

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
