package api

import (
	"time"
	"utils"
	"db"

	"github.com/golang/glog"
)

type Account struct {
	AccountNumber  string
	Balance        float64
	Available      float64
	SellStack      utils.Stack
	BuyStack       utils.Stack
	StockPortfolio map[string]float64
	SetBuyMap      map[string]float64
	BuyTriggers    map[string]float64
	SetSellMap     map[string]float64
	SellTriggers   map[string]float64
}

type BuyObj struct {
	Stock       string
	StockAmount float64
	MoneyAmount float64
}

type SetBuy struct {
	Stock       string
	MoneyAmount float64
}

type SellObj struct {
	Stock       string
	StockAmount float64
	MoneyAmount float64
}

type SetSell struct {
	Stock       string
	StockAmount float64
}

func InitializeAccount(value string) Account {
	return Account{
		AccountNumber:  value,
		Balance:        0.0,
		Available:      0.0,
		SellStack:      utils.Stack{},
		BuyStack:       utils.Stack{},
		StockPortfolio: make(map[string]float64),
		SetBuyMap:      make(map[string]float64),
		BuyTriggers:    make(map[string]float64),
		SetSellMap:     make(map[string]float64),
		SellTriggers:   make(map[string]float64),
	}
}

func (account *Account) hasStock(stock string, amount float64) bool {
	//check if the user holds the amount of stock he/she is trying to sell
	return account.StockPortfolio[stock] >= amount
}

// returns the amount that is available to the user (i.e not on hold for any transactions)
func (account *Account) getBalance() float64 {
	return account.Available
}

func (account *Account) holdMoney(amount float64) {
	if amount > 0 {
		account.Available -= amount
		//update db
		db.UpdateAvailableAccountBalance(account.AccountNumber, account.Available)
	} else {
		glog.Error("Cannot hold negative account for the account ", amount)
	}
}

func (account *Account) addMoney(amount float64) {
	account.Balance += amount
	account.Available += amount
	db.UpdateAccountBalance(account.AccountNumber, account.Balance)
	db.UpdateAvailableAccountBalance(account.AccountNumber, account.Available)
	glog.Info("This account now has ", account.Balance, account.Available)
}

func (account *Account) substractBalance(amount float64) {
	account.Balance -= amount
	db.UpdateAccountBalance(account.AccountNumber, account.Balance)
}

func (account *Account) unholdMoney(amount float64) {
	if amount > 0 {
		account.Available += amount
		//update db
		db.UpdateAvailableAccountBalance(account.AccountNumber, account.Available)
	} else {
		glog.Error("Cannot unhold negative account for the account ", amount)
	}
}

/*
	TODO: we probably need to store hold stocks separately
	i.e. the same way we're dealing with the account balance
	to be able to display accurate stock numbers per account at any given time
*/
func (account *Account) holdStock(stock string, amount float64) {
	account.StockPortfolio[stock] -= amount
}

func (account *Account) unholdStock(stock string, amount float64) {
	account.StockPortfolio[stock] += amount
}

// Start a trigger
// should pull quotes every 60 sec to check the price
// then execute BUY/SELL
func (account *Account) startBuyTrigger(stock string, transactionNum int) {
	price := GetQuote(stock, account.AccountNumber)
	limit := account.BuyTriggers[stock]

	//if there is still trigger in the map
	if limit > 0 {
		glog.Info(">>>>>>>>>>>>>>>>>>>BUY TRIGGER CHECK: >>>>>> limit: ", limit, " current: ", price)
		for price > limit {
			glog.Info("Price is still greater than the trigger limit")
			time.Sleep(60 * time.Second)
			price = GetQuote(stock, account.AccountNumber)
		}

		stockNum := account.SetBuyMap[stock]
		Buy(account, stock, stockNum, transactionNum)
		CommitBuy(account, transactionNum)
		//hacky:
		//put money back
		account.Available = account.Balance
		glog.Info("Just BOUGHT stocks for trigger #: ", stockNum)
		glog.Info("Balance: ", account.Balance, " Available: ", account.Available)
		delete(account.SetBuyMap, stock)
	}
}

func (account *Account) startSellTrigger(stock string, transactionNum int) {
	price := GetQuote(stock, account.AccountNumber)
	min := account.SellTriggers[stock]

	//if there is still trigger in the map
	if min > 0 {
		glog.Info(">>>>>>>>>>>>>>>>>>>SELL TRIGGER CHECK: >>>>>> limit: ", min, " current: ", price)
		for price < min {
			glog.Info("Price is still greater than the trigger limit")
			time.Sleep(60 * time.Second)
			price = GetQuote(stock, account.AccountNumber)
		}

		stockNum := account.SetSellMap[stock]
		Sell(account, stock, stockNum, transactionNum)
		CommitSell(account, transactionNum)
		//hacky:
		//put stock back
		account.StockPortfolio[stock] += stockNum
		glog.Info("Just SOLD stocks for trigger #: ", stockNum, stock)
		glog.Info("Accont: ", account)
		glog.Info("Stock balance: ", account.StockPortfolio[stock])
		delete(account.SetSellMap, stock)
	}
}
