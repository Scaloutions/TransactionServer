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
	PORT             = "44430"
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

	server := "CLT1"
	transactionNum := 1
	command := "QUOTE"

	// Get connection to the quote server
	conn := getConnection(server, transactionNum, command, file, userid, stock)

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
		errMsg := fmt.Sprintf("Cannot parse QS stock price %s into float64", quoteArgs[0])
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			userid,
			stock,
			0,
			errMsg)
		logging(errorEvent, file)
	}

	timestamp, err := strconv.ParseInt(quoteArgs[3], 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot parse QS timestamp %s into int64", quoteArgs[3])
		glog.Error(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			userid,
			stock,
			0,
			errMsg)
		logging(errorEvent, file)
	}

	cryptoKey := quoteArgs[4]

	quoteServer := getQuoteServer(
		server,
		transactionNum,
		timestamp,
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
		Timestamp: strconv.FormatInt(timestamp, 10),
		CryptoKey: quoteArgs[4],
	}

}

// server, transactionNum, command, file, userid, stock
func getConnection(
	server string,
	transactionNum int,
	command string,
	file *os.File,
	userid string,
	stock string) net.Conn {

	//this should create a TCP connection with the quote server
	glog.Info("Connecting to the quote server... ")
	conn, err := net.Dial(CONNECTION_TYPE, QUOTE_SERVER_API)

	if err != nil {
		errMsg := fmt.Sprintf("Error connecting to the Quote Server: something went wrong :(")
		fmt.Print(errMsg)

		errorEvent := getErrorEvent(
			server,
			transactionNum,
			command,
			userid,
			stock,
			0,
			errMsg)
		logging(errorEvent, file)
	}
	return conn
}
