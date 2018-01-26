package main

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, there!")
	testLogic()
}

func testLogic() {

	f := getFilePointer()

	account := initializeAccount("123")
	add(&account, 100, f)
	// glog.Info("Account balance after adding: ", account.getBalance())
	glog.Info(account.Balance == 100)
	buy(&account, "Apple", 10)
	// glog.Info("Available account balance after BUY: ", account.getBalance())
	glog.Info(account.Balance == 100)
	glog.Info(account.Available == 90)
	// glog.Info("Account balance after BUY: ", account.Balance)
	commitBuy(&account)
	glog.Info(account.Balance == 90)
	glog.Info(account.Available == 90)
	glog.Info(account.StockPortfolio["Apple"] == 10)
	// glog.Info("Available account balance after COMMIT BUY: ", account.getBalance())
	// glog.Info("Account balance after COMMIT BUY: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	// glog.Info("Selling 5 shares of Apple")
	sell(&account, "Apple", 5)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 90)
	commitSell(&account)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	// glog.Info("Available account balance after COMMIT SELL: ", account.getBalance())
	// glog.Info("Account balance after COMMIT SELL: ", account.Balance)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	//this should fail
	sell(&account, "Apple", 100)
	commitSell(&account)
	// glog.Info("User has ", account.StockPortfolio["Apple"], " Apple stocks.")
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	buy(&account, "Apple", 10000)
	commitBuy(&account)
	glog.Info(account.StockPortfolio["Apple"] == 5)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	buy(&account, "Apple", 1)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 94)
	cancelBuy(&account)
	glog.Info(account.Balance == 95)
	glog.Info(account.Available == 95)
	glog.Info(account.StockPortfolio["Apple"] == 5)

}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	}).Methods("GET")

	router.HandleFunc("/getQuote", echoString).Methods("GET")
	//router.HandleFunc("/api/buy", ).Methods()
	//router.HandleFunc("/api/sell", ).Methods()
	//router.HandleFunc("/api/commit_sell", ).Methods()
	//router.HandleFunc("/api/commit_buy", ).Methods()
	//router.HandleFunc("/api/cancel_buy", ).Methods()
	//router.HandleFunc("/api/cancel_sell", ).Methods()
	//router.HandleFunc("/api/set_buy_amount", ).Methods()
	//router.HandleFunc("/api/set_sell_amount", ).Methods()
	//router.HandleFunc("/api/", ).Methods()

	log.Fatal(http.ListenAndServe(":9090", router))

}
