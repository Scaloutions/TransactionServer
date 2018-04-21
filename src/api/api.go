package api

import (
	"github.com/golang/glog"
	"../db"
	"errors"
	"fmt"
)

const (
	ADD             = "ADD"
	BUY             = "BUY"
	SELL            = "SELL"
	QUOTE			= "QUOTE"
	COMMIT_BUY      = "COMMIT_BUY"
	COMMIT_SELL     = "COMMIT_SELL"
	CANCEL_BUY      = "CANCEL_BUY"
	CANCEL_SELL     = "CANCEL_SELL"
	SET_BUY_AMOUNT  = "SET_BUY_AMOUNT"
	SET_SELL_TRIGGER = "SET_SELL_TRIGGER"
	SET_BUY_TRIGGER = "SET_BUY_TRIGGER"
	SET_SELL_AMOUNT = "SET_SELL_AMOUNT"
	CANCEL_SET_BUY  = "CANCEL_SET_BUY"
	CANCEL_SET_SELL = "CANCEL_SET_SELL"
	DUMPLOG 		= "DUMPLOG"
	DISPLAY_SUMMARY = "DISPLAY_SUMMARY"
)


func Add(account *Account, amount float64, transactionNum int) error {
	if amount > 0 {
		account.addMoney(amount)
		//log transaction event
		log := getTransactionEvent(transactionNum, ADD, account.AccountNumber, amount)
		go logEvent(log)
		glog.Info("SUCCESS: Added ", amount)
		return nil
	} else {
		err := fmt.Sprintf("ERROR: Cannot add zero or negative amount to balance %s", amount)
		glog.Error(err)
		return errors.New("Cannot execute ADD: " + err)
	}
}

func GetQuote(stock string, userId string, transactionNum int) (float64, error) {
	//check cache
	cacheq, err := GetFromCache(stock)

	// found stock in the cache
	if err == nil {
		glog.Info("Got QUOTE from Redis: ", cacheq)
		// log system event
		log := getSystemEvent(transactionNum, QUOTE, userId, stock, cacheq.Price)
		go logEvent(log)
		glog.Info("LOGGING ######## ", log)

		return cacheq.Price, nil

	}

	quoteObj, err := getQuoteFromQS(userId, stock)

	glog.Info("Getting quote for: ", stock, " user: ", userId)
	if err!= nil {
		//TODO : log error event here
		// log := getErrorEvent()
		glog.Error("Failed to get Quote from the QS")
		return 0.0, err
	}

	// put it in CACHE
	glog.Info("Putting new Stock Quote into Redis Cache ", quoteObj)
	err = SetToCache(quoteObj)
	if err != nil {
		glog.Error("Error putting QUOTE into Redist cache ", quoteObj)
	}
	//LOG event as system



	log := getSystemEvent(transactionNum, QUOTE, userId, stock, quoteObj.Price)
	go logEvent(log)
	glog.Info("LOGGING ######## ", log)
	// log := getQuoteServerEvent(transactionNum, quoteObj.Timestamp, quoteObj.UserId, quoteObj.Stock, quoteObj.Price, quoteObj.CryptoKey)
	// go logEvent(log)
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
		go logEvent(log)
		return errors.New("Cannot execute BUY: " + err)
	} else {
		// pull curr stock value for that user
		// currAmount, err := db.GetUserStockAmount(account.AccountNumber, stock)
		// account.StockPortfolio[stock] = currAmount

		// if err!=nil {
			// glog.Error(err, " ", account)
			// return errors.New("Cannot execute BUY: " + err.Error())
		// }

		transaction := db.BuyObj{
			UserId:		 account.AccountNumber,
			Stock:       stock,
			MoneyAmount: amount,
			StockAmount: stockNum,
			TransactionNum: transactionNum,
		}
		//add buy transcation to the stack
		// account.BuyStack.Push(transaction)
		//hold the money
		err := db.CreateNewBuy(transaction)
		if err!=nil {
			//TODO: log error to audit server
			glog.Error(err, " ", account)
			return errors.New("Cannot execute BUY in the DB: " + err.Error())
		}

		account.holdMoney(amount)

		log := getSystemEvent(transactionNum, BUY, account.AccountNumber, stock, amount)
		go logEvent(log)
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
		transaction := db.SellObj {
			UserId:		 account.AccountNumber,
			Stock:       stock,
			MoneyAmount: amount,
			StockAmount: stockNum,
			TransactionNum: transactionNum,
		}
		//this is fine becasue commit transaction has to be executed within 60sec
		//which means that the qoute does not change
		// account.SellStack.Push(transaction)
		err := account.holdStock(stock, stockNum)
		if err !=nil {
			return err
		}

		err = db.CreateNewSell(transaction)
		/// ************!!!!!@$%#@
		//@@@@@@@@ HOLD THE STOCKKKKK
		if err!=nil {
			//TODO: log error to audit server
			glog.Error(err, " ", account)
			return errors.New("Cannot execute CREATE SELL in the DB: " + err.Error())
		}

		log := getSystemEvent(transactionNum, SELL, account.AccountNumber, stock, amount)
		go logEvent(log)
		glog.Info("Executed SELL for ", amount)
		return nil
	} else {
		err := "User does not have enough stock " + stock + " to sell."
		glog.Info(err)
		log := getErrorEvent(transactionNum, SELL, account.AccountNumber, stock, amount, err)
		go logEvent(log)
		return errors.New("Cannot execute SELL: " + err)
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
	transaction, err := db.GetBuy(account.AccountNumber)
	// if account.BuyStack.Size() > 0 {
	if err==nil {
		//go casting
		// i := account.BuyStack.Pop()
		// transaction := i.(BuyObj)
		account.updateBalance(-1*transaction.MoneyAmount)
		//add number of stocks to user
		// account.StockPortfolio[transaction.Stock] += transaction.StockAmount
		//update db record
		db.DeleteBuy(account.AccountNumber)
		err := db.AddUserStock(account.AccountNumber, transaction.Stock, transaction.StockAmount)

		if err!=nil {
			glog.Error(err, " for account:", account)
			return errors.New("Cannot execute COMMIT BUY: " + err.Error())
		}

		log := getTransactionEvent(transactionNum, COMMIT_BUY, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("SUCCESS: Executed COMMIT BUY")
		go logEvent(log)
		return nil

	} else {
		err := "No BUY transactions previously set for account."
		//TODO: figure out if we can simplify this logging with some missing parameters
		log := getErrorEvent(transactionNum, COMMIT_BUY, account.AccountNumber, "", 0, err)
		glog.Error("ERROR: No BUY transactions previously set for account: ", account.AccountNumber)
		go logEvent(log)
		return errors.New("Cannot execute COMMIT BUY: " + err)
	}
}

func CancelBuy(account *Account, transactionNum int) error {
	buy, err := db.GetBuy(account.AccountNumber)
	// if account.BuyStack.Size() > 0 {
	if err==nil {
		// i := account.BuyStack.Pop()
		// transaction := i.(BuyObj)
		//add money back to Available Balance
		err = db.DeleteBuy(account.AccountNumber)
		account.unholdMoney(buy.MoneyAmount)
		if err!=nil {
			msg := "Cannot execute CANCEL BUY: " + err.Error()
			glog.Error(msg)
			return errors.New(msg)
		}
		glog.Info("Executed CANCEL BUY")

		// log := getSystemEvent(transactionNum, CANCEL_BUY, account.AccountNumber, transaction.Stock, transaction.MoneyAmount)
		// go logEvent(log)
		return nil
	} else {
		err := "There are no BUY transcations to cancel for account " + account.AccountNumber
		log := getErrorEvent(transactionNum, CANCEL_BUY, account.AccountNumber, "", 0, err)
		glog.Error(err)
		go logEvent(log)
		return errors.New("Cannot execute CANCEL BUY: " + err)
	}
}

func CommitSell(account *Account, transactionNum int) error {
	transaction, err :=  db.GetSell(account.AccountNumber)
	// if account.SellStack.Size() > 0 {
	if err==nil {
		// i := account.SellStack.Pop()
		// transaction := i.(SellObj)
		//update db record
		err := db.UpdateUserStock(account.AccountNumber, transaction.Stock, transaction.StockAmount*-1)

		if err!=nil {
			glog.Error(err, " for account:", account)
			return errors.New("Cannot execute COMMIT SELL" + err.Error())
		}

		account.addMoney(transaction.MoneyAmount)
		err = db.DeleteSell(account.AccountNumber)

		if err!=nil {
			glog.Error(err, " for account:", account)
			return errors.New("Cannot execute COMMIT SELL" + err.Error())
		}

		log := getTransactionEvent(transactionNum, COMMIT_SELL, account.AccountNumber, transaction.MoneyAmount)
		glog.Info("Executed COMMIT SELL")
		go logEvent(log)
		return nil
	} else {
		err := "No SELL transactions previously set for account"
		glog.Error(err, " ", account.AccountNumber)
		log := getErrorEvent(transactionNum, COMMIT_SELL, account.AccountNumber, "", 0, err)
		go logEvent(log)
		return errors.New("Cannot execute COMMIT SELL: " + err)
	}
}

func CancelSell(account *Account, transactionNum int) error {
	sell, err :=  db.GetSell(account.AccountNumber)

	// if account.SellStack.Size() > 0 {
	if err == nil {
		// i := account.SellStack.Pop()
		// transaction := i.(SellObj)
		err = account.unholdStock(sell.Stock, sell.StockAmount)
		if err !=nil {
			return err
		}

		db.DeleteSell(account.AccountNumber)
		glog.Info("Executed CANCEL SELL")

		log := getSystemEvent(transactionNum, CANCEL_SELL, account.AccountNumber, sell.Stock, sell.MoneyAmount)
		go logEvent(log)
		return nil
	} else {
		err := "There are no SELL transcations to cancel for account " + account.AccountNumber
		glog.Error(err)
		log := getErrorEvent(transactionNum, CANCEL_SELL, account.AccountNumber, "", 0, err)
		go logEvent(log)
		return errors.New("Cannot execute CANCEL SELL: " + err)		
	}
}

/*
Sets a defined amount of the given stock to buy when the current stock price
is less than or equal to the BUY_TRIGGER
*/
func SetBuyAmount(account *Account, stock string, amount float64, transactionNum int) error {
	//check if there is enough money in the account
	//TODO:  CHECK DB
	available := account.getBalance()
	// if account.Available >= amount {
	if available >= amount {
		//hold money
		account.holdMoney(amount)
		err := db.AddSetBuy(account.AccountNumber, stock, amount)
		if err!=nil {
			glog.Error("DB ADD SET BUY failed for ", account)
			return err
		}
		// account.SetBuyMap[stock] += amount

		log := getSystemEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, stock, amount)
		go logEvent(log)
		glog.Info("Executed SET BUY for $", amount, " and stock ", stock)
		glog.Info("Total SET BUY on stock ", stock, " is now ", account.SetBuyMap[stock])
		return nil
	} else {
		err := "Account does not have enough money to buy stock"
		log := getErrorEvent(transactionNum, SET_BUY_AMOUNT, account.AccountNumber, "", 0, err)
		go logEvent(log)
		glog.Error(err, " ", stock)
		return errors.New("Cannot execute SET BUY: " + err)
	}
}

/* Cancels SET BUY associated with a particular stock
   TODO: verify what happens if the user set multiple SET BUY on one stock
   It shouldbe overwritten by the most recent one!
   TODO: fix this.
*/
func CancelSetBuy(account *Account, stock string, transactionNum int) error {
	setBuy, err := db.GetSetBuy(account.AccountNumber, stock)
	// if val, ok := account.SetBuyMap[stock]; ok {
	if err == nil {
		//put money back
		account.unholdMoney(setBuy.MoneyAmount)
		//cancel SET BUYs
		err = db.DeleteSetBuy(account.AccountNumber, stock)
		// delete(account.SetBuyMap, stock)
		//cancel the trigger
		// delete(account.BuyTriggers, stock)
		if err!=nil {
			glog.Info("Error deleting SET BUY for ", account, " stock: ", stock)
			return err
		}

		//TODO: check if we need to pass val here for logging
		log := getSystemEvent(transactionNum, CANCEL_SET_BUY, account.AccountNumber, stock, setBuy.MoneyAmount)
		go logEvent(log)
		glog.Info("Executed CANCEL SET BUY")
		return nil
	} else {
		err := "No SET BUY AMOUNT was previously set for this account"
		log := getErrorEvent(transactionNum, CANCEL_SET_BUY, account.AccountNumber, "", 0, err)
		go logEvent(log)
		glog.Error(err, " ", account.AccountNumber)
		return errors.New("Cannot execute CANCEL SET BUY: " + err)
	}
}

func SetBuyTrigger(account *Account, stock string, price float64, transactionNum int) error {
	//check for set buy on that stock
	setBuy, err := db.GetSetBuy(account.AccountNumber, stock)
	// if _, ok := account.SetBuyMap[stock]; ok {
	if err==nil {
		// if _, exists := account.BuyTriggers[stock]; exists {
		//TODO: Upadte Trigger value
		if setBuy.RunningTrigger {
			glog.Info("Trigger is already running!")
			// account.BuyTriggers[stock] = price
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up New SetBuy Trigger")
			//prevent race condition here TODO: rewrite
			// account.BuyTriggers[stock] = price
			//TODO: check for error and backpropogate it
			go account.startBuyTrigger(stock, price, transactionNum)
		}

		glog.Info("Set BUY trigger for ", stock, " at price ", price)
		log := getUserCmndEvent(transactionNum, SET_BUY_TRIGGER, account.AccountNumber, stock, price)
		go logEvent(log)
		return nil
	} else {
		err := "You have to SET BUY AMOUNT on stock " + stock + " first."
		glog.Error(err)
		log := getErrorEvent(transactionNum, SET_BUY_TRIGGER, account.AccountNumber, stock, price, err)
		go logEvent(log)
		return errors.New("Cannot execute SET BUY TRIGGER: " + err)		
	}
}

func SetSellAmount(account *Account, stock string, amount float64, transactionNum int) error {
	if account.hasStock(stock, amount) {
	// if account.StockPortfolio[stock] > amount {
		// account.SetSellMap[stock] += amount
		//hold stock
		account.holdStock(stock, amount)
		db.AddSetSell(account.AccountNumber, stock, amount)

		log := getSystemEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, stock, amount)
		go logEvent(log)
		glog.Info("Executed SET SELL AMOUNT for ", amount)
		return nil
	} else {
		err := "Account does not have enough stock to sell "
		log := getErrorEvent(transactionNum, SET_SELL_AMOUNT, account.AccountNumber, "", 0, err)
		go logEvent(log)
		glog.Error(err, " ", stock)
		return errors.New("Cannot execute SET SELL: " + err)		
	}
}

func SetSellTrigger(account *Account, stock string, price float64, transactionNum int) error {
	//check for set buy on that stock
	setSell, err := db.GetSetSell(account.AccountNumber, stock)
	// if _, ok := account.SetSellMap[stock]; ok {
	if err==nil {
		// if _, exists := account.SellTriggers[stock]; exists {
		if setSell.RunningTrigger {
			glog.Info("Sell Trigger is already running!")
			// account.SellTriggers[stock] = price
		} else {
			//spin up go routine trigger
			glog.Info("Spinning up SEll trigger")
			// account.SellTriggers[stock] = price
			//TODO: check for error and backpropogate it
			go account.startSellTrigger(stock, price, transactionNum)
		}
		// assuming running trigger is not an error
		glog.Info("Set SELL trigger for ", stock, " at price ", price)
		log := getUserCmndEvent(transactionNum, SET_SELL_TRIGGER, account.AccountNumber, stock, price)
		go logEvent(log)
		return nil
	} else {
		//TODO: properly log this error
		err := "You have to SET SELL AMOUNT on stock " + stock + " first."
		glog.Error(err)
		log := getErrorEvent(transactionNum, SET_SELL_TRIGGER, account.AccountNumber, stock, price, err)
		go logEvent(log)
		return errors.New("Cannot execute SET SELL TRIGGER: " + err)		
	}
}

func CancelSetSell(account *Account, stock string, transactionNum int) error {
	setSell, err := db.GetSetSell(account.AccountNumber, stock)
	// if val, ok := account.SetSellMap[stock]; ok {
	if err==nil {
		//put stock back
		account.holdStock(stock, setSell.StockAmount)
		//cancel SET SELLs
		err = db.DeleteSetBuy(account.AccountNumber, stock)
		//cancel the trigger

		log := getSystemEvent(transactionNum, CANCEL_SET_SELL, account.AccountNumber, stock, setSell.StockAmount)
		go logEvent(log)
		if err!=nil {
			return err
		}
		glog.Info("Executed CANCEL SET SELL")
		return nil
	} else {
		err := "No SET SELL AMOUNT was previously set for this account"
		log := getErrorEvent(transactionNum, CANCEL_SET_SELL, account.AccountNumber, "", 0, err)
		go logEvent(log)
		glog.Error(err, " ", account.AccountNumber)
		return errors.New("Cannot execute CANCEL SET SELL: " + err)		
	}
}

func Dumplog(transactionNum int, userId string) {

	glog.Info("Processing Display DUMPLOG Request....")
	log := getUserCmndEvent(transactionNum, DUMPLOG, userId, "NONE", 0.0)
	go logEvent(log)
	// send dumplog request
	go getDumplog()
}

func DisplaySummary(transactionNum int,
	userId string,
	stockSymbol string,
	funds float64) {

	glog.Info("Processing Display Summary Request....")
	log := getUserCmndEvent(transactionNum, DISPLAY_SUMMARY, userId, stockSymbol, funds)
	go logEvent(log)
}
