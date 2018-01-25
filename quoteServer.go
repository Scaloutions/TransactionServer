package main

import (
    "fmt"
    "net"
    "bufio"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng"
	PORT = 4444
	CONNECTION_TYPE = "tcp"
)

//function that will return the value of a requested stock
//func getQuote(stock string) float32 {
//	return 1
//}

func connectToQuoteServer(){
	//this should create a TCP connection with the quote server
	conn, err := net.Dial(CONNECTION_TYPE, QUOTE_SERVER_API)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
	}

	response, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Print("Data received: " + response)
}