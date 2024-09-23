package main

import (
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/", fs)
}
