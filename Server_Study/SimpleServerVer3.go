package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/one"] = FirstOne
	mux["/two"] = SecondOne
	wd, err := os.Getwd()
	if err != nill {
		log.Fatal(err)
	}
	mux["/static/"] = http.StripPrefix("/static/", http.FileServer(http.Dir(wd)))
	server := &http.Server{
		Addr:         ":8091",
		Handler:      &myHandle{},
		WriteTimeout: 5 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

type myHandle struct{}

func (*myHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}
	io.WriteString(w, "URL: "+r.URL.String())
}

func FirstOne(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is the FirstOne HandleFunc")
}
func SecondOne(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is the second one HandleFunc")
}
