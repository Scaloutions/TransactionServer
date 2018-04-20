package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"os"

	"github.com/golang/glog"
)

const (
	SERVER_NAME       = "TS0156"
	API_URL           = "/api"
	ACCOUNT_EVENT_URL = "/accounttransaction"
	SYSTEM_EVENT_URL  = "/systemevent"
	ERROR_EVENT_URL   = "/errorevent"
	QUOTE_SERVER_URL  = "/quoteserver"
	USER_COMMAND_URL  = "/usercommand"
	GET_DUMPLOG		  = "/api/log"
)

var (
	AUDIT_SERVER string
)

func InitializeAuditLogging() {
	testMode, _ := strconv.ParseBool(os.Getenv("DEV_ENVIRONMENT"))
	if testMode {
		AUDIT_SERVER  = os.Getenv("AUDIT_SERVER_DEV")
	} else {
		AUDIT_SERVER  = os.Getenv("AUDIT_SERVER_PROD")
	}
}

func getFundsAsString(amount float64) string {
	return fmt.Sprintf("%.2f", float64(amount))
}

func getCurrentTs() int64 {
	return time.Now().UnixNano() / 1000000
}

func getTransactionEvent(
	transactionNum int,
	action string,
	userId string,
	funds float64) AccountTransactionEvent {

	fundsAsString := getFundsAsString(funds)

	return AccountTransactionEvent{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Action:         action,
		UserId:         userId,
		Funds:          fundsAsString,
	}
}

func getSystemEvent(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64) SystemEvent {

	fundsAsString := getFundsAsString(funds)

	return SystemEvent{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         userId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
	}
}

func getQuoteServerEvent(
	transactionNum int,
	quoteServerTime int64,
	userId string,
	stockSymbol string,
	price float64,
	cryptokey string) QuoteServerEvent {

	priceAsString := getFundsAsString(price)

	return QuoteServerEvent{
		Timestamp:            getCurrentTs(),
		Server:               SERVER_NAME,
		TransactionNum:       transactionNum,
		QuoteServerEventTime: quoteServerTime,
		UserId:               userId,
		StockSymbol:          stockSymbol,
		Price:                priceAsString,
		Cryptokey:            cryptokey,
	}
}

func getErrorEvent(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64,
	errorMessage string) ErrorEvent {

	fundsAsString := getFundsAsString(funds)

	return ErrorEvent{
		Timestamp:      getCurrentTs(),
		Server:         SERVER_NAME,
		TransactionNum: transactionNum,
		Command:        command,
		UserId:         userId,
		StockSymbol:    stockSymbol,
		Funds:          fundsAsString,
		ErrorMessage:   errorMessage,
	}
}

func getUserCmndEvent(
	transactionNum int,
	command string,
	userId string,
	stockSymbol string,
	funds float64) UserCommand {
	fundsAsString := getFundsAsString(funds)

	return UserCommand {
		Timestamp:		getCurrentTs(),       
		Server:			SERVER_NAME,         
		TransactionNum: transactionNum, 
		Command: 		command,        
		UserId: 		userId,         
		StockSymbol: 	stockSymbol,    
		Funds: 			fundsAsString,          
	}
}

func logEvent(log interface{}) {

	sendLogs, _ := strconv.ParseBool(os.Getenv("LOG_EVENTS"))
	if !sendLogs {
		return
	}
	glog.Info("############## LOGGING REQUST: ", log)

	data, err := json.Marshal(log)
	URL := getUrlPath(log)

	if err != nil {
		glog.Error("Can not parse struct into JSON onject ", data)
	}
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	if err != nil {
		glog.Error("Error creating a request for the Audit server")
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		glog.Error("Error sending a POST request to Audit server")
		panic(err)
	}
	defer resp.Body.Close()
}

func getDumplog() {
	resp, err := http.Get(AUDIT_SERVER+GET_DUMPLOG)
	if err != nil {
		glog.Error("Error Sending DUMPLOG request to the Audit server")
		return
	}
	defer resp.Body.Close()
}

func getUrlPath(obj interface{}) string {
	var url bytes.Buffer
	url.WriteString(AUDIT_SERVER)
	url.WriteString(API_URL)

	switch obj.(type) {
	case SystemEvent:
		url.WriteString(SYSTEM_EVENT_URL)
	case AccountTransactionEvent:
		url.WriteString(ACCOUNT_EVENT_URL)
	case ErrorEvent:
		url.WriteString(ERROR_EVENT_URL)
	case QuoteServerEvent:
		url.WriteString(QUOTE_SERVER_URL)
	case UserCommand:
		url.WriteString(USER_COMMAND_URL)
	default:
		glog.Error("Error logging event to the audit server.")
		panic("Can not recognaize this type of event.")
	}

	return url.String()
}
