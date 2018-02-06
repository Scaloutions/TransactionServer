package main

/*
	TODO:
	getQuote
*/

import (
	"testing"

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

func TestAdd(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccount("test123")
	amount := 100.01
	transactionNum := 1
	add(&account, amount, transactionNum)
}

func TestBuyWithoutQS(t *testing.T) {

	activateMockAuditServer()
	defer httpmock.DeactivateAndReset()

	account := initializeAccount("test123")
	amount := 64.00
	stock := "S"
	stockNum := float64(3)
	transactionNum := 2
	add(&account, amount, 1)
	buyHelper(&account, amount, stock, stockNum, transactionNum)

}
