package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	if err := http.ListenAndServe(":9999", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("host: %v path: %v\n", r.Host, r.URL.Path)
		for k, v := range r.Header {
			fmt.Printf("header %v : %v\n", k, v)
		}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("body: %s\n", bytes)
	})); err != nil {
		panic(err)
	}
}
