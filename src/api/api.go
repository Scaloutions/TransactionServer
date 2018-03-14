package api

import (
	"github.com/golang/glog"
	"../db"
	"errors"
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
	CANCEL_SET_BUY  = "cancel_set_buy"
	CANCEL_SET_SELL = "cancel_set_sell"
	QUOTE			= "get_quote"
)

func Add(account *Account, amount float64, transactionNum int) error {
	if amount > 0 {
		account.addMoney(amount)
		//log transaction event
		log := getTransactionEvent(transactionNum, ADD, account.AccountNumber, amount)
		logEvent(log)
		glog.Info("SUCCESS: Added ", amount)
		return nil
	} else {
		glog.Error("ERROR: Cannot add zero or negative amount to balance ", amount)
		return errors.New("Cannot execute ADD")
	}
}

func GetQuote(stock string, userid string, transactionNum int) (float64, error) {
	quoteObj, err := getQuoteFromQS(userid, stock)
	if err!= nil {
		//TODO : log error event here
		// log := getErrorEvent()
		glog.Error("Failed to get Quote from the QS")
		return 0.0, err
	}

	log := getQuoteServerEvent(transactionNum, quoteObj.Timestamp, QUOTE, quoteObj.UserId, quoteObj.Stock, quoteObj.Price, quoteObj.CryptoKey)
	logEvent(log)
	return quoteObj.Price, nil
}

func buyHelper(
	account *Account,
	amount float64,
	stock string,
	stockNum float64,
	transactionNum int) error {

	//check balance
	if account.getBalance() < amount {
		//TODO: improve logging
		err := "Account does not have enough money to execute BUY command"
		glog.Info("Not enough money on account ", account.AccountNumber, " to buy ", stock)
		log := getErrorEvent(transactionNum, BUY, account.AccountNumber, stock, amount, err)
		logEvent(log)
		return errors.New("Cannot execute BUY")
	} else {
		// pull curr stock value for that user
		currAmount, err := db.GetUserStockAmount(account.AccountNumber, stock)
		account.StockPortfolio[stock] = currAmount

		if err!=nil {
			glog.Error(err, " ", account)
			return errors.New("Cannot execute BUY")
		}

		transaction := BuyObj{
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
		return nil
	}

}

func Buy(account *Account, stock string, amount float64, transactionNum int) error {
	//get quote and calculate number of stock
	quote, err := GetQuote(stock, account.AccountNumber, transactionNum)
	if err!= nil {
		return err
	}
	stockNum := amount / quote
	return buyHelper(account, amount, stock, stockNum, transactionNum)
}

func sellHelper(
	account *Account,
	stock string,
	amount float64,
	transactionNum int,
	stockNum float64) error {

	if account.hasStock(stock, stockNum) {
		transaction := SellObj {
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
		return nil
	} else {
		err := "User does not have enough stock to sell."
		glog.Info("Not enough stock ", stock, " to sell.")
		log := getErrorEvent(transactionNum, SELL, account.AccountNumber, stock, amount, err)
		logEvent(log)
		return errors.New("Cannot execute SELL")
	}
}

func Sell(account *Account, stock string, amount float64, transactionNum int) error {
	quote, err := GetQuote(stock, account.AccountNumber, transactionNum)
	if err!= nil {
		return err
	}
	//check if have that # of stocks
	stockNum := amount / quote
	return sellHelper(account, stock, amount, transactionNum, stockNum)
}

func CommitBuy(account *Account, transactionNum int) error {
	if account.BuyStack.Size() > 0 {
		//go casting
		i := account.BuyStack.Pop()
		transaction := i.(BuyObj)
		account.substractBalance(transaction.MoneyAmount)
		//add number of stocks to user
		account.StockPortfolio[transaction.Stock] += transaction.StockAmount
		//update db record
		err := db.UpdateUserStock(account.AccountNumber, transaction.Stock, account.StockPortfolio[transaction.Stock])

		if err!=nil {
			glog.Error(err, " for account:", account)
			return errors.New("Cannot execute COMMIT BUY")
		}

		log := getTransactionEvent(transactionNum, COMMIT_BUY, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("SUCCESS: Executed COMMIT BUY")
		logEvent(log)
		return nil

	} else {
		err := "No BUY transactions previously set for account"
		//TODO: figure out if we can simplify this logging with some missing parameters
		log := getErrorEvent(transactionNum, COMMIT_BUY, account.AccountNumber, "", 0, err)
		glog.Error("ERROR: No BUY transactions previously set for account: ", account.AccountNumber)
		logEvent(log)
		return errors.New("Cannot execute COMMIT BUY")
	}
}

func CancelBuy(account *Account, transactionNum int) error {
	if account.BuyStack.Size() > 0 {
		i := account.BuyStack.Pop()
		transaction := i.(BuyObj)
		//add money back to Available Balance
		account.unholdMoney(transaction.MoneyAmount)
		glog.Info("Executed CANCEL BUY")

		log := getSystemEvent(transactionNum, CANCEL_BUY, account.AccountNumber, transaction.Stock, transaction.MoneyAmount)
		logEvent(log)
		return nil
	} else {
		err := "There are no BUY transcations to cancel for this account"
		log := getErrorEvent(transactionNum, CANCEL_BUY, account.AccountNumber, "", 0, err)
		glog.Error(err, " ", account.AccountNumber)
		logEvent(log)
		return errors.New("Cannot execute CANCEL BUY")
	}
}

func CommitSell(account *Account, transactionNum int) error {
	if account.SellStack.Size() > 0 {
		i := account.SellStack.Pop()
		transaction := i.(SellObj)
		account.addMoney(transaction.MoneyAmount)
		//update db record
		err := db.UpdateUserStock(account.AccountNumber, transaction.Stock, account.StockPortfolio[transaction.Stock])

		if err!=nil {
			glog.Error(err, " for account:", account)
			return errors.New("Cannot execute COMMIT SELL")
		}

		log := getTransactionEvent(transactionNum, COMMIT_SELL, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("Executed COMMIT SELL")
		logEvent(log)
		return nil
	} else {
		err := "No SELL transactions previously set for account"
		glog.Error(err, " ", account.AccountNumber)
		log := getErrorEvent(transactionNum, COMMIT_SELL, account.AccountNumber, "", 0, err)
		logEvent(log)
		return errors.New("Cannot execute COMMIT SELL")
	}
}

func CancelSell(account *Account, transactionNum int) error {
	if account.SellStack.Size() > 0 {
		i := account.SellStack.Pop()
		transaction := i.(SellObj)
		account.unholdStock(transaction.Stock, transaction.StockAmount)
		glog.Info("Executed CANCEL SELL")

		log := getSystemEvent(transactionNum, CANCEL_SELL, account.AccountNumber, transaction.Stock, transaction.MoneyAmount)
		logEvent(log)
		return nil
	} else {
		err := "There are no SELL transcations to cancel for this account"
		glog.Error(err, " ", account.AccountNumber)
		log := getErrorEvent(transactionNum, CANCEL_SELL, account.AccountNumber, "", 0, err)
		logEvent(log)
		return errors.New("Cannot execute CANCEL SELL")		
	}
}

/*
Sets a defined amount of the given stock to buy when the current stock price
is less than or equal to the BUY_TRIGGER
*/
func SetBuyAmount(account *Account, stock string, amount float64, transactionNum int) error {
	//check if there is enough money in the account
	if account.Available >= amount {
		//hold money
		account.holdMoney(amount)
		account.SetBuyMap[stock] += amount

		log := getSystemEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("Executed SET BUY for $", amount, " and stock ", stock)
		glog.Info("Total SET BUY on stock ", stock, " is now ", account.SetBuyMap[stock])
		return nil
	} else {
		err := "Account does not have enough money to buy stock"
		log := getErrorEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", stock)
		return errors.New("Cannot execute SET BUY")
	}
}

/* Cancels SET BUY associated with a particular stock
   TODO: verify what happens if the user set multiple SET BUY on one stock
   It shouldbe overwritten by the most recent one!
   TODO: fix this.
*/
func CancelSetBuy(account *Account, stock string, transactionNum int) error {
	if val, ok := account.SetBuyMap[stock]; ok {
		//put money back
		account.unholdMoney(val)
		//cancel SET BUYs
		delete(account.SetBuyMap, stock)
		//cancel the trigger
		delete(account.BuyTriggers, stock)

		//TODO: check if we need to pass val here for logging
		log := getSystemEvent(transactionNum, CANCEL_SET_BUY, account.AccountNumber, stock, val)
		logEvent(log)
		glog.Info("Executed CANCEL SET BUY")
		return nil
	} else {
		err := "No SET BUY AMOUNT was previously set for this account"
		log := getErrorEvent(transactionNum, CANCEL_SET_BUY, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", account.AccountNumber)
		return errors.New("Cannot execute CANCEL SET BUY")
	}
}

func SetBuyTrigger(account *Account, stock string, price float64, transactionNum int) error {
	//check for set buy on that stock
	if _, ok := account.SetBuyMap[stock]; ok {
		if _, exists := account.BuyTriggers[stock]; exists {
			glog.Info("Trigger is already running!")
			account.BuyTriggers[stock] = price
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up SetBuy Trigger")
			//prevent race condition here TODO: rewrite
			account.BuyTriggers[stock] = price
			//TODO: check for error and backpropogate it
			go account.startBuyTrigger(stock, transactionNum)
		}

		glog.Info("Set BUY trigger for ", stock, " at price ", price)
		return nil
	} else {
		glog.Error("You have to SET BUY AMOUNT on stock ", stock, " first.")
		return errors.New("Cannot execute SET BUY TRIGGER")		
	}
}

func SetSellAmount(account *Account, stock string, amount float64, transactionNum int) error {
	if account.StockPortfolio[stock] > amount {
		account.SetSellMap[stock] += amount
		//hold stock
		account.holdStock(stock, amount)

		log := getSystemEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, stock, amount)
		logEvent(log)
		glog.Info("Executed SET SELL AMOUNT for ", amount)
		return nil
	} else {
		err := "Account does not have enough stock to sell "
		log := getErrorEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", stock)
		return errors.New("Cannot execute SET SELL")		
	}
}

func SetSellTrigger(account *Account, stock string, price float64, transactionNum int) error {
	//check for set buy on that stock
	if _, ok := account.SetSellMap[stock]; ok {
		if _, exists := account.SellTriggers[stock]; exists {
			glog.Info("Sell Trigger is already running!")
			account.SellTriggers[stock] = price
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up SEll trigger")
			account.SellTriggers[stock] = price
			//TODO: check for error and backpropogate it
			go account.startSellTrigger(stock, transactionNum)
		}
		// assuming running trigger is not an error
		glog.Info("Set SELL trigger for ", stock, " at price ", price)
		return nil
	} else {
		//TODO: properly log this error
		glog.Error("You have to SET SELL AMOUNT on stock ", stock, " first.")
		return errors.New("Cannot execute SET SELL TRIGGER")		
	}
}

func CancelSetSell(account *Account, stock string, transactionNum int) error {
	if val, ok := account.SetSellMap[stock]; ok {
		//put stock back
		account.unholdStock(stock, val)
		//cancel SET SELLs
		delete(account.SetSellMap, stock)
		//cancel the trigger
		delete(account.SellTriggers, stock)

		log := getSystemEvent(transactionNum, CANCEL_SET_SELL, account.AccountNumber, stock, val)
		logEvent(log)
		glog.Info("Executed CANCEL SET SELL")
		return nil
	} else {
		err := "No SET SELL AMOUNT was previously set for this account"
		log := getErrorEvent(transactionNum, CANCEL_SET_SELL, account.AccountNumber, "", 0, err)
		logEvent(log)
		glog.Error(err, " ", account.AccountNumber)
		return errors.New("Cannot execute CANCEL SET SELL")		
	}
}
