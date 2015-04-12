package main

import (
	"io"
	"log"
	"net/http"
)

// define a handler HelloServe
// use the little Alpha,see the result

func helloServe(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello World \n")

}

func SayHello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "bye bye \n")
}
func main() {
	http.HandleFunc("/hello", helloServe)
	http.HandleFunc("/bye", SayHello)
	err := http.ListenAndServe(":8087", nil)

	if err != nil {
		log.Fatal("serve err : ", err)
	}
}
