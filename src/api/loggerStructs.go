package api

import "encoding/xml"

type AccountTransactionEvent struct {
	XMLName        xml.Name `xml:"accountTransaction"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int      `xml:"transactionNum"`
	Action         string   `xml:"action"`
	UserId         string   `xml:"username,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type SystemEvent struct {
	XMLName        xml.Name `xml:"systemEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int      `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	UserId         string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}

type QuoteServerEvent struct {
	XMLName              xml.Name `xml:"quoteServer"`
	Timestamp            int64    `xml:"timestamp"`
	Server               string   `xml:"server"`
	TransactionNum       int      `xml:"transactionNum"`
	QuoteServerEventTime int64    `xml:"quoteServerTime"`
	UserId               string   `xml:"username,omitempty"`
	StockSymbol          string   `xml:"stockSymbol,omitempty"`
	Price                string   `xml:"price"`
	Cryptokey            string   `xml:"cryptokey"`
}

type ErrorEvent struct {
	XMLName        xml.Name `xml:"errorEvent"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int      `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	UserId         string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
	ErrorMessage   string   `xml:"errorMessage"`
}

type UserCommand struct {
	XMLName        xml.Name `xml:"userCommand"`
	Timestamp      int64    `xml:"timestamp"`
	Server         string   `xml:"server"`
	TransactionNum int      `xml:"transactionNum"`
	Command        string   `xml:"command,omitempty"`
	UserId         string   `xml:"username,omitempty"`
	StockSymbol    string   `xml:"stockSymbol,omitempty"`
	Funds          string   `xml:"funds,omitempty"`
}
