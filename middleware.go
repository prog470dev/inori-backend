package main

import (
	"fmt"
	"net/http"
)

func middle(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before f")
		f(w, r)
		fmt.Println("after f")
	}
}
