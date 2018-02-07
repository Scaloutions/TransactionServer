package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
)

const (
	SERVER_NAME       = "TS0156"
	AUDIT_SERVER      = "http://localhost:8082"
	API_URL           = "/api"
	ACCOUNT_EVENT_URL = "/accounttransaction"
	SYSTEM_EVENT_URL  = "/systemevent"
	ERROR_EVENT_URL   = "/errorevent"
	QUOTE_SERVER_URL  = "/quoteserver"
)

func getFundsAsString(amount float64) string {
	return fmt.Sprintf("%.2f", float64(amount))
}

func getCurrentTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func getTransactionEvent(
	transactionNum int,
	action string,
	userId string,
	funds float64) AccountTransactionEvent {

	fundsAsString := getFundsAsString(funds)

	return AccountTransactionEvent{
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

func getQuoteServerEvent(
	transactionNum int,
	quoteServerTime int64,
	command string,
	userId string,
	stockSymbol string,
	price float64,
	cryptokey string) QuoteServerEvent {

	priceAsString := getFundsAsString(price)

	return QuoteServerEvent{
		Timestamp:            getCurrentTs(),
		Server:               SERVER_NAME,
		TransactionNum:       transactionNum,
		QuoteServerEventTime: quoteServerTime,
		Command:              command,
		UserId:               userId,
		StockSymbol:          stockSymbol,
		Price:                priceAsString,
		Cryptokey:            cryptokey,
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
	URL := getUrlPath(log)

	if err != nil {
		glog.Error("Can not parse struct into JSON onject ", data)
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
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

func getUrlPath(obj interface{}) string {
	var url bytes.Buffer
	url.WriteString(AUDIT_SERVER)
	url.WriteString(API_URL)

	switch obj.(type) {
	case SystemEvent:
		url.WriteString(SYSTEM_EVENT_URL)
	case AccountTransactionEvent:
		url.WriteString(ACCOUNT_EVENT_URL)
	case ErrorEvent:
		url.WriteString(ERROR_EVENT_URL)
	case QuoteServerEvent:
		url.WriteString(QUOTE_SERVER_URL)
	default:
		glog.Error("Error logging event to the audit server.")
		panic("Can not recognaize this type of event.")
	}

	return url.String()
}
