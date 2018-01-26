
package main

import (
	"github.com/golang/glog"
    "fmt"
	"net"
   "strings"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng:4444"
	PORT = "444430"
	CONNECTION_TYPE = "tcp"
)

type Quote struct {
	Price string
	Stock string
	UserId string
	Timestamp string
	CryptoKey string
}

func getQuote(userid string, stock string) Quote {
	// Get connection to the quote server
	conn := getConnection()

	cstr := stock+","+userid+"\n"
	conn.Write([]byte(cstr))
	
	buff := make([]byte, 1024)
	len, _ := conn.Read(buff)

	response := string(buff[:len])
	glog.Info("Got back: ", response)

	quoteArgs := strings.Split(response, ",")
		
	//example response: 254.69,OY0,S,1516925116307,PXdxruf7H5p9Br19Si5hq+tlsP24mj6hQQbDUZi8v+s=
	// Returns: quote,sym,userid,timestamp,cryptokey\n
	return Quote {
		Price: quoteArgs[0],
		Stock: quoteArgs[1],
		UserId: quoteArgs[2],
		Timestamp: quoteArgs[3],
		CryptoKey: quoteArgs[4],
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
