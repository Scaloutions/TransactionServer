package api

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"errors"
	"github.com/golang/glog"
	"os"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng"
	PORT             = "4444"
	CONNECTION_TYPE  = "tcp"
)

var (
	QS_CONNECTION net.Conn
)

type Quote struct {
	Price     float64
	Stock     string
	UserId    string
	Timestamp int64
	CryptoKey string
}

/*func InitializeQSConn(){
	conn, err := getConnection() 
	QS_CONNECTION = conn

	if err!=nil {
		glog.Error("Cannot establish connection with the Quote Server ", err)
	}
}*/


func getQuoteFromQS(userid string, stock string) (Quote, error) {

	// Mock QuoteServer hit for local testing
	testMode, _ := strconv.ParseBool(os.Getenv("DEV_ENVIRONMENT"))
	if testMode {
		return Quote{
			Price:     1,
			Stock:     stock,
			UserId:    userid,
			Timestamp: 1516925116307,
			CryptoKey: "PXdxruf7H5p9Br19Si5hq",
		}, nil
	}

	quote := Quote{}
	// Get connection to the quote server
	// conn, err := getConnection()

	// if err!=nil {
	// QS_CONNECTION
	conn, err := getConnection()
	// if QS_CONNECTION == nil {
	if err!=nil {
		return quote, err
		// InitializeQSConn()
	}

	cstr := stock + "," + userid + "\n"
	_, err = conn.Write([]byte(cstr))

	if err!=nil {
		return quote, err
	}

	// //TODO: does this have o be 1024 bytes
	buff := make([]byte, 1024)
	// len, err := QS_CONNECTION.Read(buff)
	len, err := conn.Read(buff)

	if err != nil {
		glog.Error("Error reading data from the Quote Server")
		return quote, errors.New("Error reading the Quote.")	
	}

	response := string(buff[:len])
	glog.Info("Got back from Quote server: ", response)

	quoteArgs := strings.Split(response, ",")

	// Returns: quote,sym,userid,timestamp,cryptokey
	price, err := strconv.ParseFloat(quoteArgs[0], 64)
	if err != nil {
		glog.Error("Cannot parse QS stock price into float64 ", quoteArgs[0])
	 	return quote, errors.New("Error parsing the Quote.")
	}

	timestamp, err := strconv.ParseInt(quoteArgs[3], 10, 64)
	if err != nil {
		glog.Error("Cannot parse QS timestamp into int64 ", quoteArgs[3])
	 	return quote, errors.New("Error parsing the Quote.")
	}
	conn.Close()

	return Quote{
		Price:     price,
		Stock:     quoteArgs[1],
		UserId:    quoteArgs[2],
		Timestamp: timestamp,
		CryptoKey: strings.TrimSpace(quoteArgs[4]),
	}, nil



}

func getConnection() (net.Conn, error) {
	glog.Info("Connecting to the quote server... ")
	url := QUOTE_SERVER_API + ":" + PORT
	conn, err := net.Dial(CONNECTION_TYPE, url)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
		return nil,  errors.New("Cannot establish connection with the Quote Server")
	}
	return conn, nil
}
