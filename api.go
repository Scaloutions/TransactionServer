package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

//below are all the functions that need to be implemented in the system

func add(account *Account, amount float64, f *os.File) {

	server := "CLT1"
	transactionNum := 1
	command := "ADD"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		"",
		amount)
	logging(userCommand, f)

	if amount > 0 {

		account.addMoney(amount)
		glog.Info("Added ", amount)

	} else {

		errMsg := fmt.Sprintf("Cannot add negative amount %.2f to balance ", amount)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			"",
			amount,
			errMsg)
		logging(errorEvent, f)

	}
}

func getQuote(stock string, userid string, file *os.File) float64 {

	server := "CLT1"
	transactionNum := 2
	command := "QUOTE"

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		userid,
		stock,
		0)
	logging(userCommand, file)

	// quoteObj := getQuoteFromQS(userid, stock)
	// //TODO: log quote server hit here

	// return quoteObj.Price
	return 1
}

func buy(
	account *Account, stock string, amount float64, file *os.File) {

	//get qoute
	stockNum := amount / getQuote(stock, account.AccountNumber, file)

	server := "CLT1"
	transactionNum := 3
	command := "BUY"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		amount)
	logging(userCommand, file)

	//check balance
	if account.getBalance() < amount {

		errMsg := fmt.Sprintf(
			"Not enough money for account %s to buy %s",
			account.AccountNumber, stock)
		glog.Info(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			amount,
			errMsg)
		logging(errorEvent, file)

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
	}
}

func sell(account *Account, stock string, amount float64, file *os.File) {

	server := "CLT1"
	transactionNum := 4
	command := "SELL"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		amount)
	logging(userCommand, file)

	//check if have that # of stocks
	stockNum := amount / getQuote(stock, account.AccountNumber, file)
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

	} else {

		errMsg := fmt.Sprintf("Not enough stock %s to sell.", stock)
		glog.Info(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			amount,
			errMsg)
		logging(errorEvent, file)

	}
}

func commitBuy(account *Account, file *os.File) {

	server := "CLT1"
	transactionNum := 5
	command := "COMMIT_BUY"
	username := account.AccountNumber

	if account.BuyStack.size > 0 {
		//weird go casting
		i := account.BuyStack.Pop()
		transaction := i.(Buy)

		userCommand := getUserCommand(
			server,
			transactionNum,
			command,
			username,
			transaction.Stock,
			transaction.MoneyAmount)
		logging(userCommand, file)

		//should we check balance here insted? TODO: clarify
		account.Balance -= transaction.MoneyAmount
		//add number of stocks to user
		//TODO: refactor this line
		account.StockPortfolio[transaction.Stock] += transaction.StockAmount

	} else {

		errMsg := fmt.Sprintf("No BUY transactions previously set for account %s", username)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			"",
			0,
			errMsg)
		logging(errorEvent, file)
	}
}

func cancelBuy(account *Account, file *os.File) {

	server := "CLT1"
	transactionNum := 6
	command := "CANCEL_BUY"
	username := account.AccountNumber

	i := account.BuyStack.Pop()
	transaction := i.(Buy)
	//add money back to Available Balance
	funds := transaction.MoneyAmount

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		transaction.Stock,
		transaction.MoneyAmount)
	logging(userCommand, file)

	account.unholdMoney(funds)

}

func commitSell(account *Account, file *os.File) {

	server := "CLT1"
	transactionNum := 7
	command := "COMMIT_SELL"
	username := account.AccountNumber

	if account.SellStack.size > 0 {
		i := account.SellStack.Pop()
		transaction := i.(Sell)
		funds := transaction.MoneyAmount
		account.addMoney(funds)
		//we already holded those stocks before
		//account.StockPortfolio[transaction.Stock] -= transaction.StockAmount
		userCommand := getUserCommand(
			server,
			transactionNum,
			command,
			username,
			transaction.Stock,
			funds)
		logging(userCommand, file)
	} else {
		errMsg := fmt.Sprintf("No SELL transactions previously set for account %s", username)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			"",
			0,
			errMsg)
		logging(errorEvent, file)
	}
}

func cancelSell(account *Account, file *os.File) {

	server := "CLT1"
	transactionNum := 8
	command := "CANCEL_SELL"
	username := account.AccountNumber

	//TODO: log this
	i := account.SellStack.Pop()
	transaction := i.(Sell)
	stock := transaction.Stock

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		transaction.MoneyAmount)
	logging(userCommand, file)

	account.unholdStock(stock, transaction.StockAmount)

}

/*
Sets a defined amount of the given stock to buy when the current stock price
is less than or equal to the BUY_TRIGGER
*/
func setBuyAmount(
	account *Account, stock string, amount float64, file *os.File) {

	server := "CLT1"
	transactionNum := 9
	command := "SET_BUY_AMOUNT"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		amount)
	logging(userCommand, file)

	//check if there is enough money in the account
	if account.Available >= amount {
		//hold money
		account.holdMoney(amount)
		account.SetBuyMap[stock] += amount
		glog.Info("SetBut for $", amount, " and stock ", stock)
		glog.Info("Total SET BUY on stock ", stock, " is now ", account.SetBuyMap[stock])

	} else {
		errMsg := fmt.Sprintf(
			"Account %s does not have enough money to buy stock %s",
			username, stock)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			amount,
			errMsg)
		logging(errorEvent, file)
	}
}

/* Cancels SET BUY associated with a particular stock
   TODO: verify what happens if the user set multiple SET BUY on one stock
   ?????
*/
func cancelSetBuy(account *Account, stock string, file *os.File) {

	server := "CLT1"
	transactionNum := 10
	command := "CANCEL_SET_BUY"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		funds)
	logging(userCommand, file)

	//put money back
	funds := account.SetBuyMap[stock]
	account.unholdMoney(funds)
	//cancel SET BUYs
	delete(account.SetBuyMap, stock)
	//cancel the trigger
	delete(account.BuyTriggers, stock)

}

func setBuyTrigger(
	account *Account, stock string, price float64, file *os.File) {

	server := "CLT1"
	transactionNum := 11
	command := "SET_BUY_TRIGGER"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		account.SetBuyMap[stock])
	logging(userCommand, file)

	//check for set buy on that stock
	if _, ok := account.SetBuyMap[stock]; ok {
		account.BuyTriggers[stock] = price
		glog.Info("Set BUY trigger for ", stock, "at price ", price)

	} else {
		errMsg := fmt.Sprintf("You have to SET BUY AMOUNT on stock %s first", stock)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			account.SetBuyMap[stock],
			errMsg)
		logging(errorEvent, file)
	}
}

func setSellAmount(
	account *Account, stock string, amount float64, file *os.File) {

	server := "CLT1"
	transactionNum := 12
	command := "SET_SELL_AMOUNT"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		amount)
	logging(userCommand, file)

	if account.StockPortfolio[stock] > amount {
		account.SetSellMap[stock] += amount
		//hold stock
		account.holdStock(stock, amount)

	} else {

		errMsg := fmt.Sprintf(
			"User %s does not have enough stock %s to sell",
			username, stock)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			amount,
			errMsg)
		logging(errorEvent, file)
	}
}

func setSellTrigger(
	account *Account, stock string, price float64, file *os.File) {

	server := "CLT1"
	transactionNum := 13
	command := "SET_SELL_TRIGGER"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		account.SetSellMap[stock])
	logging(userCommand, file)

	if _, ok := account.SetSellMap[stock]; ok {
		account.SellTriggers[stock] = price
		glog.Info("Set SELL trigger for ", stock, "at price ", price)

	} else {

		errMsg := fmt.Sprintf(
			"User %s has not SET SELL AMOUNT on stock %s",
			username, stock)
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			username,
			stock,
			0,
			errMsg)
		logging(errorEvent, file)
	}

}

func cancelSetSell(
	account *Account, stock string, amount float64, file *os.File) {

	server := "CLT1"
	transactionNum := 14
	command := "CANCEL_SET_SELL"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		amount)
	logging(userCommand, file)

	funds := account.SetSellMap[stock]

	//put stock back
	account.unholdStock(stock, funds)
	//cancel SET SELLs
	delete(account.SetSellMap, stock)
	//cancel the trigger
	delete(account.SellTriggers, stock)
}

func dumplog(account *Account, filename string) {}

func dumplogAll(filename string) {}
