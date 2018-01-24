package main

import "github.com/golang/glog"
//below are all the functions that need to be implemented in the system

func add(account Account, amount float64) {
	if amount > 0{
		account.addMoney(amount)
	} else {
        glog.Error("Cannot add negative amount to balance ", amount)
	}
}

func quote(stock string) float64 { 
	return 0
}

func buy(account Account, stock string, amount float64) {}

func sell(account Account, stock string, amount float64) {}
func commitBuy(account Account) {} 

func cancelBuy(account Account) {} 

func commitSell(account Account) {} 

func cancelSell(account Account) {} 

func setBuyAmount(account Account, stock string, amount float64) {}

func cancelSetBuy(accont Account, storck string) {}

func setBuyTrigger(account Account, stock string, amount float64) {}

func setSellAmount(account Account, stock string, amount float64) {}

func setSellTrigger(account Account, stock string, amount float64) {}

func cancelSetSell(account Account, stock string, amount float64) {}

func dumplog(account Account, filename string) {}

func dumplogAll(filename string) {}


