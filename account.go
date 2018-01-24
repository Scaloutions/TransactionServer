
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

func (account *Account) getMoney() float64{
	return account.Balance
}

func (account *Account) addMoney( amount float64){
	account.Balance += amount
}