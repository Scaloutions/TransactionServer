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
	fmt.Fprintf(w, "Hi, there!")
	testLogic()
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
	buy(&account, "Apple", 1)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==94)
	cancelBuy(&account)
	glog.Info(account.Balance==95)
	glog.Info(account.Available==95)
	glog.Info(account.StockPortfolio["Apple"]==5)

}

type Response struct {
	UserId string
	PriceDollars float64
	PriceCents float64
	Command string
	CommandNumber int
	Stock string
}

//Parse request and return Response Object
func parseRequest(w http.ResponseWriter, r *http.Request){
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

	glog.Info("USERID: ", msg.UserId)
	glog.Info("PRICE: ", msg.PriceDollars)
	glog.Info("Command: ", msg.Command)
	glog.Info("CommandNo: ", msg.CommandNumber)
	glog.Info("Stock: ", msg.Stock)

	account := initializeAccount(msg.UserId)
	add(&account, msg.PriceDollars)
	glog.Info("Account balance after adding: ", account.getBalance())

	//Set Content-Type header so that clients will know how to read response
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	//Write json response back to response 
	// w.Write(msgJson)
}

func main() {
	router := mux.NewRouter()

	// router.HandleFunc("/", parseRequest)
	// router.HandleFunc("/getQuote", echoString).Methods("GET")
	router.HandleFunc("/api/add", parseRequest).Methods("POST")
	// router.HandleFunc("/api/sell", ).Methods("GET")
	// router.HandleFunc("/api/buy", ).Methods("GET")
	// router.HandleFunc("/api/commit_sell", ).Methods("GET")
	// router.HandleFunc("/api/commit_buy", ).Methods("GET")
	// router.HandleFunc("/api/cancel_buy", ).Methods("GET")
	// router.HandleFunc("/api/cancel_sell", ).Methods("GET")
	// router.HandleFunc("/api/set_buy_amount", ).Methods("GET")
	// router.HandleFunc("/api/set_sell_amount", ).Methods("GET")
	// router.HandleFunc("/api/cancel_set_buy", ).Methods("GET")
	// router.HandleFunc("/api/cancel_set_sell", ).Methods("GET")
	// router.HandleFunc("/api/set_buy_trigger", ).Methods("GET")
	// router.HandleFunc("/api/set_sell_trigger", ).Methods("GET")
	// router.HandleFunc("/api/", ).Methods()
	//router.HandleFunc("/api/", ).Methods()

	log.Fatal(http.ListenAndServe(":9090", router))

}
