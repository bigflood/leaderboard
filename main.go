package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/usercount", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, `{"count":100}`)
		default:
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, `{"score":100,"rank":123,"updated_at":"..."}`)
		case http.MethodPut:
			fmt.Fprintf(w, `{"score":100,"rank":123,"updated_at":"..."}`)
		default:
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/ranks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, `{"ranks":[{"id":"user001","score":100,"rank":123,"updated_at":"..."},...]}`)
		default:
			http.NotFound(w, r)
		}
	})

	fmt.Println("listen..")
	http.ListenAndServe(":8080", nil)
}
