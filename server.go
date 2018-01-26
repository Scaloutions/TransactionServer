package main

import (
	"github.com/golang/glog"
	"fmt"
	//"html"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
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

func parseRequest(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	userId := r.FormValue("userId")

	price, err := strconv.ParseFloat(r.FormValue("priceDollars"), 64)
	if err != nil {
		glog.Error("Cannot parse POST REQ")
	}

	// price := float64(r.FormValue("priceDollars"))
	fmt.Fprintf(w, "userId= %s\n", userId)
	fmt.Fprintf(w, "Price = %s\n", price)

	account := initializeAccount(userId)
	//call add here
	add(&account, price)
	glog.Info("Account balance after adding: ", account.getBalance())
}

func main() {
	router := mux.NewRouter()

	// router.HandleFunc("/", parseRequest)
	// router.HandleFunc("/getQuote", echoString).Methods("GET")
	router.HandleFunc("/api/buy", parseRequest).Methods("POST")
	// router.HandleFunc("/api/sell", ).Methods("GET")
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
