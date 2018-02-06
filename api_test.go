package main

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const (
	TEST_URL                     = "http://localhost:8082/api/test"
	TEST_ACCOUNT_TRANSACTION_URL = "http://localhost:8082/api/accounttransaction"
)

func activateHttpmock(url string) {

	httpmock.Activate()

	httpmock.RegisterResponder(
		"POST",
		url,
		httpmock.NewStringResponder(200, "ok"))
}

func TestAdd(t *testing.T) {

	activateHttpmock(TEST_ACCOUNT_TRANSACTION_URL)

	account := initializeAccount("test123")
	amount := 100.01
	transactionNum := 2
	add(&account, amount, transactionNum)
}
