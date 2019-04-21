package main

import (
	"math/rand"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		r := rand.Intn(100)
		if r < 20 {
			http.Error(w, "Random error", 500)
		} else {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		}
	})
	http.ListenAndServe(":8080", nil)
}
