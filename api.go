package main

import "github.com/golang/glog"
//below are all the functions that need to be implemented in the system

func add(account *Account, amount float64) {
	if amount > 0{
		account.addMoney(amount)
		glog.Info("Added ", amount)
	} else {
        glog.Error("Cannot add negative amount to balance ", amount)
	}
}

func quote(stock string) float64 { 
	return 1
}

func buy(account *Account, stock string, amount float64) {
	//get qoute
	total := quote(stock) * amount
	//check balance
	if account.getBalance() < total {
		//TODO: improve logging
		glog.Info("Not enough money for account ",account, "to buy ", stock)
	} else {
		transaction := Buy{
			Stock: stock,
			Amount: amount,
		}
		//add buy transcation to the stack
		account.BuyStack.Push(transaction)
		//hold the money
		account.holdMoney(total)
	}
}

func sell(account Account, stock string, amount float64) {
	//check if have that # of stocks
	if account.hasStock(stock, amount){
		//get quote
		total := quote(stock) * amount
		transaction := Sell {
			Stock: stock,
			Amount: total,
		}
		account.SellStack.Push(transaction)

	} else {
		//TODO: improve logging
		glog.Info("Not enough stock ", stock, "to sell.")
	}
}
func commitBuy(account Account) {} 

func cancelBuy(account Account) {
	//TODO: log this
	account.BuyStack.Pop()
} 

func commitSell(account Account) {} 

func cancelSell(account Account) {
	//TODO: log this
	account.SellStack.Pop()
} 

func setBuyAmount(account Account, stock string, amount float64) {}

func cancelSetBuy(accont Account, storck string) {}

func setBuyTrigger(account Account, stock string, amount float64) {}

func setSellAmount(account Account, stock string, amount float64) {}

func setSellTrigger(account Account, stock string, amount float64) {}

func cancelSetSell(account Account, stock string, amount float64) {}

func dumplog(account Account, filename string) {}

func dumplogAll(filename string) {}


