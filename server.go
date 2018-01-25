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
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))}).Methods("GET")

	router.HandleFunc("/getQuote", echoString).Methods("GET")

	log.Fatal(http.ListenAndServe(":9090", router))

}
