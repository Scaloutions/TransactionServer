
package main

import (
	"github.com/golang/glog"
    "fmt"
	"net"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng"
	PORT = "444430"
	CONNECTION_TYPE = "tcp"
)

type Quote struct {
	Price float64
	Stock string
	UserId string
	Timestamp int64
	CryptoKey string
}

func getQuoteFromQS(userid string, stock string) Quote {
	/*
	// Get connection to the quote server
	conn := getConnection()

	cstr := stock+","+userid+"\n"
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

	timestamp, err := strconv.ParseInt(quoteArgs[3], 10, 64)
	if err != nil {
		glog.Error("Cannot parse QS timestamp into int64 ", quoteArgs[3])
	}

	return Quote {
		Price: price,
		Stock: quoteArgs[1],
		UserId: quoteArgs[2],
		Timestamp: timestamp,
		CryptoKey: quoteArgs[4],
	}*/

	//this is just to mock quote server response for testing purposes
	return Quote {
		Price: 1,
		Stock: "S",
		UserId: "Agent007",
		Timestamp: 1516925116307,
		CryptoKey: "PXdxruf7H5p9Br19Si5hq",
	}

}

func getConnection() net.Conn {
	glog.Info("Connecting to the quote server... ")	
	url := QUOTE_SERVER_API + ":" + PORT
	conn, err := net.Dial(CONNECTION_TYPE, url)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
	}
	return conn
}
