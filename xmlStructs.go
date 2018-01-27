package main

import (
	"encoding/xml"
)

type UserCommand struct {
	XMLName        xml.Name `xml:"userCommand"`
	Timestamp      int64    `xml:"timestamp,omitempty"`
	Server         string   `xml:"server,omitempty"`
	TransactionNum int      `xml:"transactionNum,omitempty"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type AccountTransaction struct {
	XMLName        xml.Name `xml:"accountTransaction"`
	Timestamp      int64    `xml:"timestamp,omitempty"`
	Server         string   `xml:"server,omitempty"`
	TransactionNum int      `xml:"transactionNum,omitempty"`
	Action         string   `xml:"action,omitempty"`
	Username       string   `xml:"username,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type SystemEvent struct {
	XMLName        xml.Name `xml:"systemEvent"`
	Timestamp      int64    `xml:"timestamp,omitempty"`
	Server         string   `xml:"server,omitempty"`
	TransactionNum int      `xml:"transactionNum,omitempty"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type QuoteServer struct {
	XMLName         xml.Name `xml:"quoteServer"`
	Timestamp       int64    `xml:"timestamp,omitempty"`
	Server          string   `xml:"server,omitempty"`
	TransactionNum  int      `xml:"transactionNum,omitempty"`
	QuoteServerTime int64    `xml:"quoteServerTime,omitempty"`
	Command         string   `xml:"command,omitempty"`
	Username        string   `xml:"username,omitempty"`
	StockSymbol     string   `xml:"stockSymbol,omitempty"`
	Price           string   `xml:"price,omitempty"`
	Cryptokey       string   `xml:"cryptokey,omitempty"`
}

type ErrorEvent struct {
	XMLName        xml.Name `xml:"errorEvent"`
	Timestamp      int64    `xml:"timestamp,omitempty"`
	Server         string   `xml:"server,omitempty"`
	TransactionNum int      `xml:"transactionNum,omitempty"`
	Command        string   `xml:"command,omitempty"`
	Username       string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
	ErrorMessage   string   `xml:"errorMessage,omitempty"`
}

func getUserCommand(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds float64) UserCommand {

	fundsAsString := getFundsAsString(funds)

	return UserCommand{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString}
}

func getAccountTransaction(
	server string,
	transactionNum int,
	action string,
	username string,
	funds float64) AccountTransaction {

	fundsAsString := getFundsAsString(funds)

	return AccountTransaction{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Action:         action,
		Username:       username,
		Funds:          fundsAsString}
}

func getSystemEvent(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds float64) SystemEvent {

	fundsAsString := getFundsAsString(funds)

	return SystemEvent{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString}
}

func getQuoteServer(
	server string,
	transactionNum int,
	quoteServerTime int64,
	command string,
	username string,
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
		Username:        username,
		StockSymbol:     stockSymbol,
		Price:           priceAsString,
		Cryptokey:       cryptokey}
}

func getErrorEvent(
	server string,
	transactionNum int,
	command string,
	username string,
	stockSymbol string,
	funds float64,
	errorMessage string) ErrorEvent {

	fundsAsString := getFundsAsString(funds)

	return ErrorEvent{
		Timestamp:      getCurrentTs(),
		Server:         server,
		TransactionNum: transactionNum,
		Command:        command,
		Username:       username,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
		ErrorMessage:   errorMessage}
}