package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi, there!")
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/getQuote", echoString)

	log.Fatal(http.ListenAndServe(":9090", nil))

}
