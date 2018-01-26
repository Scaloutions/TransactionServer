
package main

import (
    "fmt"
	"net"
//look into how to do requests with bufio
// might be easier to read ?
//    "bufio"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng:4444"
	//TODO: append port instead
	PORT = 4444
	CONNECTION_TYPE = "tcp"
)

//function that will return the value of a requested stock
//func getQuote(stock string) float32 {
//	return 1
//}

func connectToQuoteServer(){
	//this should create a TCP connection with the quote server
	fmt.Println("Connecting...")
	conn, err := net.Dial(CONNECTION_TYPE, QUOTE_SERVER_API)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
	}

	fmt.Println("Reading response...")

	//send request here
	cstr := "oY01WVirLr,S"

	conn.Write([]byte(cstr+"\n"))
	
	buff := make([]byte, 1024)
	len, _ := conn.Read(buff)

	//gives back : 254.69,OY0,S,1516925116307,PXdxruf7H5p9Br19Si5hq+tlsP24mj6hQQbDUZi8v+s=
	response := string(buff[:len])
	fmt.Println("Got back: ", response)

	//	response, _ := bufio.NewReader(conn).ReadString('\n')
	conn.Close()
}

//main here just for testing
func main() {
	connectToQuoteServer()
	fmt.Println("Check connection now...")	
}

