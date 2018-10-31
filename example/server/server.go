package main

import "log"
import "net/http"

var counter = 0

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if counter >= 0 {
		w.WriteHeader(http.StatusInternalServerError)
		counter = 0
		return
	}
	// w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("pong")); err != nil {
		log.Printf("[WARN] can't send pong: %s", err)
	}
	// counter++
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
