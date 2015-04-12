package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/myHandle", &myHandler{})
	//mux.Handle(pattern, handler)
	err := http.ListenAndServe(":8089", mux)
	if err != nil {
		log.Fatal(err)
	}

}

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "URL: "+r.URL.String())
	io.WriteString(w, "Hello world\n")
}
