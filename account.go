package main

import (
	"github.com/golang/glog"
	"time"
)

type Account struct {
	AccountNumber string
	Balance float64
	Available float64
	SellStack Stack
	BuyStack Stack
	StockPortfolio map[string]float64
	SetBuyMap map[string]float64
	BuyTriggers map[string]float64
	SetSellMap map[string]float64
	SellTriggers map[string]float64
}

type Buy struct {
	Stock string
	StockAmount float64
	MoneyAmount float64
}

type SetBuy struct {
	Stock string
	MoneyAmount float64
}

type Sell struct {
	Stock string
	StockAmount float64
	MoneyAmount float64
}

type SetSell struct {
	Stock string
	StockAmount float64
}

func initializeAccount(value string) Account {
	return Account{
		AccountNumber: value,
		Balance: 0.0,
		Available: 0.0,
		SellStack: Stack{},
		BuyStack: Stack{},
		StockPortfolio: make(map[string]float64),
		SetBuyMap: make(map[string]float64),
		BuyTriggers: make(map[string]float64),
		SetSellMap: make(map[string]float64),
		SellTriggers: make(map[string]float64),
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
	}
}

func (account *Account) unholdMoney(amount float64) {
	if amount > 0 {
		account.Available += amount
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
	account.StockPortfolio[stock] -= amount
}

func (account *Account) addMoney(amount float64) {
	account.Balance += amount
	account.Available += amount
	glog.Info("This account now has ", account.Balance, account.Available)
}

// Start a trigger
// should pull quotes every 60 sec to check the price
// then execute BUY/SELL
// unix.Nono timestamp
func (account *Account) startBuyTrigger(stock string, limit float64) {
	price := getQuote(stock, account.AccountNumber)
	//limit := trigger.MoneyAmount

	for price > limit {
		time.Sleep(60 * time.Millisecond)
		price = getQuote(stock, account.AccountNumber)
		//sleep for 60 sec
	}

	stockNum := account.SetBuyMap[stock]
	buy(account, stock, stockNum)
	//remove st buy
	delete(account.SetBuyMap, stock)	
}