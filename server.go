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

	// TODO: add error checking here
	c.JSON(200, gin.H{
		"transaction_num": req.CommandNumber,
		"user_id": req.UserId,
	})
}

func successfulResponse(c *gin.Context, tranNum int, userId string) {
	c.JSON(200, gin.H{
		"transaction_num": tranNum,
		"user_id": userId,
	})
}

func errorResponse(c *gin.Context, tranNum int, userId string) {
	c.JSON(500, gin.H{
		"transaction_num": tranNum,
		"user_id": userId,
	})
}

func getQuoteReq(c *gin.Context) {
	req := getParams(c)
	// TODO:
	// This doesn't work because it is a GET request so we need to obtain params differently!!
	// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	glog.Info("Request params for get quote: ", req)

	glog.Info("\n Executing QUOTE: ", req)
	quote, err := api.GetQuote(req.Stock, req.UserId, req.CommandNumber)

	if err!=nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		c.JSON(200, gin.H{
			"transaction_num": req.CommandNumber,
			"user_id": req.UserId,
			"quote": quote,
		})
	}
}

func addReq(c *gin.Context) {
	req := getParams(c)

	var account *api.Account
	account = getUser(req.UserId)
	glog.Info("\n Executing ADD: ", req)
	err := api.Add(account, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func buyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing BUY ", req)
	err := api.Buy(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func sellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SELL ", req)
	err := api.Sell(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func commitSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)
	glog.Info("\n Executing COMMIT SELL ", req)
	err := api.CommitSell(account, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func commitBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing COMMIT BUY ", req)
	err := api.CommitBuy(account, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func cancelBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL BUY ", req)
	err := api.CancelBuy(account, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func cancelSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SELL ", req)
	err := api.CancelSell(account, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func setBuyAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET BUY AMOUNT ", req)
	err := api.SetBuyAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func setSellAmountReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET SELL AMOUNT ", req)
	err := api.SetSellAmount(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func cancelSetBuyReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SET BUY ", req)
	err := api.CancelSetBuy(account, req.Stock, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func cancelSetSellReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing CANCEL SET SELL ", req)
	err := api.CancelSetSell(account, req.Stock, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func setBuyTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET BUY TRIGGER ", req)
	err := api.SetBuyTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
}

func setSellTriggerReq(c *gin.Context) {
	req := getParams(c)
	account := getUser(req.UserId)

	glog.Info("\n Executing SET SELL TRIGGER ", req)
	err := api.SetSellTrigger(account, req.Stock, req.PriceDollars, req.CommandNumber)

	if err != nil {
		errorResponse(c,req.CommandNumber, req.UserId)
	} else {
		successfulResponse(c,req.CommandNumber, req.UserId)
	}
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
