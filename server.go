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
	glog.Info("\n\n############################### INFO: Executing authenticate... ", req.CommandNumber)
	account := initializeAccount(req.UserId)
	UserMap[req.UserId] = &account
	glog.Info("\n############################### SUCCESS: Authentication Successful!")
	glog.Info("##### Account Balance: ", account.Balance, " Available: ", account.Available)
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

	glog.Info("\n\n############################### INFO: Executing SELL FOR... ", req.PriceDollars, req.CommandNumber)
	sell(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SELL Successful")
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

	glog.Info("\n\n############################### INFO: Executing COMMIT BUY ", req.CommandNumber)
	commitBuy(account)
	glog.Info("\n############################### SUCCESS: COMMIT BUY Successful")
}

func CancelBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL BUY ", req.CommandNumber)
	cancelBuy(account)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL BUY Successful")
}

func CancelSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SELL ", req.CommandNumber)
	cancelSell(account)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SELL Successful")
}

func SetBuyAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET BUY AMOUNT ", req.CommandNumber)
	setBuyAmount(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET BUY AMOUNT Successful")
}

func SetSellAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET SELL AMOUNT ", req.CommandNumber)
	setSellAmount(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET SELL AMOUNT Successful")
}

func CancelSetBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SET BUY ", req.CommandNumber)
	cancelSetBuy(account, req.Stock)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SET BUY Successful")
}

func CancelSetSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SET SELL ", req.CommandNumber)
	cancelSetSell(account, req.Stock)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SET SELL Successful")
}

func SetBuyTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET BUY TRIGGER ", req.CommandNumber)
	setBuyTrigger(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
}

func SetSellTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET SELL TRIGGER ", req.CommandNumber)
	setSellTrigger(account, req.Stock, req.PriceDollars)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET SELL TRIGGER Successful")
}

func DumplogReq(c *gin.Context) {
	req := getParams(c)

	glog.Info("SAVING XML LOG FILE")
	c.String(http.StatusOK, "Getting logs...")
}

func main() {
	router := gin.Default()
	// router := mux.NewRouter()
	flag.Usage = usage
	flag.Parse()

	router.GET("/api/test", echoString)
	// routPOSTunc("/getQuote", echoString
	router.POST("/api/authenticate", authenticateUser)
	router.POST("/api/add", addReq)
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
