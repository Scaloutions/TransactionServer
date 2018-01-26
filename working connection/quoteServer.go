package main

import (
    "fmt"
    "net"
    "bufio"
    "strings"
)

const (
	QUOTE_SERVER_API = "quoteserve.seng"
	PORT = "4444"
	CONNECTION_TYPE = "tcp"
)

type Quote struct {
	A string
	B string
	C string
	D string
	CryptoKey string
}

//function that will return the value of a requested stock
//func getQuote(stock string) float32 {
//	return 1
//}

func getQuote() Quote {
	
	// Get connection to the quote server
	connQS := getConnection()

	// Listen on all interfaces
	ln, _ := net.Listen("tcp", ":44430")
	conn, _ := ln.Accept()

	//for{
		message, _ := bufio.NewReader(conn).ReadString('\n')
		cstr := string(message)
		fmt.Print("Message received: ", cstr)
		
		connQS.Write([]byte(cstr+"\n"))
	
		buff := make([]byte, 1024)
		len, _ := connQS.Read(buff)

		//gives back : 254.69,OY0,S,1516925116307,PXdxruf7H5p9Br19Si5hq+tlsP24mj6hQQbDUZi8v+s=

		response := string(buff[:len])
		fmt.Println("Got back: ", response)

		quoteArr := strings.Split(response, ",")
		
		return Quote {
			A: quoteArr[0],
			B: quoteArr[1],
			C: quoteArr[2],
			D: quoteArr[3],
			CryptoKey: quoteArr[4],
		}
	//}
	

	
}

func getConnection() net.Conn {
	//this should create a TCP connection with the quote server
	fmt.Println("Connecting to the quote server... ")	
	url := QUOTE_SERVER_API + ":" + PORT
	conn, err := net.Dial(CONNECTION_TYPE, url)

	if err != nil {
		fmt.Print("Error connecting to the Quote Server: somthing went wrong :(")
	}
	return conn
}

func main() {
	quote := getQuote()
	fmt.Printf("Received quote infomation: ", quote)
}






/*


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
//    	fmt.Print("Data received: " + response)
	conn.Close()
}

*/
