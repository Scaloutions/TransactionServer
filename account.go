
package main

type Account struct {
	AccountNumber string
	Balance float64
}

func initializeAccount() Account {
	return Account{}
}

func getMoney(account *Account) float64{
	return account.Balance
}

func addMoney(account *Account, amount float64){
	account.Balance += amount
}