package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng:4444"
	PORT             = "444430"
	CONNECTION_TYPE  = "tcp"
)

type Quote struct {
	Price     float64
	Stock     string
	UserId    string
	Timestamp string
	CryptoKey string
}

func getQuoteFromQS(userid string, stock string, file *os.File) Quote {
	// Get connection to the quote server
	conn := getConnection()

	cstr := stock + "," + userid + "\n"
	conn.Write([]byte(cstr))

	//TODO: does this have o be 1024 bytes
	buff := make([]byte, 1024)
	len, _ := conn.Read(buff)

	response := string(buff[:len])
	glog.Info("Got back: ", response)

	quoteArgs := strings.Split(response, ",")

	//example response: 254.69,OY0,S,1516925116307,PXdxruf7H5p9Br19Si5hq+tlsP24mj6hQQbDUZi8v+s=
	// Returns: quote,sym,userid,timestamp,cryptokey\n
	price, err := strconv.ParseFloat(quoteArgs[0], 64)
	if err != nil {
		glog.Error("Cannot parse QS stock price into float64 ", quoteArgs[0])
	}

	server := "CLT1"
	transactionNum := 1
	command := "ADD"
	cryptoKey := quoteArgs[4]
	var quoteServerTime int64

	quoteServer := getQuoteServer(
		server,
		transactionNum,
		quoteServerTime,
		command,
		userid,
		stock,
		price,
		cryptoKey)
	logging(quoteServer, file)

	return Quote{
		Price:     price,
		Stock:     quoteArgs[1],
		UserId:    quoteArgs[2],
		Timestamp: quoteArgs[3],
		CryptoKey: cryptoKey,
	}
}

func getConnection() net.Conn {
	//this should create a TCP connection with the quote server
	glog.Info("Connecting to the quote server... ")
	url := QUOTE_SERVER_API + ":" + PORT
	conn, err := net.Dial(CONNECTION_TYPE, url)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
	}
	return conn
}
