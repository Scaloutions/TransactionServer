package main

import (
	"github.com/golang/glog"
	"fmt"
	"html"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, there!")
	testLogic()
}

func testLogic(){
	account := initializeAccount("123")
	add(&account, 100)
	glog.Info("Account balance after adding: ", account.getBalance())
	buy(&account, "Apple", 10)
	glog.Info("Available account balance after BUY: ", account.getBalance())
	glog.Info("Account balance after BUY: ", account.Balance)
	commitBuy(&account)
	glog.Info("Available account balance after COMMIT BUY: ", account.getBalance())
	glog.Info("Account balance after COMMIT BUY: ", account.Balance)
	glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	sell(&account, "Apple", 5)
	commitSell(&account)
	glog.Info("Available account balance after COMMIT SELL: ", account.getBalance())
	glog.Info("Account balance after COMMIT SELL: ", account.Balance)
	glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	//this should fail
	sell(&account, "Apple", 100)
	commitSell(&account)
	glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))}).Methods("GET")

	router.HandleFunc("/getQuote", echoString).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", router))

}
