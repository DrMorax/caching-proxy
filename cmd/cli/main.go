package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/*", Proxy)
	http.ListenAndServe(":4000", nil)
	fmt.Println("Server started on port 4000")
}
