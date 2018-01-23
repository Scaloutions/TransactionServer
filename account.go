
package main

type Account struct {
	AccountNumber string
	Balance float64
}

func initializeAccount(value string) Account {
	return Account{
		AccountNumber: value,
		Balance: 0.0,
	}
}

func getMoney(account *Account) float64{
	return account.Balance
}

func addMoney(account *Account, amount float64){
	account.Balance += amount
}