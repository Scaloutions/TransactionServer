package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"db"
	"api"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func usage() {
	fmt.Println("usage: example -logtostderr=true -stderrthreshold=[INFO|WARN|FATAL|ERROR] -log_dir=[string]\n")
	flag.PrintDefaults()
}

func echoString(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to DayTrading Inc!")
}

var UserMap = make(map[string]*api.Account)

type Request struct {
	UserId        string
	PriceDollars  float64
	PriceCents    float64
	Command       string
	CommandNumber int
	Stock         string
}

func authenticateUser(userId string) {
	_, err := db.GetUser(userId)
	/*
		FOR THE PROJECT XML requirements
		automatically create users for testing
	*/
	if err!=nil {
		db.CreateNewUser(userId, "", "", "")
		db.CreateNewAccount(userId)
	}

	account := api.GetAccount(userId)
	UserMap[userId] = &account
	glog.Info("\nSUCCESS: Authentication Successful!")
	glog.Info("\nAccount Balance: ", account.Balance, " Available: ", account.Available, "User: ", userId)
}

// Gets user from the memory: assumes we authenticate user first
func getUser(userId string) *api.Account {
	if user, ok := UserMap[userId]; ok {
		//do something here
		return user
	} else {
		authenticateUser(userId)
		return UserMap[userId]
	}
}

/*
	TODO: add API point for creating user record with personal info
*/
func createUser(userId string, name string, email string, address string) {
	db.CreateNewUser(userId, name, email, address)
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

func authReq(c *gin.Context) {
	req := getParams(c)
	glog.Info("\n Executing AUTHENTICATE for user: ", req.UserId)
	authenticateUser(req.UserId)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func getQuoteReq(c *gin.Context) {
	req := getParams(c)

	glog.Info("\n Executing QUOTE: ", req)
	api.GetQuote(req.Stock, req.UserId, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func addReq(c *gin.Context) {
	req := getParams(c)

	var account *api.Account
	account = getUser(req.UserId)
	glog.Info("\n Executing ADD: ", req)
	api.Add(account, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func buyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing BUY ", req)
	api.Buy(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func sellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SELL ", req)
	api.Sell(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func commitSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing COMMIT SELL ", req)
	
	api.CommitSell(account, req.CommandNumber)
	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})	
}

func commitBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing COMMIT BUY ", req)
	api.CommitBuy(account, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func cancelBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL BUY ", req)
	api.CancelBuy(account, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func cancelSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SELL ", req)
	api.CancelSell(account, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func setBuyAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET BUY AMOUNT ", req)
	api.SetBuyAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func setSellAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET SELL AMOUNT ", req)
	api.SetSellAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func cancelSetBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SET BUY ", req)
	api.CancelSetBuy(account, req.Stock, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func cancelSetSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SET SELL ", req)
	api.CancelSetSell(account, req.Stock, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func setBuyTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET BUY TRIGGER ", req)
	api.SetBuyTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func setSellTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET SELL TRIGGER ", req)
	api.SetSellTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)

	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
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

	// db connection
	db.InitializeDB()

	api := router.Group("/api")
	{
		api.GET("/test", echoString)
		api.GET("/dumplog", dumplogReq)
		api.GET("/get_quote", getQuoteReq)
		api.POST("/authenticate", authReq)
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
