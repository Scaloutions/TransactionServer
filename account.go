package main

import (
	"github.com/golang/glog"
)

type Account struct {
	AccountNumber string
	Balance float64
	Available float64
	SellStack Stack
	BuyStack Stack
}

type Buy struct {
	Stock string
	Amount float64
}

type Sell struct {
	Stock string
	Amount float64
}

func initializeAccount(value string) Account {
	return Account{
		AccountNumber: value,
		Balance: 0.0,
		Available: 0.0,
		SellStack: Stack{},
		BuyStack: Stack{},
	}
}

func (account *Account) hasStock(stock string, amount float64) bool {
	//check if the user holds the amount of stock he/she is trying to sell
	return true
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

func (account *Account) addMoney(amount float64) {
	account.Balance += amount
	account.Available += amount
	glog.Info("This account now has ", account.Balance, account.Available)
}