package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {

	mux := http.NewServeMux()
	mux.Handle("/", &myHandle{})
	wd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(wd))))
	err = http.ListenAndServe(":8090", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type myHandle struct{}

func (*myHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL : "+r.URL.String())
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
