package main

import (
	"fmt"
	"time"
)

func getFundsAsString(amount float64) string {
	if amount == 0 {
		return ""
	}
	return fmt.Sprintf("%.2f", float64(amount))
}

func getCurrentTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func getUserCommand(
	server string,
	transactionNum int,
	command string,
	UserId string,
	stockSymbol string,
	funds float64) UserCommand {

	fundsAsString := getFundsAsString(funds)

	return UserCommand{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         UserId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString}
}

func getAccountTransaction(
	server string,
	transactionNum int,
	action string,
	UserId string,
	funds float64) AccountTransaction {

	fundsAsString := getFundsAsString(funds)

	return AccountTransaction{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Action:         action,
		UserId:         UserId,
		Funds:          fundsAsString}
}

func getSystemEvent(
	server string,
	transactionNum int,
	command string,
	UserId string,
	stockSymbol string,
	funds float64) SystemEvent {

	fundsAsString := getFundsAsString(funds)

	return SystemEvent{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         UserId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString}
}

func getQuoteServer(
	server string,
	transactionNum int,
	quoteServerTime int64,
	command string,
	UserId string,
	stockSymbol string,
	price float64,
	cryptokey string) QuoteServer {

	priceAsString := getFundsAsString(price)

	return QuoteServer{
		Timestamp:       getCurrentTs(),
		Server:          server,
		TransactionNum:  transactionNum,
		QuoteServerTime: quoteServerTime,
		Command:         command,
		UserId:          UserId,
		StockSymbol:     stockSymbol,
		Price:           priceAsString,
		Cryptokey:       cryptokey}
}

func getErrorEvent(
	server string,
	transactionNum int,
	command string,
	UserId string,
	stockSymbol string,
	funds float64,
	errorMessage string) ErrorEvent {

	fundsAsString := getFundsAsString(funds)

	return ErrorEvent{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         UserId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
		ErrorMessage:   errorMessage}
}
