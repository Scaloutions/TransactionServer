package main

/*
	TODO:
	getQuote
	sell
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	TEST_URL                     = "http://localhost:8082/api/test"
	TEST_ACCOUNT_TRANSACTION_URL = "http://localhost:8082/api/accounttransaction"
	TEST_SYSTEM_EVENT_URL        = "http://localhost:8082/api/systemevent"
)

func activateHttpmock(url string) {

	httpmock.Activate()

	httpmock.RegisterResponder(
		"POST",
		url,
		httpmock.NewStringResponder(200, "ok"))
}

func activateMockAuditServer() {
	activateHttpmock(TEST_SYSTEM_EVENT_URL)
	activateHttpmock(TEST_ACCOUNT_TRANSACTION_URL)
}

func initializeAccountForTesting(amount float64) *Account {

	account := initializeAccount("test123")
	transactionNum := 1
	add(&account, amount, transactionNum)

	return &account
}

func buyStockForTesting(account *Account) *Account {
	amount := float64(64)
	stock := "S"
	transactionNum := 2
	stockNum := float64(4)
	buyHelper(account, amount, stock, stockNum, transactionNum)
	return account
}

func TestAdd(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	amount := 100.01
	account := initializeAccountForTesting(amount)
	assert.Equal(t, amount, account.Available)
	assert.Equal(t, amount, account.Balance)
}

func TestBuyWithoutQS(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	amount := 64.00
	account := initializeAccountForTesting(amount)
	assert.Equal(t, amount, account.Balance)
	assert.Equal(t, amount, account.Available)
	targetAmount := float64(0)
	stock := "S"
	stockNum := float64(10)
	transactionNum := 2
	buyHelper(account, amount, stock, stockNum, transactionNum)
	assert.Equal(t, targetAmount, account.Available)
	assert.Equal(t, amount, account.Balance) // 0 only when buy operation is committed
	assert.False(t, account.hasStock(stock, stockNum))
	// has stock only after buy is committed

}

func TestCommitBuy(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccountForTesting(100)
	buyHelper(account, float64(64), "S", float64(4), 2)
	assert.Equal(t, float64(100), account.Balance)

	commitBuy(account, 3)
	assert.Equal(t, float64(36), account.Balance)

}