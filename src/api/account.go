package api

import (
	"time"
	"../utils"
	"../db"

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

func GetAccount(userId string) Account {
	dbAccount, err := db.GetAccount(userId)

	if err!=nil {
		glog.Error(err, " ", userId)
	}

	return Account{
		AccountNumber:  dbAccount.UserId,
		Balance:        dbAccount.Balance,
		Available:      dbAccount.Available,
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
	currAmount, err := db.GetUserStockAmount(account.AccountNumber, stock)
	glog.Info("From DB: User ", account.AccountNumber, " has ", stock, " amount as : ", currAmount)

	if err!=nil {
		glog.Error(err, " ", account)
	}
	return currAmount >= amount
}

// returns the amount that is available to the user (i.e not on hold for any transactions)
func (account *Account) getBalance() float64 {
	dbAccount, err := db.GetAccount(account.AccountNumber)

	if err!=nil {
		glog.Error(err, " ", account)
	}
	return dbAccount.Available
}

func (account *Account) holdMoney(amount float64) {
	if amount > 0 {
		db.UpdateAvailableAccountBalance(account.AccountNumber, amount*-1)
	} else {
		glog.Error("Cannot hold negative account for the account ", amount)
	}
}

func (account *Account) addMoney(amount float64) {
	glog.Info("Updating account balance in the DB for user: ", account.AccountNumber)
	err := db.AddMoneyToAccount(account.AccountNumber, amount)
		
	if err!=nil {
		glog.Error(err, " for account:", account)
		return
	}

	glog.Info("This account "+ account.AccountNumber +" now has ", account.Balance, " available: ", account.Available)
}

func (account *Account) updateBalance(amount float64) {
	err := db.UpdateAccountBalance(account.AccountNumber, amount)

	if err!=nil {
		glog.Error(err)
	}
}

func (account *Account) unholdMoney(amount float64) {
	if amount > 0 {
		err := db.UpdateAvailableAccountBalance(account.AccountNumber, amount)
		if err!=nil {
			glog.Error(err)
		}
	} else {
		glog.Error("Cannot unhold negative account for the account ", amount)
	}
}

/*
	TODO: we probably need to store hold stocks separately
	i.e. the same way we're dealing with the account balance
	to be able to display accurate stock numbers per account at any given time
*/
func (account *Account) holdStock(stock string, amount float64) error {
	err := db.UpdateAvailableUserStock(account.AccountNumber, stock, amount*-1)
	if err!=nil {
		glog.Error("Failed to Hold STOCK for ", account)
		return err
	}
	return nil
}

func (account *Account) unholdStock(stock string, amount float64) error {
	account.StockPortfolio[stock] += amount
	err := db.UpdateAvailableUserStock(account.AccountNumber, stock, amount)
	if err!=nil {
		glog.Error("Failed to UnHold STOCK for ", account)
		return err
	}
	return nil
}

// Start a trigger
// should pull quotes every 60 sec to check the price
// then execute BUY/SELL
func (account *Account) startBuyTrigger(stock string, limit float64, transactionNum int) error {
	price, err := GetQuote(stock, account.AccountNumber, transactionNum)
	if err!= nil {
		return err
	}

	setBuy, err := db.GetSetBuy(account.AccountNumber, stock)
	if err!= nil {
		return err
	}
	//if there is still trigger in the map
	if limit > 0 {
		glog.Info(">>>>>>>>>>>>>>>>>>>BUY TRIGGER CHECK: >>>>>> limit: ", limit, " current: ", price)
		for price > limit {
			glog.Info("Price is still greater than the trigger limit")
			time.Sleep(60 * time.Second)
			setBuy, err = db.GetSetBuy(account.AccountNumber, stock)
			glog.Info("Triggers: Associated SET BUY: ", setBuy)
			if err!=nil {
				glog.Info("BUY TRIGGER CANCELLED NO SET BUY SET ANYMORE")
				return nil
			}
			price, err = GetQuote(stock, account.AccountNumber, transactionNum)
			if err!= nil {
				return err
			}
		}

		if err!=nil {
			return err
		}
		glog.Info("Buying Stock as ", setBuy.MoneyAmount, "\\", price)

		stockNum := setBuy.MoneyAmount / price
		account.updateBalance(-1*setBuy.MoneyAmount)
		err := db.AddUserStock(account.AccountNumber, stock, stockNum)
		glog.Info("Just BOUGHT stocks for trigger #: ", stockNum)
		err = db.DeleteSetBuy(account.AccountNumber, stock)

		if err!=nil {
			glog.Info("Error deleting SET BUY for ", account, " stock: ", stock)
			return err
		}
	}

	return nil
}

func (account *Account) startSellTrigger(stock string, min float64, transactionNum int) error {
	price, err := GetQuote(stock, account.AccountNumber, transactionNum)
	if err!= nil {
		return err
	}

	setSell, err := db.GetSetSell(account.AccountNumber, stock)
	if err!= nil {
		return err
	}
	//if there is still trigger in the map
	if min > 0 {
		glog.Info(">>>>>>>>>>>>>>>>>>>SELL TRIGGER CHECK: >>>>>> limit: ", min, " current: ", price)
		for price < min {
			glog.Info("Price is still greater than the trigger limit")
			time.Sleep(60 * time.Second)
			setSell, err = db.GetSetSell(account.AccountNumber, stock)
			if err!=nil {
				glog.Info("SELL TRIGGER CANCELLED NO SET SELL SET ANYMORE")
				return nil
			}
			price, err = GetQuote(stock, account.AccountNumber, transactionNum)
			if err!= nil {
				return err
			}
		}

		revenue := price * setSell.StockAmount
		//decrease stock val 
		db.UpdateUserStock(account.AccountNumber, stock, -1*setSell.StockAmount)
		//update balance
		db.AddMoneyToAccount(account.AccountNumber, revenue)
		glog.Info("Just SOLD stocks for trigger #: ", stock, " amount: ", revenue)
		glog.Info("Accont: ", account)
		err = db.DeleteSetSell(account.AccountNumber, stock)
		if err!=nil {
			return err
		}
	}
	return nil
}
