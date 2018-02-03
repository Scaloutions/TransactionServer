package main

import (
	"github.com/golang/glog"
)

const (
	ADD             = "add"
	BUY             = "buy"
	SELL            = "sell"
	COMMIT_BUY      = "commit_buy"
	COMMIT_SELL     = "commit_sell"
	CANCEL_BUY      = "cancel_buy"
	CANCEL_SELL     = "cancel_sell"
	SET_BUY_AMOUNT  = "set_buy_amount"
	SET_SELL_AMOUNT = "set_sell_amount"
)

func add(account *Account, amount float64, transactionNum int) {
	if amount > 0 {
		account.addMoney(amount)
		//TODO: log userid instead of account number
		//this logs transaction event
		log := getTransactionEvent(transactionNum, ADD, account.AccountNumber, amount)
		logEvent(log)
		glog.Info("SUCCESS: Added ", amount)
	} else {
		glog.Error("ERROR: Cannot add negative amount to balance ", amount)
	}
}

func getQuote(stock string, userid string) float64 {
	quoteObj := getQuoteFromQS(userid, stock)
	//TODO: log quote server hit here
	return quoteObj.Price
}

func buy(account *Account, stock string, amount float64, transactionNum int) {
	//get qoute
	stockNum := amount / getQuote(stock, account.AccountNumber)
	//check balance
	if account.getBalance() < amount {
		//TODO: improve logging
		err := "Account does not have enough money to execute BUY command"
		glog.Info("Not enough money on account ", account.AccountNumber, " to buy ", stock)
		log := getErrorEvent(transactionNum, BUY, account.AccountNumber, stock, amount, err)
		logEvent(log)
	} else {
		transaction := Buy{
			Stock:       stock,
			MoneyAmount: amount,
			StockAmount: stockNum,
		}
		//add buy transcation to the stack
		account.BuyStack.Push(transaction)
		//hold the money
		account.holdMoney(amount)

		log := getSystemEvent(transactionNum, BUY, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("SUCCESS: Executed BUY for ", amount)
	}
}

func sell(account *Account, stock string, amount float64, transactionNum int) {
	//check if have that # of stocks
	stockNum := amount / getQuote(stock, account.AccountNumber)
	if account.hasStock(stock, stockNum) {
		transaction := Sell{
			Stock:       stock,
			MoneyAmount: amount,
			StockAmount: stockNum,
		}
		//this is fine becasue commit transaction has to be executed within 60sec
		//which means that the qoute does not change
		account.SellStack.Push(transaction)
		account.holdStock(stock, stockNum)

		log := getSystemEvent(transactionNum, BUY, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("Executed SELL for ", amount)
	} else {
		//TODO: improve logging
		err := "User doesn not have enough stock to sell."
		glog.Info("WARNING: Not enough stock ", stock, " to sell.")
		log := getErrorEvent(transactionNum, SELL, account.AccountNumber, stock, amount, err)
		logEvent(log)
	}
}

func commitBuy(account *Account, transactionNum int) {
	if account.BuyStack.size > 0 {
		//weird go casting
		i := account.BuyStack.Pop()
		transaction := i.(Buy)
		//should we check balance here insted? TODO: clarify
		account.Balance -= transaction.MoneyAmount
		//add number of stocks to user
		account.StockPortfolio[transaction.Stock] += transaction.StockAmount

		log := getTransactionEvent(transactionNum, COMMIT_BUY, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("SUCCESS: Executed COMMIT BUY")
		logEvent(log)

	} else {
		err := "No BUY transactions previously set for account"
		//TODO: figure out if we can simplify this logging with some missing parameters
		log := getErrorEvent(transactionNum, COMMIT_BUY, account.AccountNumber, "", 0, err)
		glog.Error("ERROR: No BUY transactions previously set for account: ", account.AccountNumber)
		logEvent(log)
	}
}

func cancelBuy(account *Account, transactionNum int) {
	if account.BuyStack.size > 0 {
		i := account.BuyStack.Pop()
		transaction := i.(Buy)
		//add money back to Available Balance
		account.unholdMoney(transaction.MoneyAmount)
		glog.Info("Executed CANCEL BUY")

		log := getSystemEvent(transactionNum, CANCEL_BUY, account.AccountNumber, transaction.Stock, transaction.MoneyAmount)
		logEvent(log)
	} else {
		err := "There are no BUY transcations to cancel for this account"
		log := getErrorEvent(transactionNum, CANCEL_BUY, account.AccountNumber, "", 0, err)
		glog.Error(err, " ", account.AccountNumber)
		logEvent(log)
	}
}

func commitSell(account *Account, transactionNum int) {
	if account.SellStack.size > 0 {
		i := account.SellStack.Pop()
		transaction := i.(Sell)
		account.addMoney(transaction.MoneyAmount)

		log := getTransactionEvent(transactionNum, COMMIT_SELL, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("Executed COMMIT SELL")
		logEvent(log)
	} else {
		err := "No SELL transactions previously set for account"
		glog.Error(err, " ", account.AccountNumber)
		log := getErrorEvent(transactionNum, COMMIT_SELL, account.AccountNumber, "", 0, err)
		logEvent(log)
	}
}

func cancelSell(account *Account, transactionNum int) {
	if account.SellStack.size > 0 {
		i := account.SellStack.Pop()
		transaction := i.(Sell)
		account.unholdStock(transaction.Stock, transaction.StockAmount)
		glog.Info("Executed CANCEL SELL")

		log := getSystemEvent(transactionNum, CANCEL_SELL, account.AccountNumber, transaction.Stock, transaction.MoneyAmount)
		logEvent(log)
	} else {
		err := "There are no SELL transcations to cancel for this account"
		glog.Error(err, " ", account.AccountNumber)
		log := getErrorEvent(transactionNum, CANCEL_SELL, account.AccountNumber, "", 0, err)
		logEvent(log)
	}
}

/*
Sets a defined amount of the given stock to buy when the current stock price
is less than or equal to the BUY_TRIGGER
*/
func setBuyAmount(account *Account, stock string, amount float64, transactionNum int) {
	//check if there is enough money in the account
	if account.Available >= amount {
		//hold money
		account.holdMoney(amount)
		account.SetBuyMap[stock] += amount

		log := getSystemEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("Executed SET BUY for $", amount, " and stock ", stock)
		glog.Info("Total SET BUY on stock ", stock, " is now ", account.SetBuyMap[stock])
	} else {
		err := "Account does not have enough money to buy stock"
		log := getErrorEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", stock)
	}
}

/* Cancels SET BUY associated with a particular stock
   TODO: verify what happens if the user set multiple SET BUY on one stock
   It shouldbe overwritten by the most recent one!
   TODO: fix this.
*/
func cancelSetBuy(account *Account, stock string) {
	//put money back
	account.unholdMoney(account.SetBuyMap[stock])
	//cancel SET BUYs
	delete(account.SetBuyMap, stock)
	//cancel the trigger
	delete(account.BuyTriggers, stock)
	glog.Info("Executed CANCEL SET BUY")
}

func setBuyTrigger(account *Account, stock string, price float64, transactionNum int) {
	//check for set buy on that stock
	if _, ok := account.SetBuyMap[stock]; ok {
		//TODO: this is hacky we need to properly check for the key here
		if _, exists := account.BuyTriggers[stock]; exists {
			glog.Info("Trigger is already running!")
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up go routine")
			//prevent race condition here TODO: rewrite
			account.BuyTriggers[stock] = price
			go account.startBuyTrigger(stock, transactionNum)
		}

		account.BuyTriggers[stock] = price
		glog.Info("Set BUY trigger for ", stock, " at price ", price)
	} else {
		glog.Error("You have to SET BUY AMOUNT on stock ", stock, " first.")
	}
}

func setSellAmount(account *Account, stock string, amount float64, transactionNum int) {
	if account.StockPortfolio[stock] > amount {
		account.SetSellMap[stock] += amount
		//hold stock
		account.holdStock(stock, amount)

		log := getSystemEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("Executed SET SELL AMOUNT for ", amount)
	} else {
		err := "Account does not have enough stock to sell "
		log := getErrorEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", stock)
	}
}

func setSellTrigger(account *Account, stock string, price float64, transactionNum int) {
	//check for set buy on that stock
	if _, ok := account.SetSellMap[stock]; ok {
		//TODO: this is hacky we need to properly check for the key here
		if _, exists := account.SellTriggers[stock]; exists {
			glog.Info("Sell Trigger is already running!")
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up go routine SEll trigger")
			account.SellTriggers[stock] = price
			go account.startSellTrigger(stock, transactionNum)
		}

		account.SellTriggers[stock] = price
		glog.Info("Set SELL trigger for ", stock, " at price ", price)
	} else {
		glog.Error("You have to SET SELL AMOUNT on stock ", stock, " first.")
	}
}

func cancelSetSell(account *Account, stock string) {
	//put stock back
	account.unholdStock(stock, account.SetSellMap[stock])
	//cancel SET SELLs
	delete(account.SetSellMap, stock)
	//cancel the trigger
	delete(account.SellTriggers, stock)
	glog.Info("Executed CANCEL SET SELL")
}

func dumplog(account *Account, filename string) {}

func dumplogAll(filename string) {}
