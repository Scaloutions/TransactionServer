package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"
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
>>>>>>> origin/ParseHttpReqs
	
	f := getFilePointer()

	account := initializeAccount("123")
	add(&account, 100, f)
	// glog.Info("Account balance after adding: ", account.getBalance())
	glog.Info(account.Balance == 100)
	buy(&account, "Apple", 10, f)
	// glog.Info("Available account balance after BUY: ", account.getBalance())
	glog.Info(account.Balance == 100)
	glog.Info(account.Available == 90)
	// glog.Info("Account balance after BUY: ", account.Balance)
	commitBuy(&account, f)
	glog.Info(account.Balance == 90)
	glog.Info(account.Available == 90)
	glog.Info(account.StockPortfolio["Apple"] == 10)
	// glog.Info("Available account balance after COMMIT BUY: ", account.getBalance())
	// glog.Info("Account balance after COMMIT BUY: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	// glog.Info("Selling 5 shares of Apple")
	sell(&account, "Apple", 5, f)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 90)
	commitSell(&account, f)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	// glog.Info("Available account balance after COMMIT SELL: ", account.getBalance())
	// glog.Info("Account balance after COMMIT SELL: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	//this should fail
	sell(&account, "Apple", 100, f)
	commitSell(&account, f)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	buy(&account, "Apple", 10000, f)
	commitBuy(&account, f)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	buy(&account, "Apple", 1, f)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 94)
	cancelBuy(&account, f)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	glog.Info(account.StockPortfolio["Apple"] == 5)

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
		glog.Info("USER: ", account.AccountNumber, account.Balance)
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
	case "add":
		add(account, msg.PriceDollars)
		glog.Info("Account balance after adding: ", account.getBalance())
		// UserMap[msg.UserId] = account
	case "buy":
		buy(account, msg.Stock, msg.PriceDollars)
	case "commit_sell":
		commitSell(account)
	case "commit_buy":
		commitBuy(account)
	case "cance_buy":
		cancelSell(account)
	case "cancel_sell":
		cancelSell(account)
	case "set_buy_amount":
		setSellAmount(account, msg.Stock, msg.PriceDollars)
	case "set_sell_amount":
		setSellAmount(account, msg.Stock, msg.PriceDollars)
	case "cancel_set_buy":
		cancelSetBuy(account, msg.Stock)
	case "cancel_set_sell":
		cancelSetSell(account, msg.Stock)
	case "set_buy_trigger":
		setBuyTrigger(account, msg.Stock, msg.PriceDollars)
	case "set_sell_trigger":
		setSellTrigger(account, msg.Stock, msg.PriceDollars)
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

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")

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
>>>>>>> origin/ParseHttpReqs

	log.Fatal(http.ListenAndServe(":9090", router))

}
