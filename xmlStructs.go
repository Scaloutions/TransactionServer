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
