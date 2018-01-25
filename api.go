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
	stockNum := amount / quote(stock)
	//check balance
	if account.getBalance() < amount {
		//TODO: improve logging
		glog.Info("Not enough money for account ",account, "to buy ", stock)
	} else {
		transaction := Buy{
			Stock: stock,
			MoneyAmount: amount,
			StockAmount: stockNum,
		}
		//add buy transcation to the stack
		account.BuyStack.Push(transaction)
		//hold the money
		account.holdMoney(amount)
	}
}

func sell(account *Account, stock string, amount float64) {
	//check if have that # of stocks
	stockNum := amount / quote(stock)
	if account.hasStock(stock, stockNum){
		transaction := Sell {
			Stock: stock,
			MoneyAmount: amount,
			StockAmount: stockNum, 
		}
		account.SellStack.Push(transaction)

	} else {
		//TODO: improve logging
		glog.Info("Not enough stock ", stock, "to sell.")
	}
}

func commitBuy(account *Account) {
	if account.BuyStack.size >0 {
		//weird go casting
		i := account.BuyStack.Pop()
		transaction := i.(Buy)
		//should we check balance here insted? TODO: clarify
		account.Balance -= transaction.MoneyAmount
		//add number of stocks to user
		//TODO: refactor this line
		account.StockPortfolio[transaction.Stock] += transaction.StockAmount 

	} else {
		glog.Error("No BUY transactions previously set for account: ", account.AccountNumber)
	}
} 

func cancelBuy(account *Account) {
	//TODO: log this
	account.BuyStack.Pop()
} 

func commitSell(account *Account) {
	if account.SellStack.size > 0{
		i := account.SellStack.Pop()
		transaction := i.(Sell)
		account.Balance += transaction.MoneyAmount
		account.StockPortfolio[transaction.Stock] -= transaction.StockAmount 

	} else {
		glog.Error("No SELL transactions previously set for account: ", account.AccountNumber)
	}
} 

func cancelSell(account *Account) {
	//TODO: log this
	account.SellStack.Pop()
} 

/*
Sets a defined amount of the given stock to buy when the current stock price 
is less than or equal to the BUY_TRIGGER
*/
func setBuyAmount(account *Account, stock string, amount float64) {
	//check if there is enough money in the account
	if account.Available >= amount {
		//hold money
		account.Available -= amount
		account.SetBuyMap[stock] += amount
		glog.Info("SetBut for $", amount, " and stock ", stock)
		glog.Info("Total SET BUY on stock ", stock, " is now ", account.SetBuyMap[stock])
	} else {
		glog.Error("Account does not have enough money to buy stock ", stock)
	}
}

func cancelSetBuy(accont Account, stock string) {}

func setBuyTrigger(account *Account, stock string, price float64) {
	//check for set buy on that stock
	if _, ok := account.SetBuyMap[stock]; ok {
		account.BuyTriggers[stock] = price
		glog.Info("Set BUY trigger for ", stock, "at price ", price)
	} else {
		glog.Error("You have to SET BUY AMOUNT on stock ", stock, " first.")
	}
}

func setSellAmount(account Account, stock string, amount float64) {}

func setSellTrigger(account Account, stock string, amount float64) {}

func cancelSetSell(account Account, stock string, amount float64) {}

func dumplog(account Account, filename string) {}

func dumplogAll(filename string) {}


