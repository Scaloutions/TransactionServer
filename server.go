package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func usage() {
	fmt.Println("usage: example -logtostderr=true -stderrthreshold=[INFO|WARN|FATAL|ERROR] -log_dir=[string]\n")
	flag.PrintDefaults()
}

func echoString(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to DayTrading Inc! \n Running some tests...")
	// fmt.Fprintf(w, "Hi, there! Running test function..")
	testLogic()
}

// User map
var UserMap = make(map[string]*Account)

func authenticateUser(c *gin.Context) {
	req := getParams(c)
	account := initializeAccount(req.UserId)
	UserMap[req.UserId] = &account
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

func getParams(c *gin.Context) Request {
	request := Request{}
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		glog.Error("Error processing request: %s", err)
	}

	err = json.Unmarshal(body, &request)
	if err != nil {
		glog.Error("Error parsing JSON: %s", err)
	}

	return request
}

func addReq(c *gin.Context) {
	req := getParams(c)

	var account *Account
	account = getUser(req.UserId)
	glog.Info("\n\n############################### INFO: Executing ADD FOR... ", req.PriceDollars, req.CommandNumber)
	glog.Info(req)
	glog.Info(account)
	add(account, req.PriceDollars)
}

func BuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing BUY FOR... ", req.PriceDollars, req.CommandNumber)
	buy(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: BUY Successful")
}

func SellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func CommitSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing COMMIT SELL ", req.CommandNumber)
	commitSell(account)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: COMMIT SELL Successful")
}
func CommitBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func CancelBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func CancelSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func SetBuyAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func CancelSetBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func CancelSetSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func SetBuyTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func SetSellTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}
func DumplogReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
}

// 	//TODO: rewrite this!!
// 	switch req.Command {
// 	case "authenticate":
// 		glog.Info("\n\n############################### INFO: Executing authenticate... ", req.CommandNumber)
// 		authenticateUser(req.UserId)
// 		glog.Info("\n############################### SUCCESS: Authentication Successful!")
// 	case "add":
// 		glog.Info("\n\n############################### INFO: Executing ADD FOR... ", req.PriceDollars, req.CommandNumber)
// 		add(account, req.PriceDollars)
// 		glog.Info("SUCCESS: Account Balance: ", account.Balance, " Available: ", account.Available)
// 	case "buy":
// 	case "commit_sell":
// 		glog.Info("\n\n############################### INFO: Executing COMMIT SELL ", req.CommandNumber)
// 		commitSell(account)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: COMMIT SELL Successful")
// 	case "commit_buy":
// 		glog.Info("\n\n############################### INFO: Executing COMMIT BUY ", req.CommandNumber)
// 		commitBuy(account)
// 		glog.Info("\n############################### SUCCESS: COMMIT BUY Successful")
// 		// glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		// glog.Info("Account Stocks: ", account.StockPortfolio["S"])
// 	case "cancel_buy":
// 		glog.Info("\n\n############################### INFO: Executing CANCEL BUY ", req.CommandNumber)
// 		cancelSell(account)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: CANCEL BUY Successful")
// 	case "cancel_sell":
// 		glog.Info("\n\n############################### INFO: Executing CANCEL SELL ", req.CommandNumber)
// 		cancelSell(account)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("Account Stocks: ", account.StockPortfolio["S"])
// 		glog.Info("\n############################### SUCCESS: CANCEL SELL Successful")
// 	case "set_buy_amount":
// 		glog.Info("\n\n############################### INFO: Executing SET BUY AMOUNT ", req.CommandNumber)
// 		setSellAmount(account, req.Stock, req.PriceDollars)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: SET BUY AMOUNT Successful")
// 	case "set_sell_amount":
// 		glog.Info("\n\n############################### INFO: Executing SET SELL AMOUNT ", req.CommandNumber)
// 		setSellAmount(account, req.Stock, req.PriceDollars)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: SET SELL AMOUNT Successful")
// 	case "cancel_set_buy":
// 		glog.Info("\n\n############################### INFO: Executing CANCEL SET BUY ", req.CommandNumber)
// 		cancelSetBuy(account, req.Stock)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: CANCEL SET BUY Successful")
// 	case "cancel_set_sell":
// 		glog.Info("\n\n############################### INFO: Executing CANCEL SET SELL ", req.CommandNumber)
// 		cancelSetSell(account, req.Stock)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: CANCEL SET SELL Successful")
// 	case "set_buy_trigger":
// 		glog.Info("\n\n############################### INFO: Executing SET BUY TRIGGER ", req.CommandNumber)
// 		setBuyTrigger(account, req.Stock, req.PriceDollars)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 	case "set_sell_trigger":
// 		glog.Info("\n\n############################### INFO: Executing SET SELL TRIGGER ", req.CommandNumber)
// 		setSellTrigger(account, req.Stock, req.PriceDollars)
// 		glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
// 		glog.Info("\n############################### SUCCESS: SET SELL TRIGGER Successful")
// 	case "dumplog":
// 		glog.Info("SAVING XML LOG FILE")
// 	default:
// 		panic("Oh noooo we can't process this request :(")

// 	}

// 	//Set Content-Type header so that clients will know how to read response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	//Write json response back to response
// 	// w.Write(reqJson)
// 	// return req
// }

func main() {
	router := gin.Default()
	// router := mux.NewRouter()
	flag.Usage = usage
	flag.Parse()

	router.GET("/api/test", echoString)
	// routPOSTunc("/getQuote", echoString
	router.POST("/api/authenticate", authenticateUser)
	router.POST("/api/add", addRequest)
	// router.POST("/api/sell", parseRequest)
	// router.POST("/api/buy", parseRequest)
	// router.POST("/api/commit_sell", parseRequest)
	// router.POST("/api/commit_buy", parseRequest)
	// router.POST("/api/cancel_buy", parseRequest)
	// router.POST("/api/cancel_sell", parseRequest)
	// router.POST("/api/set_buy_amount", parseRequest)
	// router.POST("/api/set_sell_amount", parseRequest)
	// router.POST("/api/cancel_set_buy", parseRequest)
	// router.POST("/api/cancel_set_sell", parseRequest)
	// router.POST("/api/set_buy_trigger", parseRequest)
	// router.POST("/api/set_sell_trigger", parseRequest)
	// router.HandleFunc("/api/", ).Methods("POST")
	//router.HandleFunc("/api/", ).Methods("POST")

	// log.Fatal(http.ListenAndServe(":9090", router))
	log.Fatal(router.Run(":9090"))

}
