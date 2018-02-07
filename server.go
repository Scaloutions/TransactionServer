package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"api"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func usage() {
	fmt.Println("usage: example -logtostderr=true -stderrthreshold=[INFO|WARN|FATAL|ERROR] -log_dir=[string]\n")
	flag.PrintDefaults()
}

func echoString(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to DayTrading Inc! \n Running some tests...")
	testLogic()
}

// User map
var UserMap = make(map[string]*api.Account)

func authenticateUser(c *gin.Context) {
	req := getParams(c)
	glog.Info("\n\n############################### INFO: Executing authenticate... ", req.CommandNumber)
	account := api.InitializeAccount(req.UserId)
	UserMap[req.UserId] = &account
	glog.Info("\n############################### SUCCESS: Authentication Successful!")
	glog.Info("##### Account Balance: ", account.Balance, " Available: ", account.Available)
}

func testLogic() {
	account := api.InitializeAccount("123")
	api.Add(&account, 100, 1)
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

func getUser(userId string) *api.Account {
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

func getQuoteReq(c *gin.Context) {
	req := getParams(c)

	glog.Info("\n\n############################### INFO: Executing QUOTE FOR... ", req.Stock)
	api.GetQuote(req.Stock, req.UserId)
	glog.Info("\n############################### SUCCESS: QUOTE Successful")
}

func addReq(c *gin.Context) {
	req := getParams(c)

	var account *api.Account
	account = getUser(req.UserId)
	glog.Info("\n\n############################### INFO: Executing ADD FOR... ", req.PriceDollars, req.CommandNumber)
	glog.Info(req)
	glog.Info(account)
	api.Add(account, req.PriceDollars, req.CommandNumber)
}

func buyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing BUY FOR... ", req.PriceDollars, req.CommandNumber)
	api.Buy(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: BUY Successful")
}

func sellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SELL FOR... ", req.PriceDollars, req.CommandNumber)
	api.Sell(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SELL Successful")
}

func commitSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing COMMIT SELL ", req.CommandNumber)
	api.CommitSell(account, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: COMMIT SELL Successful")
}
func commitBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing COMMIT BUY ", req.CommandNumber)
	api.CommitBuy(account, req.CommandNumber)
	glog.Info("\n############################### SUCCESS: COMMIT BUY Successful")
}

func cancelBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL BUY ", req.CommandNumber)
	api.CancelBuy(account, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL BUY Successful")
}

func cancelSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SELL ", req.CommandNumber)
	api.CancelSell(account, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SELL Successful")
}

func setBuyAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET BUY AMOUNT ", req.CommandNumber)
	api.SetBuyAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET BUY AMOUNT Successful")
}

func setSellAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET SELL AMOUNT ", req.CommandNumber)
	api.SetSellAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET SELL AMOUNT Successful")
}

func cancelSetBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SET BUY ", req.CommandNumber)
	api.CancelSetBuy(account, req.Stock, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SET BUY Successful")
}

func cancelSetSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing CANCEL SET SELL ", req.CommandNumber)
	api.CancelSetSell(account, req.Stock, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: CANCEL SET SELL Successful")
}

func setBuyTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET BUY TRIGGER ", req.CommandNumber)
	api.SetBuyTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
}

func setSellTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n\n############################### INFO: Executing SET SELL TRIGGER ", req.CommandNumber)
	api.SetSellTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)
	glog.Info("Account Balance: ", account.Balance, " Available: ", account.Available)
	glog.Info("\n############################### SUCCESS: SET SELL TRIGGER Successful")
}

func dumplogReq(c *gin.Context) {
	// req := getParams(c)

	glog.Info("SAVING XML LOG FILE")
	c.String(http.StatusOK, "Getting logs...")
}

func main() {
	router := gin.Default()
	//glog initialization flags
	flag.Usage = usage
	flag.Parse()

	api := router.Group("/api")
	{
		api.GET("/test", echoString)
		api.GET("/dumplog", dumplogReq)
		api.GET("/get_quote", getQuoteReq)
		api.POST("/authenticate", authenticateUser)
		api.POST("/add", addReq)
		api.POST("/buy", buyReq)
		api.POST("/sell", sellReq)
		api.POST("/commit_sell", commitSellReq)
		api.POST("/commit_buy", commitBuyReq)
		api.POST("/cancel_buy", cancelBuyReq)
		api.POST("/cancel_sell", cancelSellReq)
		api.POST("/set_buy_amount", setBuyAmountReq)
		api.POST("/set_sell_amount", setSellAmountReq)
		api.POST("/cancel_set_buy", cancelSetBuyReq)
		api.POST("/cancel_set_sell", cancelSetSellReq)
		api.POST("/set_buy_trigger", setBuyTriggerReq)
		api.POST("/set_sell_trigger", setSellTriggerReq)
	}

	log.Fatal(router.Run(":9090"))

}
