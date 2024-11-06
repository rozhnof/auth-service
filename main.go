package main

import (
	"io"
	"net/http"
	"runtime"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	runtime.GOMAXPROCS(6)
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
