package main

import (
	"fmt"
	"net/http"
	"shorturl"
)

var u *shorturl.URLStore

func init() {
	u = shorturl.NewURLStore()
}

func main() {
	http.HandleFunc("/", Redirect)
	http.HandleFunc("/add", Add)

	// start server
	http.ListenAndServe(":8009", nil)
}

func Add(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		fmt.Fprint(w, AddForm)
		return
	}
	key := u.Put(url)
	fmt.Fprintf(w, "localhost:8009/%s", key)
}

const AddForm = `<!DOCTYPE html>
<form method=POST action=/add>
URL: <input type=text name=url>
<input type=submit value=Add>
</form>
`

func Redirect(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[1:]
	url := u.Get(key)
	if url == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}
