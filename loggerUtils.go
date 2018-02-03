package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
)

const (
	SERVER_NAME  = "TS0156"
	AUDIT_SERVER = "localhost:9090/log_event"
)

func getFundsAsString(amount float64) string {
	return fmt.Sprintf("%.2f", float64(amount))
}

func getCurrentTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func getUserCommand(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64) UserCommand {

	fundsAsString := getFundsAsString(funds)

	return UserCommand{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         userId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
	}
}

func getTransactionEvent(
	transactionNum int,
	action string,
	userId string,
	funds float64) AccountTransaction {

	fundsAsString := getFundsAsString(funds)

	return AccountTransaction{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Action:         action,
		UserId:         userId,
		Funds:          fundsAsString,
	}
}

func getSystemEvent(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64) SystemEvent {

	fundsAsString := getFundsAsString(funds)

	return SystemEvent{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         userId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
	}
}

func getQuoteServer(
	transactionNum int,
	quoteServerTime int64,
	command string,
	userId string,
	stockSymbol string,
	price float64,
	cryptokey string) QuoteServer {

	priceAsString := getFundsAsString(price)

	return QuoteServer{
		Timestamp:       getCurrentTs(),
		Server:          SERVER_NAME,
		TransactionNum:  transactionNum,
		QuoteServerTime: quoteServerTime,
		Command:         command,
		UserId:          userId,
		StockSymbol:     stockSymbol,
		Price:           priceAsString,
		Cryptokey:       cryptokey,
	}
}

func getErrorEvent(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64,
	errorMessage string) ErrorEvent {

	fundsAsString := getFundsAsString(funds)

	return ErrorEvent{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         userId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
		ErrorMessage:   errorMessage,
	}
}

func logEvent(log interface{}) {
	data, err := json.Marshal(log)
	if err != nil {
		glog.Error("Can not parse struct into JSON onject ", data)
	}
	req, err := http.NewRequest("POST", AUDIT_SERVER, bytes.NewBuffer(data))
	if err != nil {
		glog.Error("Error creating a request for the Audit server")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		glog.Error("Error sending a POST request to Audit server")
		panic(err)
	}
	defer resp.Body.Close()
}
