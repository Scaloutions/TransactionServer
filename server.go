package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
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
	glog.Info("INFO: Retrieving user from the db..")
}

func testLogic() {
	account := initializeAccount("123")
	add(&account, 100)
	glog.Info("Balance: ", account.Balance)
	glog.Info(account.Available)
}

type Request struct {
	UserId        string
	PriceDollars  float64
	PriceCents    float64
	Command       string
	CommandNumber int
	Stock         string
}

func getUser(userId string) *Account {
	return UserMap[userId]
}

func getParams(w http.ResponseWriter, r *http.Request) Request {
	request := Request{}
	//Parse json request body and use it to set fields on user
	//Note that user is passed as a pointer variable so that it's fields can be modified
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		panic(err)
	}

	return request
}

func addRequest(w http.ResponseWriter, r *http.Request) {
	req := getParams(w, r)

	// var account *Account
	account := getUser(req.UserId)
	glog.Info("\n\n############################### INFO: Executing ADD FOR... ", req.PriceDollars, req.CommandNumber)
	add(account, req.PriceDollars)
}

//Parse request and return Response Object
func parseRequest(w http.ResponseWriter, r *http.Request) {
	req := getParams(w, r) //initialize empty user

	//Parse json request body and use it to set fields on user
	//Note that user is passed as a pointer variable so that it's fields can be modified
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil{
	// 	panic(err)
	// }

	var account *Account
	if req.Command != "authenticate" {
		account = getUser(req.UserId)
	}

	//TODO: rewrite this!!
	switch req.Command {
	case "authenticate":
		glog.Info("\n\n############################### INFO: Executing authenticate... ", req.CommandNumber)
		authenticateUser(req.UserId)
		glog.Info("\n############################### SUCCESS: Authentication Successful!")
	case "add":
		glog.Info("\n\n############################### INFO: Executing ADD FOR... ", req.PriceDollars, req.CommandNumber)
		add(account, req.PriceDollars)
		glog.Info("SUCCESS: Account Balance: ", account.Balance, " Available: ", account.Available)
	case "buy":
		glog.Info("\n\n############################### INFO: Executing BUY FOR... ", req.PriceDollars, req.CommandNumber)
		buy(account, req.Stock, req.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: BUY Successful")
	case "commit_sell":
		glog.Info("\n\n############################### INFO: Executing COMMIT SELL ", req.CommandNumber)
		commitSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: COMMIT SELL Successful")
	case "commit_buy":
		glog.Info("\n\n############################### INFO: Executing COMMIT BUY ", req.CommandNumber)
		commitBuy(account)
		glog.Info("\n############################### SUCCESS: COMMIT BUY Successful")
		// glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		// glog.Info("Account Stocks: ", account.StockPortfolio["S"])
	case "cancel_buy":
		glog.Info("\n\n############################### INFO: Executing CANCEL BUY ", req.CommandNumber)
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: CANCEL BUY Successful")
	case "cancel_sell":
		glog.Info("\n\n############################### INFO: Executing CANCEL SELL ", req.CommandNumber)
		cancelSell(account)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
		glog.Info("\n############################### SUCCESS: CANCEL SELL Successful")
	case "set_buy_amount":
		glog.Info("\n\n############################### INFO: Executing SET BUY AMOUNT ", req.CommandNumber)
		setSellAmount(account, req.Stock, req.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: SET BUY AMOUNT Successful")
	case "set_sell_amount":
		glog.Info("\n\n############################### INFO: Executing SET SELL AMOUNT ", req.CommandNumber)
		setSellAmount(account, req.Stock, req.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: SET SELL AMOUNT Successful")
	case "cancel_set_buy":
		glog.Info("\n\n############################### INFO: Executing CANCEL SET BUY ", req.CommandNumber)
		cancelSetBuy(account, req.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: CANCEL SET BUY Successful")
	case "cancel_set_sell":
		glog.Info("\n\n############################### INFO: Executing CANCEL SET SELL ", req.CommandNumber)
		cancelSetSell(account, req.Stock)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: CANCEL SET SELL Successful")
	case "set_buy_trigger":
		glog.Info("\n\n############################### INFO: Executing SET BUY TRIGGER ", req.CommandNumber)
		setBuyTrigger(account, req.Stock, req.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	case "set_sell_trigger":
		glog.Info("\n\n############################### INFO: Executing SET SELL TRIGGER ", req.CommandNumber)
		setSellTrigger(account, req.Stock, req.PriceDollars)
		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
		glog.Info("\n############################### SUCCESS: SET SELL TRIGGER Successful")
	case "dumplog":
		glog.Info("SAVING XML LOG FILE")
	default:
		panic("Oh noooo we can't process this request :(")

	}

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response
	// w.Write(reqJson)
	// return req
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
	router.HandleFunc("/api/commit_sell", parseRequest).Methods("POST")
	router.HandleFunc("/api/commit_buy", parseRequest).Methods("POST")
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
