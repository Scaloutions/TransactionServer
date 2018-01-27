package main

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

//below are all the functions that need to be implemented in the system

const (
	addAction    = "add"
	removeAction = "remove"
)

func add(
	account *Account,
	amount float64,
	f *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
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
		glog.Info("SUCCESS: Added ", amount)

		accountTransaction := getAccountTransaction(
			server,
			transactionNum,
			addAction,
			username,
			amount)
		logging(accountTransaction, f)

	} else {
		glog.Error("ERROR: Cannot add negative amount to balance ", amount)
		errMsg := fmt.Sprintf("Cannot add negative amount %.2f to balance ", amount)
		// glog.Error(errMsg)

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

func getQuote(
	stock string,
	userid string,
	file *os.File,
	transactionNum int,
	command string) float64 {
	quoteObj := getQuoteFromQS(userid, stock)

	server := "CLT1"

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		userid,
		stock,
		0)
	logging(userCommand, file)

	//TODO: log quote server hit here
	return quoteObj.Price
	// return 1
}

func buy(
	account *Account,
	stock string,
	amount float64,
	file *os.File,
	transactionNum int,
	command string) {

	//get qoute
	stockNum := amount / getQuote(stock, account.AccountNumber, file, transactionNum, "QUOTE")

	server := "CLT1"
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
			glog.Info("WARNING: Not enough money for account ",account.AccountNumber, " to buy ", stock)

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
		glog.Info("SUCCESS: Executed BUY for ", amount)
	}
}

func sell(account *Account, stock string, amount float64,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
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
	stockNum := amount / getQuote(stock, account.AccountNumber, file, transactionNum, "QUOTE")
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
		glog.Info("Executed SELL for ", amount)

	} else {
		errMsg := fmt.Sprintf("Not enough stock %s to sell.", stock)
		glog.Info("WARNING: Not enough stock ", stock, " to sell.")

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

func commitBuy(account *Account, file *os.File, transactionNum int, command string) {

	server := "CLT1"
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

		funds := transaction.MoneyAmount
		//should we check balance here insted? TODO: clarify
		account.Balance -= funds

		accountTransaction := getAccountTransaction(
			server,
			transactionNum,
			removeAction,
			username,
			funds)
		logging(accountTransaction, file)

		//add number of stocks to user
		//TODO: refactor this line
		account.StockPortfolio[transaction.Stock] += transaction.StockAmount
		glog.Info("SUCCESS: Executed COMMIT BUY")

	} else {

		errMsg := fmt.Sprintf("No BUY transactions previously set for account %s", username)
		gglog.Error("ERROR: No BUY transactions previously set for account: ", account.AccountNumber)

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

func cancelBuy(account *Account, file *os.File, transactionNum int, command string) {

	server := "CLT1"
	username := account.AccountNumber

	//TODO: log this
	if account.BuyStack.size > 0 {
		 i := account.BuyStack.Pop()
		 transaction := i.(Buy)
		 //add money back to Available Balance
		account.unholdMoney(transaction.MoneyAmount)
		glog.Info("Executed CANCEL BUY")

		userCommand := getUserCommand(
			server,
			transactionNum,
			command,
			username,
			transaction.Stock,
			transaction.MoneyAmount)
		logging(userCommand, file)
	} else {
		glog.Error("No BUY transactions previously set for account: ", account.AccountNumber)
	}
} 

func commitSell(account *Account, file *os.File, transactionNum int, command string) {

	server := "CLT1"
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

		glog.Info("Executed COMMIT SELL")

		accountTransaction := getAccountTransaction(
			server,
			transactionNum,
			addAction,
			username,
			funds)
		logging(accountTransaction, file)

	} else {

		errMsg := fmt.Sprintf(
			"No SELL transactions previously set for account %s", username)
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

func cancelSell(account *Account, file *os.File, transactionNum int, command string) {

	server := "CLT1"
	username := account.AccountNumber

	//TODO: log this
	glog.Info("Executing CANCEL SELL")
	if account.SellStack.size > 0 {
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
		account.unholdStock(transaction.Stock, transaction.StockAmount)
		glog.Info("Executed CANCEL SELL")
	} else {
		glog.Error("There are no SELL commands associated with this user.")
	}
} 

/*
Sets a defined amount of the given stock to buy when the current stock price
is less than or equal to the BUY_TRIGGER
*/
func setBuyAmount(account *Account, stock string, amount float64,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
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
		glog.Info("Executed SET BUY for $", amount, " and stock ", stock)
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
func cancelSetBuy(account *Account, stock string,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
	username := account.AccountNumber

	//put money back
	funds := account.SetBuyMap[stock]
	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		funds)
	logging(userCommand, file)

	account.unholdMoney(funds)

	//cancel SET BUYs
	delete(account.SetBuyMap, stock)
	//cancel the trigger
	delete(account.BuyTriggers, stock)
	glog.Info("Executed CANCEL SET BUY")
}

func setBuyTrigger(account *Account, stock string, price float64,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
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
		//TODO: this is hacky we need to properly check for the key here
		if _, exists := account.BuyTriggers[stock]; exists {
			glog.Info("Trigger is already running!")
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up go routine")
			//prevent race condition here TODO: rewrite
			account.BuyTriggers[stock] = price
			go account.startBuyTrigger(stock)
		}

		account.BuyTriggers[stock] = price
		glog.Info("Set BUY trigger for ", stock, " at price ", price)
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

func setSellAmount(account *Account, stock string, amount float64,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
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
		glog.Info("Executed SET SELL AMOUNT for ", amount)
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

func setSellTrigger(account *Account, stock string, price float64,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
	username := account.AccountNumber

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		account.SetSellMap[stock])
	logging(userCommand, file)

	//check for set buy on that stock
	if _, ok := account.SetSellMap[stock]; ok {
		//TODO: this is hacky we need to properly check for the key here
		if _, exists := account.SellTriggers[stock]; exists {
			glog.Info("Sell Trigger is already running!")
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up go routine SEll trigger")
			account.SellTriggers[stock] = price
			go account.startSellTrigger(stock)
		}

		account.SellTriggers[stock] = price
		glog.Info("Set SELL trigger for ", stock, " at price ", price)
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

// func setSellTrigger(account *Account, stock string, price float64) {
// 	if _, ok := account.SetSellMap[stock]; ok {
// 		account.SellTriggers[stock] = price
// 		glog.Info("Executed SET SELL trigger for ", stock, "at price ", price)
// 	} else {
// 		glog.Error("You have to SET SELL AMOUNT on stock ", stock, " first.")
// 	}

// }

func cancelSetSell(account *Account, stock string,
	file *os.File,
	transactionNum int,
	command string) {

	server := "CLT1"
	username := account.AccountNumber

	funds := account.SetSellMap[stock]

	userCommand := getUserCommand(
		server,
		transactionNum,
		command,
		username,
		stock,
		funds)
	logging(userCommand, file)

	//put stock back
	account.unholdStock(stock, funds)
	//cancel SET SELLs
	delete(account.SetSellMap, stock)
	//cancel the trigger
	delete(account.SellTriggers, stock)
	glog.Info("Executed CANCEL SET SELL")
}

func dumplog(account *Account, filename string) {}

func dumplogAll(filename string) {}
