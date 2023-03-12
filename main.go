package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

var cb *CircuitBreaker

func main() {
	log.Println("Hello, World")
	cb = NewCircuitBreaker(3, 3, 30*time.Second)
	http.HandleFunc("/hello", errorFunc)
	http.ListenAndServe(":8000", nil)
	log.Println("Bye")
}

func errorFunc(w http.ResponseWriter, r *http.Request) {
	cb.Execute(func() error {
		pass, _ := strconv.ParseBool(r.Header.Get("x-pass"))
		log.Print(pass)
		if pass {
			return nil
		}
		return errors.New("service returned some error")
	})
}
