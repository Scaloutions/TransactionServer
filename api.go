package main

//below are all the functions that need to be implemented in the system

func add(acoount string, amount float64) {}

func quote(stock string) float64 { 
	return 0
}

func buy(account string, stock string, amount float64) {}

func sell(account string, stock string, amount float64) {}

func commitBuy(account string) {} 

func cancelBuy(account string) {} 

func commitSell(account string) {} 

func cancelSell(account string) {} 

func setBuyAmount(account string, stock string, amount float64) {}

func cancelSetBuy(accont string, storck string) {}

func setBuyTrigger(account string, stock string, amount float64) {}

func setSellAmount(account string, stock string, amount float64) {}

func setSellTrigger(account string, stock string, amount float64) {}

func cancelSetSell(account string, stock string, amount float64) {}

func dumplog(account string, filename string) {}

func dumplogAll(filename string) {}


